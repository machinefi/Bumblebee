package id

import (
	"context"
	"net"
	"time"

	"github.com/iotexproject/Bumblebee/base/types"
	"github.com/iotexproject/Bumblebee/base/types/snowflake_id"
	"github.com/iotexproject/Bumblebee/x/contextx"
	"github.com/iotexproject/Bumblebee/x/misc/must"
)

var (
	start, _ = time.Parse(time.RFC3339, "2022-10-13T22:10:24Z")
	sff      = snowflake_id.NewSnowflakeFactory(12, 10, 1, start)
)

type Generator interface {
	ID() (uint64, error)
}

func FromIP(ip net.IP) (Generator, error) {
	return sff.NewSnowflake(snowflake_id.WorkerIDFromIP(ip))
}

func FromLocalIP() (Generator, error) {
	return sff.NewSnowflake(snowflake_id.WorkerIDFromLocalIP())
}

type keyGenerator struct{}

func WithGenerator(ctx context.Context, g Generator) context.Context {
	return contextx.WithValue(ctx, keyGenerator{}, g)
}

func WithGeneratorContext(g Generator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, keyGenerator{}, g)
	}
}

func GeneratorFromContext(ctx context.Context) (Generator, bool) {
	g, ok := ctx.Value(keyGenerator{}).(Generator)
	return g, ok
}

func MustGeneratorFromContext(ctx context.Context) Generator {
	g, ok := ctx.Value(keyGenerator{}).(Generator)
	must.BeTrue(ok)
	return g
}

type SFIDGenerator interface {
	NewSFID() (types.SFID, error)
	NewSFIDs(n int) (types.SFIDs, error)
}

type keySFIDGenerator struct{}

func WithSFIDGenerator(ctx context.Context, g SFIDGenerator) context.Context {
	return contextx.WithValue(ctx, keySFIDGenerator{}, g)
}

func WithSFIDGeneratorContext(g SFIDGenerator) contextx.WithContext {
	return func(ctx context.Context) context.Context {
		return contextx.WithValue(ctx, keySFIDGenerator{}, g)
	}
}

func SFIDGeneratorFromContext(ctx context.Context) (SFIDGenerator, bool) {
	g, ok := ctx.Value(keySFIDGenerator{}).(SFIDGenerator)
	return g, ok
}

func MustSFIDGeneratorFromContext(ctx context.Context) SFIDGenerator {
	g, ok := ctx.Value(keySFIDGenerator{}).(SFIDGenerator)
	must.BeTrue(ok)
	return g
}
