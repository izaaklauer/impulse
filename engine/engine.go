package engine

import (
    "impulse/chamber"
)

type Engine interface {
    Create(spec chamber.Spec) error
    List() ([]chamber.Status, error)
}
