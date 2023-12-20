package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ToDo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type returnMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []ToDo `json:"data"`
}

type returnMessagesError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var todoId int = 1

var Todos []ToDo

func addNewTodoHandler(w http.ResponseWriter, r *http.Request) {
	// if not POST method, return 405

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var addNewTodo ToDo
	err := json.NewDecoder(r.Body).Decode(&addNewTodo)
	if err != nil {
		// return json error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		errMessage := returnMessagesError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	addNewTodo.ID = todoId
	todoId++

	Todos = append(Todos, addNewTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	successMessage := returnMessage{
		Status:  http.StatusCreated,
		Message: "Success",
	}

	fmt.Printf("%# v\n", Todos)
	json.NewEncoder(w).Encode(successMessage)
}

func listAllTodoHandler(w http.ResponseWriter, r *http.Request) {
	// if not GET method, return 405
	if r.Method != "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		errMessage := returnMessagesError{
			Status:  http.StatusMethodNotAllowed,
			Message: "Method not allowed",
		}

		json.NewEncoder(w).Encode(errMessage)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	returnMessage := returnMessage{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    Todos,
	}

	json.NewEncoder(w).Encode(returnMessage)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// if not DELETE method, return 405
	if r.Method != "DELETE" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		errMessage := returnMessagesError{
			Status:  http.StatusMethodNotAllowed,
			Message: "Method not allowed",
		}
		json.NewEncoder(w).Encode(errMessage)
		return
	}

	// get id from url
	// delete todo from todos
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		errMessage := returnMessagesError{
			Status:  http.StatusBadRequest,
			Message: "Missing id",
		}
		json.NewEncoder(w).Encode(errMessage)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, todo := range Todos {
		fmt.Printf("%# v\n", todo)
		fmt.Printf("%# v\n", intId)
		if todo.ID == intId {
			Todos = append(Todos[:i], Todos[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			successMessage := returnMessage{
				Status:  http.StatusOK,
				Message: "Success",
			}
			json.NewEncoder(w).Encode(successMessage)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	errMessage := returnMessagesError{
		Status:  http.StatusNotFound,
		Message: "Not found",
	}
	json.NewEncoder(w).Encode(errMessage)
}

func main() {
	// view html
	http.Handle("/", http.FileServer(http.Dir("./views/")))

	// endpoint
	http.HandleFunc("/api/todos", listAllTodoHandler)
	http.HandleFunc("/api/todos/add", addNewTodoHandler)
	http.HandleFunc("/api/todos/delete", deleteTodoHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
