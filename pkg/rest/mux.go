package rest

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

type Mux struct {
	chi.Mux
	Validator   StructValidator
	Transformer *mold.Transformer
}

func NewMux() *Mux {
	return &Mux{
		Mux:         *chi.NewRouter(),
		Validator:   NewStructValidator(),
		Transformer: modifiers.New(),
	}
}

// PrepareParams transforms and validates parameters. Submit only pointer values.
func (m Mux) PrepareParams(ctx context.Context, params any) error {
	if err := m.Transformer.Struct(ctx, params); err != nil {
		return err
	}
	return m.Validator.Validate(ctx, params)
}
