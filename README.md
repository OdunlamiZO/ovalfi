# Ovalfi

This is a project based assessment given to Odunlami Zacchaeus by Oval Labs Inc as part of the internship recruitment process.

## API Documentation

The following default users have been created: **Kenny**, **Ttilaayo**, **Viktor**, **Shegzy**, **Mhiddey** are their respective username, all having **ovalfi** as password.

The following are the available endpoint:

### /login

This is used for user authentication, returns a jwt signed access token.

```
POST /login?username=<username>&password=<password>
```

The returned signed jwt token is required to be sent through the **Authorization** header as **Bearer <access_token>** for every other request.

### /task/create

This is used to create a new task. Requires two attributes, title and description.

```
POST /task/create
{
    "title": "Portfolio Website",
    "description": "Create my portfolio website"
}
```

### /task/{id}/update

This is used to update a task. **id** is the task id of the task to be updated. You can update any of **title**, **description**, **status**, or even all. Status can only have values **Todo**, **In Progress**, and **Completed**.

```
PUT /task/{id}/update
{
    "status": "In Progress"
}
```

### /task/{id}/complete

This marks the task with the specified task **id** variable as completed.

```
PUT /task/{id}/complete
```

### /task/{id}

This is used to retrieve the task with the specified task **id** variable.

```
GET /task/{id}
```

### /tasks

This is used to retrieve a list of the currently authenticated user's tasks.

```
GET /tasks
```

### /task/{id}/delete

```
This is used to delete a task. **id** is the task id of the task to be deleted.
```

## Starting Application

Run the following command from the project directory to start the application

```
go run ./cmd/ovalfi/main.go
```

Or you can build a docker image and run using the following:

```
docker build --tag ovalfi .

docker run --publish 8080:8080 ovalfi
```

## Running Test

Run the following command from the project directory for test

```
go test ./test/unit -v
```