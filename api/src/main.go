package main

import (
	"log"
	"net/http"

	"restful/src/database"
	"restful/src/tasks"
)

func main() {
	database.ConnectDB()
	defer database.DB.Close()

	http.HandleFunc("/create-task", tasks.CreateTask)
	http.HandleFunc("/get-task", tasks.GetTask)
	http.HandleFunc("/get-all-tasks", tasks.GetAllTasks)
	http.HandleFunc("/update-task", tasks.UpdateTask)
	http.HandleFunc("/delete-task", tasks.DeleteTask)

	log.Println("server started!")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println("failed to start server!")
	}
}
