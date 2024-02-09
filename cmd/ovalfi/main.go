package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/odunlamizo/ovalfi/api/handler"
	"github.com/odunlamizo/ovalfi/api/middleware"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/login", middleware.ApplyContentJson(http.HandlerFunc(handler.Login))).Methods("POST")
	r.Handle("/task/create", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.CreateTask)))).Methods("POST")
	r.Handle("/task/{id}/update", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.UpdateTask)))).Methods("PUT")
	r.Handle("/task/{id}/complete", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.MarkTaskComplete)))).Methods("PUT")
	r.Handle("/task/{id}", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.GetTask)))).Methods("GET")
	r.Handle("/tasks", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.GetAllTasks)))).Methods("GET")
	r.Handle("/task/{id}/delete", middleware.ApplyContentJson(middleware.Authorize(http.HandlerFunc(handler.DeleteTask)))).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}
