package engine

import (
    "impulse/chamber"
)

type Engine interface {
    Create(spec chamber.Spec) (chamber.Status, error)
    List() ([]chamber.Status, error)
}
