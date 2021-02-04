package engine

import (
    "context"
    "impulse/chamber"
)

type Engine interface {
    Create(ctx context.Context, spec chamber.Spec) error
    List(ctx context.Context) ([]chamber.Status, error)
}
