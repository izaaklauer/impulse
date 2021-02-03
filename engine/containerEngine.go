package engine

import (
    "fmt"
    "github.com/opencontainers/runc/libcontainer"
    _ "github.com/opencontainers/runc/libcontainer/nsenter"
    "github.com/sirupsen/logrus"
    "impulse/chamber"
    "log"
    "os"
    "runtime"
)

func init() {
   if len(os.Args) > 1 && os.Args[1] == "init" {
       runtime.GOMAXPROCS(1)
       runtime.LockOSThread()
       factory, _ := libcontainer.New("")
       if err := factory.StartInitialization(); err != nil {
           logrus.Fatal(err)
       }
       panic("--this line should have never been executed, congratulations--")
   }
}

type ContainerEngine struct {
    factory libcontainer.Factory
}

func NewContainerEngine() (*ContainerEngine, error){
    factory, err := libcontainer.New("/opt/impulse/containers", libcontainer.Cgroupfs, libcontainer.InitArgs(os.Args[0], "init"))
    if err != nil {
        return nil, fmt.Errorf("failed to instantiate libcontainer factory: %v", err)
    }
    engine := ContainerEngine{
        factory: factory,
    }
    return &engine, nil
}

func (ce *ContainerEngine) Create(spec chamber.Spec) error {
    log.Printf("creating vm for app %s..", spec.App)
    
    container, err := ce.factory.Create(spec.App, containerConfig("/root/pythontest/rootfs"))
    if err != nil {
        return fmt.Errorf("failed to create container: %v", err)
    }
    
    log.Printf("Got this container: %+v", container)
    
    return nil
}

func (ce *ContainerEngine) List() ([]chamber.Status, error) {
    return []chamber.Status{
        {
            Status: "running",
        },
    }, nil
}
