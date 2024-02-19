package handlers

import (
	"github.com/izzanzahrial/skeleton/internal/interface/http/authentication"
	"github.com/izzanzahrial/skeleton/internal/interface/http/user"
)

type Handlers struct {
	Auth *authentication.Handler
	User *user.Handler
}

// type HandlersConfiguration func(h *Handlers) error

// func New(cfgs ...HandlersConfiguration) (*Handlers, error) {
// 	h := &Handlers{}

// 	for _, cfg := range cfgs {
// 		if err := cfg(h); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return h, nil
// }

// func WithAuthHandler(ah *authentication.Handler) HandlersConfiguration {
// 	return func(h *Handlers) error {
// 		h.Auth = ah
// 		return nil
// 	}
// }

func NewHandlers(ah *authentication.Handler, uh *user.Handler) *Handlers {
	return &Handlers{
		Auth: ah,
		User: uh,
	}
}
