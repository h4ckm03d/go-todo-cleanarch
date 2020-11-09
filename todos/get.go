package todos

import (
	"context"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
)

type get struct {
	repository rel.Repository
}

func (g get) Get(ctx context.Context, todo *Todo, id uint) error {
	if err := g.repository.Find(ctx, todo, where.Eq("id", id)); err != nil {
		return err
	}

	return nil
}
