package api

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/jamesog/example-http-api/database"
	"github.com/jamesog/example-http-api/database/mock"
)

func setup() (*API, database.Storage) {
	db, _ := mock.NewDB()
	api := NewService(db)
	return api, db
}

func TestGetUsers(t *testing.T) {
	api, _ := setup()
	tests := []struct {
		Name     string
		WantBody []byte
		WantCode int
	}{
		{
			Name:     "alice",
			WantBody: []byte(`[{"id":1,"user":"Alice Example"}]`),
			WantCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			api.GetUsersHTTP(w, req)
			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)

			if resp.StatusCode != tt.WantCode {
				t.Errorf("got code %d; want %d", resp.StatusCode, tt.WantCode)
			}
			if string(body) != string(tt.WantBody) {
				t.Errorf("got body %s; want %s", body, tt.WantBody)
			}
		})
	}
}
