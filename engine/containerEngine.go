package engine

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/containerd/go-runc"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "strconv"

    //"log"
    "golang.org/x/sys/unix"
    "impulse/chamber"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
)

func init() {
    log.SetReportCaller(false)
    log.SetLevel(log.DebugLevel)
}

type ContainerEngine struct {
    runc                *runc.Runc
    containerRuntimeDir string
    baseImageDir        string
    guestImageDir       string
    portAllocator       *PortAllocator
}

func NewContainerEngine() (*ContainerEngine, error) {
    runc := &runc.Runc{
        Root:         "/run/runc",
        Log:          "/tmp/runclog.log",
        LogFormat:    runc.Text,
        PdeathSignal: unix.SIGKILL,
        Debug:        true,
    }
    engine := ContainerEngine{
        runc:                runc,
        containerRuntimeDir: "/opt/impulse/run/containers", // TODO: why does /var/run/impulse not work?
        baseImageDir:        "/vagrant/images/base",
        guestImageDir:       "/vagrant/images/guest",
        portAllocator:       NewPortAllocator(5000, 9999),
    }
    return &engine, nil
}

func (ce *ContainerEngine) getContainerId(spec chamber.Spec) string {
    return fmt.Sprintf("%s-%s", spec.App, spec.Version)
}

func (ce *ContainerEngine) getBundlePath(spec chamber.Spec) string {
    return filepath.Join(ce.containerRuntimeDir, ce.getContainerId(spec))
}

func (ce *ContainerEngine) getRootFsPath(spec chamber.Spec) string {
    return filepath.Join(ce.getBundlePath(spec), "rootfs")
}

func (ce *ContainerEngine) Create(spec chamber.Spec) (status chamber.Status, err error) {

    status.Id = ce.getContainerId(spec)
    bundlePath := ce.getBundlePath(spec)
    rootfsPath := ce.getRootFsPath(spec)

    log := log.WithFields(log.Fields{"containerId": status.Id})

    // Create container bundle and rootfs path
    if err := os.MkdirAll(rootfsPath, 0755); err != nil {
        return status, errors.Wrapf(err, "failed to create container root directory: %q", bundlePath)
    }

    // Create and write the oci container spec as config.json for runc
    status.Port = ce.portAllocator.GetPort()
    containerSpec, err := getOciContainerSpec(spec, status.Id, status.Port, rootfsPath)
    if err != nil {
        return status, errors.Wrapf(err, "failed to generate oci container spec")
    }
    containerSpecBytes, err := json.Marshal(containerSpec)
    if err != nil {
        return status, errors.Wrap(err, "failed to marshal container spec to json")
    }

    runcConfigPath := filepath.Join(bundlePath, "config.json")
    if err := ioutil.WriteFile(runcConfigPath, containerSpecBytes, 0755); err != nil {
        return status, errors.Wrapf(err, "Failed to write runc config to %s", runcConfigPath)
    }

    // Setup base image
    baseImagePath := filepath.Join(ce.baseImageDir, fmt.Sprintf("%s.tar", spec.Runtime))
    if _, err := os.Stat(baseImagePath); os.IsNotExist(err) {
        return status, fmt.Errorf("unable to find base image at path %s", baseImagePath)
    }
    baseImageExtractCmd := exec.Command("tar", "-C", rootfsPath, "-xf", baseImagePath)

    // TODO: Costs 300 millis. Clever filesystem tricks could make this faster
    log.Debugf("Extracting base image from %s to %s...", baseImagePath, rootfsPath)
    if err := baseImageExtractCmd.Run(); err != nil {
        return status, errors.Wrapf(err, "failed to extract base image %s to rootfs at %s", baseImagePath, rootfsPath)
    }

    // Setup guest image
    guestImagePath := filepath.Join(ce.guestImageDir, fmt.Sprintf("%s.tar", status.Id))
    if _, err := os.Stat(guestImagePath); os.IsNotExist(err) {
        return status, fmt.Errorf("unable to find guest image at path %s", guestImagePath)
    }
    guestDir := filepath.Join(rootfsPath, "/opt/skyhook/guest") // TODO: make this configurable per base image?
    if err := os.MkdirAll(guestDir, 0755); err != nil {
        return status, errors.Wrapf(err, "failed to create guest dir inside contaier at %s", guestDir)
    }
    guestImageExtractCmd := exec.Command("tar", "-C", guestDir, "-xf", guestImagePath)
    log.Debugf("Extracting guest image from %s to %s...", guestImagePath, guestDir)
    if err := guestImageExtractCmd.Run(); err != nil {
        return status, errors.Wrapf(err, "failed to extract guest image %s to rootfs at %s", guestImagePath, guestDir)
    }

    containerIO, err := runc.NewSTDIO()
    if err != nil {
        return status, errors.Wrap(err, "failed to create container io")
    }

    opts := &runc.CreateOpts{
        IO:      containerIO,
        PidFile: filepath.Join(bundlePath, "init.pid"),
    }

    go func() {
        // cleanup
        defer func() {
            ce.portAllocator.ReleasePort(status.Port)
            if err := os.RemoveAll(bundlePath); err != nil {
                log.Errorf("Failed to clean up container bundle at %s: %v", bundlePath, err)
            }
            log.Infof("Cleaup complete")
        }()
         
        
        // TODO: Create container in separate step then start container here, to catch errors earlier and return
        log.Debug("Running runc container...")
        i, err := ce.runc.Run(context.Background(), status.Id, bundlePath, opts)
        if err != nil {
            //return fmt.Errorf("failed to create container: %v", err)
            log.Error("failed to create container: ", err)
        }
        log.Infof("Container exited with code %d", i)
    }()

    status.Status = "STARTING"
    log.Debug("Container creation started")
    return status, nil
}

func (ce *ContainerEngine) List() (statuses []chamber.Status, err error) {
    containerStatuses, err := ce.runc.List(context.Background())
    if err != nil {
        return statuses, errors.Wrapf(err, "failed to list runc containers")
    }
    for _, runcStatus := range(containerStatuses) {
        runcState, err := ce.runc.State(context.Background(), runcStatus.ID)
        if err != nil {
            return statuses, errors.Wrapf(err, "failed to get runc state for container %s", runcState.ID)
        }
        
        // Extract chamber spec from annotations
        spec := chamber.Spec{}
        chamberSpecKey := "skyhook_spec" // TODO: constant or something
        chamberSpecJson, ok := runcState.Annotations[chamberSpecKey]
        if !ok {
            return statuses, fmt.Errorf("container %s missing annotation %s", runcState.ID, chamberSpecKey)
        }
        b := []byte(chamberSpecJson)
        if err := json.Unmarshal(b, &spec); err != nil {
            return statuses, errors.Wrapf(err, "failed to unmarshal spec from runc state for container %s", runcState.ID)
        }
        
        // Extract port from annotations
        portKey := "skyhook_port"
        portStr, ok := runcState.Annotations[portKey]
        if !ok {
            return statuses, fmt.Errorf("container %s missing annotation %s", runcState.ID, portKey)
        }
        port, err := strconv.Atoi(portStr)
        if err != nil {
            return statuses, errors.Wrapf(err, "failed to convert port value %s to int", port)
        }
        
        status := chamber.Status{
            Id: runcState.ID,
            Status: runcState.Status,
            CreatedTimeMillis: runcState.Created.UnixNano() / 1000000,
            Spec: spec,
            Port: port,
        }
        statuses = append(statuses, status)
    }
    return statuses, nil
}
