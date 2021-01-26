package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}
type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task 1",
		Content: "blablabla",
	},
}

// create a task..
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "invalid task")
	}
	json.Unmarshal(reqbody, &newTask)
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	// http header type..
	w.Header().Set("Content-Type", "application/json")
	// if it's correct..
	w.WriteHeader(http.StatusCreated)
	// send added data..
	json.NewEncoder(w).Encode(newTask)
}

// find all tasks..
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// find one task..
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// parse integer..
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

// delete a task..
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been deleted", taskID)
		}
	}
}

// update a task..
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = task.ID
			tasks = append(tasks, updatedTask)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			fmt.Fprintf(w, "The task with ID %v has been updated", taskID)
		}
	}

}

// index route..
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "testeando api")

}

func main() {

	//routes..
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	// initialize server..
	fmt.Println("ejecutando servidor en puerto 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
