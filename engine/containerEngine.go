package engine

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/containerd/go-runc"
    "github.com/pkg/errors"
    log "github.com/sirupsen/logrus"
    "golang.org/x/sys/unix"
    "impulse/chamber"
    "io"
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
    runc *runc.Runc
    containerDir string
}

func NewContainerEngine() (*ContainerEngine, error){
    runc := &runc.Runc {
        Root: "/run/runc",
        Log: "/tmp/runclog.log",
        LogFormat: runc.Text,
        PdeathSignal: unix.SIGKILL,
        Debug: true,
    }
    engine := ContainerEngine{
        runc: runc,
        containerDir: "/opt/impulse/containers",
    }
    return &engine, nil
}

func (ce *ContainerEngine) Create(ctx context.Context, spec chamber.Spec) error {
    containerId := fmt.Sprintf("%s-%s",spec.App, spec.Version)
    
    log := log.WithFields(log.Fields{"containerId": containerId})
    
    //bundlePath := filepath.Join(ce.runc.Root, containerId)
    bundlePath := filepath.Join(ce.containerDir, containerId)
    rootfsPath := filepath.Join(bundlePath, "rootfs")
    
    // Create container bundle and rootfs path
    if err := os.MkdirAll(rootfsPath, 0755); err != nil {
        return errors.Wrapf(err, "failed to create container root directory: %q", bundlePath)
    }
    
    // Create and write the oci container spec as config.json for runc
    containerSpec := getContainerSpec(containerId, rootfsPath)
    containerSpecBytes, err := json.Marshal(containerSpec)
    if err != nil {
        return errors.Wrap(err, "failed to marshal container spec to json")
    }
    ioutil.WriteFile(filepath.Join(bundlePath, "config.json"), containerSpecBytes, 0755)
    
    // Add the rootfs to the bundle
    archivePath := "/vagrant/python-alpine-docker-export.tar"
    cmd := exec.Command("tar", "-C", rootfsPath, "-xvf", archivePath)
    
    log.Debugf("Extracting rootfs from %s to %s...", archivePath, rootfsPath)
    if err := cmd.Run(); err != nil {
        return errors.Wrapf(err, "failed to extract archive %s to rootfs at %s", archivePath, rootfsPath)
    }
    
    containerIO, err := runc.NewSTDIO()
    if err != nil {
        return errors.Wrap(err, "failed to create container io")
    }
    
    opts := &runc.CreateOpts{
        IO: containerIO,
        PidFile: filepath.Join(bundlePath, "init.pid"),
    }
    log.Debug("Creating runc container...")
    
    
    // TODO: Make this log streaming optional
    go func() {
        _, err := io.Copy(os.Stdout, containerIO.Stdout())
        if err != nil {
           log.Error("failed to read stdout: ", err)
        }
    }()

    go func() {
        _, err := io.Copy(os.Stderr, containerIO.Stderr())
        if err != nil {
            log.Error("failed to read stderr: ", err)
        }
    }()
    
    // TODO: runc create does not return (https://github.com/containerd/go-runc/issues/31). Forking a new gorouine is a bad solution.
    i, err := ce.runc.Run(context.Background(), containerId, bundlePath, opts);
    if err != nil {
        //return fmt.Errorf("failed to create container: %v", err)
        log.Error("failed to create container: ", err)
    }
    log.Infof("Container exited with code %d", i)
    log.Debug("Created container")
    
    return nil
}



func (ce *ContainerEngine) List(ctx context.Context) ([]chamber.Status, error) {
    return []chamber.Status{
        {
            Status: "running",
        },
    }, nil
}
