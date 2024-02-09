package database

import "github.com/odunlamizo/ovalfi/internal/model"

var (
	size  int          = 0
	tasks []model.Task // a slice to model a database
)

// adds the specified task to database
func AddTask(task *model.Task) {
	size = size + 1
	task.Id = size
	tasks = append(tasks, *task)
}

// removes the specified task from database
func DeleteTask(task model.Task) {
	for index, t := range tasks {
		if t.Id == task.Id {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}
}

// retrieves task with specified id from database
func GetTask(id int) model.Task {
	var task model.Task
	for _, t := range tasks {
		if t.Id == id {
			task = t
		}
	}
	return task
}

// retrieves tasks owned by the specified user from the database
func GetAllTasks(user model.User) []model.Task {
	var userTasks []model.Task
	for _, t := range tasks {
		if t.User == user.Username {
			userTasks = append(userTasks, t)
		}
	}
	return userTasks
}

// updates task
func UpdateTask(task *model.Task) {
	for index, t := range tasks {
		if t.Id == task.Id {
			if task.Title == "" {
				task.Title = t.Title
			}
			if task.Description == "" {
				task.Description = t.Description
			}
			tasks[index] = *task
			break
		}
	}
}
