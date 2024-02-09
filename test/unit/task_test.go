package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/odunlamizo/ovalfi/api/handler"
	"github.com/odunlamizo/ovalfi/api/middleware"
	"github.com/odunlamizo/ovalfi/internal/database"
	"github.com/odunlamizo/ovalfi/internal/model"
	"github.com/odunlamizo/ovalfi/internal/util"
	"github.com/steinfletcher/apitest"
)

func TestCreateTask(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/task/create", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.CreateTask)))).Methods("POST")
	ts := httptest.NewServer(r)
	defer ts.Close()
	s, _ := util.CreateToken("Ttilaayo")
	t.Run("successful task creation", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Post("/task/create").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			JSON(`{
				"title": "Portfolio Website", 
				"description": "Create my portfolio website"
			}`).
			Expect(t).
			Body(`{
				"id": 1,
				"title": "Portfolio Website",
				"description": "Create my portfolio website",
				"status": "Todo"
			}`).
			Status(http.StatusOK).
			End()
	})

	t.Run("missing authorization token", func(t *testing.T) { // this applies for all other endpoints
		apitest.New().
			Handler(r).
			Post("/task/create").
			JSON(`{
				"title": "Portfolio Website", 
				"description": "Create my portfolio website"
			}`).
			Expect(t).
			Body(`{
				"message": "Authorization token not provided"
			}`).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("failed token validation", func(t *testing.T) { // this applies for all other endpoints
		apitest.New().
			Handler(r).
			Post("/task/create").
			Header("Authorization", fmt.Sprintf("Bearer %s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")).
			JSON(`{
				"title": "Portfolio Website", 
				"description": "Create my portfolio website"
			}`).
			Expect(t).
			Body(`{
				"message": "Error validating token"
			}`).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("missing attributes", func(t *testing.T) { // this applies for all other endpoints
		apitest.New().
			Handler(r).
			Post("/task/create").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			JSON(`{
				"description": "Create my portfolio website"
			}`).
			Expect(t).
			Body(`{
				"message": "Missing attribute(s)"
			}`).
			Status(http.StatusBadRequest).
			End()
	})
}

func TestUpdateTask(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/task/{id}/update", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.UpdateTask)))).Methods("PUT")
	ts := httptest.NewServer(r)
	defer ts.Close()
	s, _ := util.CreateToken("Ttilaayo")
	t.Run("successful task update", func(t *testing.T) {
		apitest.New().
			Handler(r).
			Put("/task/1/update").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			JSON(`{
				"status": "In Progress"
			}`).
			Expect(t).
			Body(`{
				"id": 1,
				"title": "Portfolio Website",
				"description": "Create my portfolio website",
				"status": "In Progress"
			}`).
			Status(http.StatusOK).
			End()
	})
	t.Run("missing task", func(t *testing.T) { // this applies to all other endpoints having this error
		apitest.New().
			Handler(r).
			Put("/task/2/update").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			JSON(`{
				"status": "In Progress"
			}`).
			Expect(t).
			Body(`{
				"message": "Missing task"
			}`).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("access to task denied", func(t *testing.T) { // this applies to all other endpoints having this error
		database.AddTask(&model.Task{
			Title:       "Update Resume",
			Description: "Add my recent work experience",
			Status:      model.TODO,
			User:        "Viktor",
		})
		apitest.New().
			Handler(r).
			Put("/task/2/update").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			JSON(`{
				"status": "In Progress"
			}`).
			Expect(t).
			Body(`{
				"message": "Access denied"
			}`).
			Status(http.StatusForbidden).
			End()
	})
}

func TestGetTask(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/task/{id}", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.GetTask)))).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	s, _ := util.CreateToken("Ttilaayo")
	apitest.New().
		Handler(r).
		Get("/task/1").
		Header("Authorization", fmt.Sprintf("Bearer %s", s)).
		Expect(t).
		Body(`{
			"id": 1,
			"title": "Portfolio Website",
			"description": "Create my portfolio website",
			"status": "In Progress"
		}`).
		Status(http.StatusOK).
		End()
}

func TestGetAllTask(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/tasks", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.GetAllTasks)))).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	t.Run("returns list of task", func(t *testing.T) {
		s, _ := util.CreateToken("Ttilaayo")
		apitest.New().
			Handler(r).
			Get("/tasks").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			Expect(t).
			Body(`[
			{
				"id": 1,
				"title": "Portfolio Website",
				"description": "Create my portfolio website",
				"status": "In Progress"
			}
		]`).
			Status(http.StatusOK).
			End()
	})
	t.Run("empty task list", func(t *testing.T) {
		s, _ := util.CreateToken("Mhiddey")
		apitest.New().
			Handler(r).
			Get("/tasks").
			Header("Authorization", fmt.Sprintf("Bearer %s", s)).
			Expect(t).
			Body(`{
			"message": "Empty task list"
		}`).
			Status(http.StatusOK).
			End()
	})
}

func TestMarkTaskComplete(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/task/{id}/complete", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.MarkTaskComplete)))).Methods("PUT")
	ts := httptest.NewServer(r)
	defer ts.Close()
	s, _ := util.CreateToken("Ttilaayo")
	apitest.New().
		Handler(r).
		Put("/task/1/complete").
		Header("Authorization", fmt.Sprintf("Bearer %s", s)).
		Expect(t).
		Body(`{
			"id": 1,
			"title": "Portfolio Website",
			"description": "Create my portfolio website",
			"status": "Completed"
		}`).
		Status(http.StatusOK).
		End()
}

func TestDeleteTask(t *testing.T) {
	r := mux.NewRouter()
	r.Handle("/task/{id}/delete", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.DeleteTask)))).Methods("DELETE")
	ts := httptest.NewServer(r)
	defer ts.Close()
	s, _ := util.CreateToken("Ttilaayo")
	apitest.New().
		Handler(r).
		Delete("/task/1/delete").
		Header("Authorization", fmt.Sprintf("Bearer %s", s)).
		Expect(t).
		Body(`{
			"message": "Task deleted successfully"
		}`).
		Status(http.StatusOK).
		End()
}
