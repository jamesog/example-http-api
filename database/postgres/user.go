package postgres

import (
	"context"

	"github.com/jamesog/example-http-api/database"
)

// Users returns a user from the database.
func (db *DB) Users(ctx context.Context) ([]database.User, error) {
	rows, err := db.QueryxContext(ctx, "SELECT id, name FROM app_user")
	if err != nil {
		db.log.Error().Err(err).Msg("failed to query database")
		return nil, err
	}

	// Use a zero-length slice rather than an uninitialized slice.
	// If the slice is uninitialized it will be nil. JSON encoding
	// in the API will then return "null" instead of "[]".
	users := make([]database.User, 0)
	for rows.Next() {
		var user database.User
		if err := rows.StructScan(&user); err != nil {
			db.log.Err(err).Msg("StructScan error")
		}

		users = append(users, user)
	}

	return users, nil
}
