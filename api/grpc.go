package api

//go:generate protoc api.proto --go_out=plugins=grpc:.

import (
	context "context"
)

// GetUsers is the protobuf interface for retrieving users.
func (api *API) GetUsers(ctx context.Context, in *UserRequest) (*UserList, error) {
	dbUsers, err := api.db.Users(ctx)
	if err != nil {
		return nil, err
	}

	// Protobuf calls the ID field Id, so we can't just do a type conversion
	// from the database.User type to the Protobuf type. Instead we loop over
	// the returned slice of users, creating a new slice and returning that
	// wrapped in the protobuf UserList.
	//
	// This could perhaps be tidied up by creating a new access method in the
	// postgres package, but that would introduce dependencies between the api
	// and postgres packages which is not desirable. The Protobuf definitions
	// could equally go in their own package outside of the api package for
	// true separation of concerns, allowing postgres to link to it.
	users := make([]*User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = &User{
			Id:   int64(u.ID),
			Name: u.Name,
		}
	}

	return &UserList{Users: users}, nil
}
