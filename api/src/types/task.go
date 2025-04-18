package types

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Completed   bool   `json:"completed"`
}

type GetTask struct {
	Title string `json:"title"`
}

type UpdateTask struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Completed   bool   `json:"completed"`
}

type CreateTaskResponse struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type GetTaskResponse struct {
	Message string `json:"description"`
}

type UpdateTaskResponse struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Completed   bool   `json:"completed"`
	Message     string `json:"message"`
}

type DeleteTaskResponse struct {
	Message string `json:"message"`
}
