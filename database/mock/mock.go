package mock

import (
	"context"

	"github.com/jamesog/example-http-api/database"
)

type mockdb struct{}

// NewDB returns a new mock instance.
func NewDB() (*mockdb, error) {
	return &mockdb{}, nil
}

func (m *mockdb) Users(context.Context) ([]database.User, error) {
	u := []database.User{
		database.User{
			ID:   1,
			Name: "Alice Example",
		},
	}
	return u, nil
}
