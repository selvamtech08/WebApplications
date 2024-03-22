package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Todo struct
type ToDo struct {
	ID        string    `json:"id"`
	Item      string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Completed bool      `json:"status"`
}

var todos []ToDo

// show all todo items
func GetTodos(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, todos)
}

// Create new todo item
func AddTodo(ctx *gin.Context) {
	var newTodo ToDo
	newTodo.CreatedAt = time.Now() //current time
	newTodo.UpdatedAt = time.Now()
	if err := ctx.BindJSON(&newTodo); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "couldn't add the item", "msg": err})
		return
	}
	todos = append(todos, newTodo)
	ctx.IndentedJSON(http.StatusAccepted, gin.H{"message": "added succesffuly!!"})
}

// function to get todo item by id
func GetItemByID(id string) (*ToDo, error) {
	for i, _ := range todos {
		if todos[i].ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("incorrect item id")
}

// show single todo item using id(path parameter)
func GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := GetItemByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, item)
}

// update the todo status
func UpdateTodoStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	item, err := GetItemByID(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "couldn't update the todo item"})
		return
	}
	item.UpdatedAt = time.Now() //updated time
	item.Completed = !item.Completed
	ctx.IndentedJSON(http.StatusAccepted, item)
}

// remove the todo item
func RemoveTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	for idx, _ := range todos {
		if todos[idx].ID == id {
			todos = append(todos[0:idx], todos[idx+1:]...)
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"info": "todo item removed!"})
			return
		}
	}
	ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "couldn't find the todo item"})
}

// show incompleted todo items
func InCompletedTodo(ctx *gin.Context) {
	var inprogressToDo []ToDo
	for _, val := range todos {
		if !val.Completed {
			inprogressToDo = append(inprogressToDo, val)
		}
	}
	ctx.IndentedJSON(http.StatusOK, inprogressToDo)
}

func main() {
	router := gin.Default()
	router.GET("/todos", GetTodos)
	router.POST("/todos", AddTodo)
	router.GET("/todo/:id", GetTodo)
	router.PATCH("/todo/:id", UpdateTodoStatus)
	router.DELETE("/todo/:id", RemoveTodo)
	router.GET("/todos/incompleted", InCompletedTodo)
	router.Run("localhost:8070")
}
