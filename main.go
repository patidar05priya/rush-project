package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var tasks []Task
var idCounter = 1

const taskFile = "tasks.json"

func main() {
	// CLI mode if arguments are provided
	loadTasks()
	if len(os.Args) > 1 {
		handleCLI()
		return
	}

	// Web server mode (API + static UI)
	r := mux.NewRouter()
	r.HandleFunc("/api/tasks", listTasks).Methods("GET")
	r.HandleFunc("/api/tasks", addTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", deleteTask).Methods("DELETE")
	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	// Enable CORS for local testing (optional, for cross-origin requests)
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleCLI() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	addTask := addCmd.String("task", "", "Task name to add")
	deleteID := deleteCmd.Int("id", 0, "Task ID to delete")

	if len(os.Args) < 2 {
		printCLIUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTask == "" {
			fmt.Println("Error: -task is required")
			addCmd.Usage()
			os.Exit(1)
		}
		add(*addTask)
	case "list":
		listCmd.Parse(os.Args[2:])
		list()
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Println("Error: -id is required")
			deleteCmd.Usage()
			os.Exit(1)
		}
		delete(*deleteID)
	default:
		printCLIUsage()
		os.Exit(1)
	}
}

func printCLIUsage() {
	fmt.Println("Usage: todo-hybrid <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  add -task <name>    Add a new task")
	fmt.Println("  list                List all tasks")
	fmt.Println("  delete -id <id>     Delete a task by ID")
}

// Shared Task Logic
func add(taskName string) {
	task := Task{ID: idCounter, Name: taskName}
	tasks = append(tasks, task)
	idCounter++
	saveTasks()
	fmt.Printf("Added task: %d - %s\n", task.ID, task.Name)
}

func list() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	fmt.Println("Tasks:")
	for _, task := range tasks {
		fmt.Printf("%d - %s\n", task.ID, task.Name)
	}
}

func delete(id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			fmt.Printf("Deleted task ID %d\n", id)
			return
		}
	}

	fmt.Printf("Task ID %d not found\n", id)
}

// API Handlers
func listTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func addTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	task.ID = idCounter
	idCounter++
	tasks = append(tasks, task)
	saveTasks() // Save to tasks.json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	fmt.Println()
	saveTasks()
	http.Error(w, "Task not found", http.StatusNotFound)
}
func loadTasks() {
	file, err := ioutil.ReadFile(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{} // Initialize empty if file doesn't exist
			return
		}
		log.Fatal(err)
	}
	if len(file) > 0 {
		if err := json.Unmarshal(file, &tasks); err != nil {
			log.Fatal(err)
		}
	}
	// Update idCounter based on existing tasks
	if len(tasks) > 0 {
		maxID := 0
		for _, task := range tasks {
			if task.ID > maxID {
				maxID = task.ID
			}
		}
		idCounter = maxID + 1
	}
}
func saveTasks() {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(taskFile, data, 0644); err != nil {
		log.Fatal(err)
	}
}
