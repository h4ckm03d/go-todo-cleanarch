package todostest

import (
	context "context"

	rel "github.com/go-rel/rel"
	todos "github.com/h4ckm03d/go-todo-cleanarch/todos"
	mock "github.com/stretchr/testify/mock"
)

// MockFunc function.
type MockFunc func(service *Service)

// Mock apply mock todo functions.
func Mock(service *Service, funcs ...MockFunc) {
	for i := range funcs {
		if funcs[i] != nil {
			funcs[i](service)
		}
	}
}

// MockSearch util.
func MockSearch(result []todos.Todo, filter todos.Filter, err error) MockFunc {
	return func(service *Service) {
		service.On("Search", mock.Anything, mock.Anything, filter).
			Return(func(ctx context.Context, out *[]todos.Todo, filter todos.Filter) error {
				*out = result
				return err
			})
	}
}

// MockCreate util.
func MockCreate(result todos.Todo, err error) MockFunc {
	return func(service *Service) {
		service.On("Create", mock.Anything, mock.Anything).
			Return(func(ctx context.Context, out *todos.Todo) error {
				*out = result
				return err
			})
	}
}

// MockUpdate util.
func MockUpdate(result todos.Todo, err error) MockFunc {
	return func(service *Service) {
		service.On("Update", mock.Anything, mock.Anything, mock.Anything).
			Return(func(ctx context.Context, out *todos.Todo, changeset rel.Changeset) error {
				if result.ID != out.ID {
					panic("inconsistent id")
				}

				*out = result
				return err
			})
	}
}

// MockClear util.
func MockClear() MockFunc {
	return func(service *Service) {
		service.On("Clear", mock.Anything)
	}
}

// MockDelete util.
func MockDelete() MockFunc {
	return func(service *Service) {
		service.On("Delete", mock.Anything, mock.Anything)
	}
}

// MockGet util.
func MockGet(result todos.Todo, err error) MockFunc {
	return func(service *Service) {
		service.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Return(func(ctx context.Context, out *todos.Todo, id uint) error {
				*out = result
				return err
			})
	}
}
