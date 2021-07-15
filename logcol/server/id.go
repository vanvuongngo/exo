package server

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	"github.com/deref/exo/chrono"
	ulid "github.com/oklog/ulid/v2"
)

type idGen struct {
	mu      sync.Mutex
	entropy *ulid.MonotonicEntropy
}

func newIdGen() *idGen {
	seed := int64(chrono.NowNano(context.Background()))
	randRead := rand.New(rand.NewSource(seed))
	return &idGen{
		entropy: ulid.Monotonic(randRead, 0),
	}
}

func (gen *idGen) nextId(ctx context.Context) ([]byte, error) {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	ts := ulid.Timestamp(chrono.Now(ctx))
	id, err := ulid.New(ts, gen.entropy)
	if err != nil {
		return nil, err
	}
	return id.MarshalBinary()
}

func parseID(id []byte) (string, error) {
	var asULID ulid.ULID
	if copy(asULID[:], id) != 16 {
		return "", errors.New("invalid length")
	}

	return asULID.String(), nil
}
