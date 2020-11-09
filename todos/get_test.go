package todos

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/reltest"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	var (
		ctx             = context.TODO()
		repository      = reltest.New()
		service         = New(repository, nil)
		id         uint = 1
		result          = Todo{ID: id, Title: "Sleep"}
		todo       Todo
	)

	repository.ExpectFind(where.Eq("id", id)).Result(result)

	assert.NotPanics(t, func() {
		err := service.Get(ctx, &todo, id)
		assert.Equal(t, result, todo)
		assert.Nil(t, err)
	})

	repository.AssertExpectations(t)
}

func TestGetErr(t *testing.T) {
	var (
		ctx             = context.TODO()
		repository      = reltest.New()
		service         = New(repository, nil)
		id         uint = 1
		todo       *Todo
	)

	repository.ExpectFind(
		where.Eq("id", id),
	).Error(rel.ErrNotFound)

	assert.NotPanics(t, func() {
		err := service.Get(ctx, todo, id)
		assert.Equal(t, rel.ErrNotFound, err)
		assert.Nil(t, todo)
	})

	repository.AssertExpectations(t)
}
