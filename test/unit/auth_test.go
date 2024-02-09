package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/odunlamizo/ovalfi/api/handler"
	"github.com/odunlamizo/ovalfi/api/middleware"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath/jsonpath"
)

func TestLogin(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/login", middleware.ApplyContentJson(http.HandlerFunc(handler.Login))).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("successful login", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/login").
			Query("username", "Ttilaayo").
			Query("password", "ovalfi").
			Expect(t).
			Assert(func(resp *http.Response, req *http.Request) error {
				if err := jsonpath.Present("access_token", resp.Body); err != nil {
					return err
				}
				return nil
			}).
			Status(http.StatusOK).
			End()
	})
	t.Run("wrong password", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/login").
			Query("username", "Ttilaayo").
			Query("password", "ovafi").
			Expect(t).
			Body(`{
				"message": "Incorrect password"
			}`).
			Status(http.StatusUnauthorized).
			End()
	})
}
