package engine

import (
    "fmt"
    "impulse/chamber"
    "log"
)

type ContainerEngine struct {
    
}

func NewContainerEngine() (*ContainerEngine, error){
    return &ContainerEngine{}, nil
}

func (ce *ContainerEngine) Create(spec chamber.Spec) error {
    log.Printf("creating vm for app %s..", spec.App)
    return fmt.Errorf("Failed to create")
}

func (ce *ContainerEngine) List() ([]chamber.Status, error) {
    return []chamber.Status{
        {
            Status: "running",
        },
    }, nil
}
