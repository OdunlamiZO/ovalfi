package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/odunlamizo/ovalfi/internal/database"
	"github.com/odunlamizo/ovalfi/internal/model"
)

// creates a new task with status TODO and writes the task to the response
func CreateTask(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error reading request body"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	var task model.Task
	if err := json.Unmarshal(requestBody, &task); err != nil {
		errorMessage := model.ResponseMessage{Message: "Error decoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if task.Title == "" || task.Description == "" {
		errorMessage := model.ResponseMessage{Message: "Missing attribute(s)"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	task.Status = model.TODO
	task.User = r.Context().Value("user").(string)
	database.AddTask(&task)
	responseBody, err := json.Marshal(task)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// deletes task with the specified id param in the request
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	strId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strId)
	task := database.GetTask(id)
	if task == (model.Task{}) {
		errorMessage := model.ResponseMessage{Message: "Missing task"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if task.User != r.Context().Value("user").(string) {
		errorMessage := model.ResponseMessage{Message: "Access denied"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusForbidden)
		w.Write(responseBody)
		return
	}
	database.DeleteTask(task)
	responseMessage := model.ResponseMessage{Message: "Task deleted successfully"}
	responseBody, _ := json.Marshal(responseMessage)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// writes task with the specified id param in the request to the response
func GetTask(w http.ResponseWriter, r *http.Request) {
	strId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strId)
	task := database.GetTask(id)
	if task == (model.Task{}) {
		errorMessage := model.ResponseMessage{Message: "Missing task"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if task.User != r.Context().Value("user").(string) {
		errorMessage := model.ResponseMessage{Message: "Access denied"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusForbidden)
		w.Write(responseBody)
		return
	}
	responseBody, err := json.Marshal(task)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// writes a list of user's task to the response
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	user := database.GetUser(r.Context().Value("user").(string))
	tasks := database.GetAllTasks(user)
	if tasks == nil {
		errorMessage := model.ResponseMessage{Message: "Empty task list"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
		return
	}
	responseBody, err := json.Marshal(tasks)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// marks the task with the specified id param in the request as complete and writes the updated task to the response
func MarkTaskComplete(w http.ResponseWriter, r *http.Request) {
	strId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strId)
	task := database.GetTask(id)
	if task == (model.Task{}) {
		errorMessage := model.ResponseMessage{Message: "Missing task"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if task.User != r.Context().Value("user").(string) {
		errorMessage := model.ResponseMessage{Message: "Access denied"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusForbidden)
		w.Write(responseBody)
		return
	}
	task.Status = model.COMPLETED
	database.UpdateTask(&task)
	responseBody, err := json.Marshal(task)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// updates the task with the specified id and writes the updated task to the response
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error reading request body"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	var task model.Task
	if err := json.Unmarshal(requestBody, &task); err != nil {
		errorMessage := model.ResponseMessage{Message: "Error decoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if task.Status != model.TODO && task.Status != model.IN_PROGRESS && task.Status != model.COMPLETED {
		errorMessage := model.ResponseMessage{Message: "Invalid status value"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	strId := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(strId)
	databaseTask := database.GetTask(id)
	if databaseTask == (model.Task{}) {
		errorMessage := model.ResponseMessage{Message: "Missing task"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
		return
	}
	if databaseTask.User != r.Context().Value("user").(string) {
		errorMessage := model.ResponseMessage{Message: "Access denied"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusForbidden)
		w.Write(responseBody)
		return
	}
	task.Id = databaseTask.Id
	task.User = databaseTask.User
	database.UpdateTask(&task)
	responseBody, err := json.Marshal(task)
	if err != nil {
		errorMessage := model.ResponseMessage{Message: "Error encoding JSON"}
		responseBody, _ := json.Marshal(errorMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseBody)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
