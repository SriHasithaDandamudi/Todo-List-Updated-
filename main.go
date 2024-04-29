package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

type List struct {
	Object string
	Finish bool
}

type PageInfo struct {
	Title string
	Todos []List
}

func list(w http.ResponseWriter, r *http.Request) {
	data := PageInfo{
		Title: "ToDo List",
		Todos: todos,
	}
	tmpl.Execute(w, data)
}

func create(w http.ResponseWriter, r *http.Request) {
	// Parse form data to get the new todo task
	r.ParseForm()
	todoObject := r.Form.Get("todo")

	// Append the new task to the existing list of todos
	todos = append(todos, List{Object: todoObject, Finish: false})

	// Redirect back to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func update(w http.ResponseWriter, r *http.Request) {
	// Parse form data to get the index of the completed task
	r.ParseForm()
	index := r.Form.Get("index")

	// Convert index to integer
	i, err := strconv.Atoi(index)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	// Update the Finish status of the corresponding task
	if i >= 0 && i < len(todos) {
		todos[i].Finish = true
	} else {
		http.Error(w, "Index out of range", http.StatusBadRequest)
		return
	}

	// Redirect back to the homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func remove(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	selectedTasks := r.Form["todo[]"] // Change to "todo[]" to handle multiple selections

	// Create a map to keep track of selected task indices
	selectedIndices := make(map[int]bool)
	for _, taskIndexStr := range selectedTasks {
		taskIndex, err := strconv.Atoi(taskIndexStr)
		if err != nil {
			http.Error(w, "Invalid task index", http.StatusBadRequest)
			return
		}
		selectedIndices[taskIndex] = true
	}

	// Filter out selected tasks from the todos slice
	updatedTodos := []List{}
	for i, todo := range todos {
		if !selectedIndices[i] {
			updatedTodos = append(updatedTodos, todo)
		}
	}
	todos = updatedTodos

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var todos = []List{
	{Object: "Write scripts", Finish: true},
	{Object: "Shoot video", Finish: false},
	{Object: "Edit the video", Finish: false},
}

func main() {
	mux := http.NewServeMux()
	tmpl = template.Must(template.ParseFiles("templates/mark.gohtml"))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", list)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/update", update)
	mux.HandleFunc("/remove", remove)

	/*go func() {
		for {
			removeCompleted()
			time.Sleep(time.Hour)
		}
	}() */

	log.Fatal(http.ListenAndServe(":8080", mux))
}
