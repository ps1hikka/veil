package auth

import (
	"encoding/hex"
	"fmt"
	"strings"

	"gateway/internal/config"
)

type Authenticator struct {
	allowed map[[16]byte]config.Client
}

func New(cfg config.Config) (*Authenticator, error) {
	m := make(map[[16]byte]config.Client, len(cfg.Clients))
	for _, c := range cfg.Clients {
		id, err := parseUUID(c.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid client id %q: %w", c.ID, err)
		}
		m[id] = c
	}
	return &Authenticator{allowed: m}, nil
}

func (a *Authenticator) Allow(id [16]byte) (config.Client, bool) {
	c, ok := a.allowed[id]
	return c, ok
}

func parseUUID(s string) ([16]byte, error) {
	var out [16]byte
	b, err := hex.DecodeString(strings.ReplaceAll(s, "-", ""))
	if err != nil {
		return out, err
	}
	if len(b) != 16 {
		return out, fmt.Errorf("invalid UUID length: %d", len(b))
	}
	copy(out[:], b)
	return out, nil
}
