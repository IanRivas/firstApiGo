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
	ID int `json:ID`
	Name string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks {
	{
		ID: 1,
		Name: "Task One",
		Content: "Some Content",
	},
}

func createTask(w http.ResponseWriter, r *http.Request){
	var newTask task
	reqBody,err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task")
	}
	//este lee el json y se lo asigna a new task
	json.Unmarshal(reqBody,&newTask)

	//creando el id incremental
	newTask.ID = len(tasks) + 1
	//agregando a tasks la nueva tarea
	tasks = append(tasks,newTask)
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	// le respondo al cliente la tarea que acabo de crear
	json.NewEncoder(w).Encode(newTask)

}

func getTask(w http.ResponseWriter, r *http.Request){
	//req.params
	vars := mux.Vars(r)
	//el id que tomo es un string , con esto lo convertimos a number
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for _,task := range tasks{
		if task.ID == taskID {
			w.Header().Set("Content-Type","application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}
	for i,task := range tasks {
		if task.ID == taskID {
			// agarro todos los elementos antes de la tarea a eliniar y las de despues las contateno asi elimino la tarea
			tasks = append(tasks[:i], tasks[i + 1:]...)
			fmt.Fprintf(w,"The task with ID %v has been remove succesfully", taskID)
		}
	}
}

func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	//este crea un json para mostrarlo
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil{
		fmt.Fprintf(w,"Invalid ID")
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i,t := range tasks {
		if t.ID == taskID {
			//elimino
			tasks = append(tasks[:i],tasks[i+1:]...)
			updatedTask.ID = taskID
			//agrego
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "The task with ID %v has been updated successfully", taskID)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w,"Welcome to my API")
}

func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/",indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}",getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}",deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}",updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000",router))
}