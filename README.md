# rush-project

A simple Go app to manage tasks via CLI or web browser.
Demo link: https://patidar05priya.github.io/rush-project/ 
## Features
- CLI: Add, list, and delete tasks.
- Web UI: Add, view, and delete tasks.
- REST API: Supports web UI.
- Persistence: Tasks saved in `tasks.json`.

## Requirements
- Go (1.16+)
- Browser (for web UI)
- Git

## Install dependencies:

go mod init rush-project
go get github.com/gorilla/mux

## Available commands:

### Add a task:

```./rush-project add -task "Buy milk"```

#### Output:

```
Added task: 1 - Buy milk
```

### list tasks: 

```./rush-project list```

#### Output:

```
Tasks:
1 - Buy milk
2 - Do homework
```


### Delete a task:


```./rush-project delete -id 1```

#### Output:


```
Deleted task ID 1

```
Note: All tasks are saved automatically to tasks.json. 


## Web UI
Run the server:



```
go run main.go
```


```http://localhost:8080```

### Web Features:

Type a task and click Add Task to add it.

View all existing tasks.

Click Delete next to a task to remove it.

Technical Details:

The web UI communicates via REST API endpoints (/api/tasks).

All data is persisted in tasks.json.