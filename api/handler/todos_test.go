package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/reltest"
	"github.com/h4ckm03d/go-todo-cleanarch/api/handler"
	"github.com/h4ckm03d/go-todo-cleanarch/todos"
	"github.com/h4ckm03d/go-todo-cleanarch/todos/todostest"
	"github.com/stretchr/testify/assert"
)

func TestTodos_Index(t *testing.T) {
	var (
		trueb = true
	)

	tests := []struct {
		name            string
		status          int
		path            string
		response        string
		mockTodosSearch func(todos *todostest.Service)
	}{
		{
			name:     "ok",
			status:   http.StatusOK,
			path:     "/",
			response: `[{"id":1, "title":"Sleep", "completed":false, "order":0, "url":"todos/1", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}]`,
			mockTodosSearch: todostest.MockSearch(
				[]todos.Todo{{ID: 1, Title: "Sleep"}},
				todos.Filter{},
				nil,
			),
		},
		{
			name:     "with keyword and filter completed",
			status:   http.StatusOK,
			path:     "/?keyword=Wake&completed=true",
			response: `[{"id":2, "title":"Wake", "completed":true, "order":0, "url":"todos/2", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}]`,
			mockTodosSearch: todostest.MockSearch(
				[]todos.Todo{{ID: 2, Title: "Wake", Completed: true}},
				todos.Filter{Keyword: "Wake", Completed: &trueb},
				nil,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("GET", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosSearch)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.JSONEq(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Create(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		path            string
		payload         string
		response        string
		location        string
		mockTodosCreate func(todos *todostest.Service)
	}{
		{
			name:     "created",
			status:   http.StatusCreated,
			path:     "/",
			payload:  `{"title": "Sleep"}`,
			response: `{"id":1, "title":"Sleep", "completed":false, "order":0, "url":"todos/1", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}`,
			location: "/1",
			mockTodosCreate: todostest.MockCreate(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
		},
		{
			name:     "validation error",
			status:   http.StatusUnprocessableEntity,
			path:     "/",
			payload:  `{"title": ""}`,
			response: `{"error":"Title can't be blank"}`,
			mockTodosCreate: todostest.MockCreate(
				todos.Todo{Title: "Sleep"},
				todos.ErrTodoTitleBlank,
			),
		},
		{
			name:     "bad request",
			status:   http.StatusBadRequest,
			path:     "/",
			payload:  ``,
			response: `{"error":"Bad Request"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				body       = strings.NewReader(test.payload)
				req, _     = http.NewRequest("POST", test.path, body)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosCreate)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.Equal(t, test.location, rr.Header().Get("Location"))
			assert.JSONEq(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Show(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		path         string
		response     string
		isPanic      bool
		mockTodosGet func(todos *todostest.Service)
	}{
		{
			name:     "ok",
			status:   http.StatusOK,
			path:     "/1",
			response: `{"id":1, "title":"Sleep", "completed":false, "order":0, "url":"todos/1", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}`,
			mockTodosGet: todostest.MockGet(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
		},
		{
			name:     "not found",
			status:   http.StatusNotFound,
			path:     "/1",
			response: `{"error":"Record not found"}`,
			mockTodosGet: todostest.MockGet(
				todos.Todo{Title: "Sleep"},
				rel.ErrNotFound,
			),
		},
		{
			name:    "panic",
			path:    "/1",
			isPanic: true,
			mockTodosGet: todostest.MockGet(
				todos.Todo{Title: "Sleep"},
				errors.New("other errors"),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("GET", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosGet)

			if test.isPanic {
				assert.Panics(t, func() {
					handler.ServeHTTP(rr, req)
				})
			} else {
				handler.ServeHTTP(rr, req)
				assert.Equal(t, test.status, rr.Code)
				assert.JSONEq(t, test.response, rr.Body.String())
			}

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Update(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		path            string
		payload         string
		response        string
		mockTodosUpdate func(todos *todostest.Service)
		mockTodosGet    func(todos *todostest.Service)
	}{
		{
			name:     "ok",
			status:   http.StatusOK,
			path:     "/1",
			payload:  `{"title": "Wake"}`,
			response: `{"id":1, "title":"Wake", "completed":false, "order":0, "url":"todos/1", "created_at":"0001-01-01T00:00:00Z", "updated_at":"0001-01-01T00:00:00Z"}`,
			mockTodosUpdate: todostest.MockUpdate(
				todos.Todo{ID: 1, Title: "Wake"},
				nil,
			),
			mockTodosGet: todostest.MockGet(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
		},
		{
			name:     "validation error",
			status:   http.StatusUnprocessableEntity,
			path:     "/1",
			payload:  `{"title": ""}`,
			response: `{"error":"Title can't be blank"}`,
			mockTodosGet: todostest.MockGet(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
			mockTodosUpdate: todostest.MockUpdate(
				todos.Todo{ID: 1, Title: ""},
				todos.ErrTodoTitleBlank,
			),
		},
		{
			name:     "bad request",
			status:   http.StatusBadRequest,
			path:     "/1",
			payload:  ``,
			response: `{"error":"Bad Request"}`,
			mockTodosGet: todostest.MockGet(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				body       = strings.NewReader(test.payload)
				req, _     = http.NewRequest("PATCH", test.path, body)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosUpdate, test.mockTodosGet)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.JSONEq(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Destroy(t *testing.T) {
	tests := []struct {
		name            string
		status          int
		path            string
		response        string
		mockTodosDelete func(todos *todostest.Service)
		mockTodosGet    func(todos *todostest.Service)
	}{
		{
			name:     "ok",
			status:   http.StatusNoContent,
			path:     "/1",
			response: "",
			mockTodosGet: todostest.MockGet(
				todos.Todo{ID: 1, Title: "Sleep"},
				nil,
			),
			mockTodosDelete: todostest.MockDelete(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("DELETE", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosDelete, test.mockTodosGet)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			assert.Equal(t, test.response, rr.Body.String())

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}

func TestTodos_Clear(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		path           string
		response       string
		mockTodosClear func(todos *todostest.Service)
	}{
		{
			name:           "created",
			status:         http.StatusNoContent,
			path:           "/",
			response:       "",
			mockTodosClear: todostest.MockClear(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				req, _     = http.NewRequest("DELETE", test.path, nil)
				rr         = httptest.NewRecorder()
				repository = reltest.New()
				todos      = &todostest.Service{}
				handler    = handler.NewTodos(repository, todos)
			)

			todostest.Mock(todos, test.mockTodosClear)

			handler.ServeHTTP(rr, req)

			assert.Equal(t, test.status, rr.Code)
			if test.response != "" {
				assert.JSONEq(t, test.response, rr.Body.String())
			} else {
				assert.Equal(t, "", rr.Body.String())
			}

			repository.AssertExpectations(t)
			todos.AssertExpectations(t)
		})
	}
}
