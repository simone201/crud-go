package api

import (
	mocks "HTTPChiSqlite/mocks/db"
	"HTTPChiSqlite/model"
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListPeople(t *testing.T) {
	repo := new(mocks.PeopleRepository)

	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, &time.Location{})
	createdDate := time.Date(2001, 1, 1, 0, 0, 0, 0, &time.Location{})
	updatedDate := time.Date(2002, 1, 1, 0, 0, 0, 0, &time.Location{})
	repo.On("FindAllPeople").Return([]model.Person{
		{Id: 0, Name: "a", Birth: birthDate, CreatedAt: createdDate, UpdatedAt: updatedDate},
		{Id: 1, Name: "b", Birth: birthDate, CreatedAt: createdDate, UpdatedAt: updatedDate},
	}, nil)

	peopleRepository = repo
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/people", nil)

	rctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	listPeople(w, r)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t,
		"[{\"id\": 0, \"name\": \"a\", \"birth\": \"2000-01-01T00:00:00Z\", \"createdAt\": \"2001-01-01T00:00:00Z\", \"updatedAt\": \"2002-01-01T00:00:00Z\"}, "+
			"{\"id\": 1, \"name\": \"b\", \"birth\": \"2000-01-01T00:00:00Z\", \"createdAt\": \"2001-01-01T00:00:00Z\", \"updatedAt\": \"2002-01-01T00:00:00Z\"}]",
		func(body io.ReadCloser) string {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(body)
			return buf.String()
		}(w.Result().Body),
	)
}
