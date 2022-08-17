package main

import (
	"net/http"
	"strconv"
	"web-back-go/pkg/structs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

func (app *App) getAllTodo(c *gin.Context) {
	query := `SELECT id, active, done, content FROM todo WHERE active = 1`
	rows, err := app.db.Query(query)
	if err != nil {
		log.Errorf("Could not get all todo from db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}
	defer rows.Close()
	var todos []structs.Todo
	for rows.Next() {
		var todo structs.Todo
		if err := rows.Scan(&todo.ID, &todo.Active, &todo.Done, &todo.Content); err != nil {
			log.Errorf("Could not scan todo from db : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Unknown error, sorry",
			})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
	})
}

func (app *App) getTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Illegal todo ID : %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Illegal todo ID",
		})
		return
	}

	todo, err := getTodoFromDb(app.db, id)
	if err != nil {
		log.Errorf("Could not get todo from db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"todo": todo,
	})
}

func getTodoFromDb(db *sqlx.DB, id int) (structs.Todo, error) {
	query := `SELECT id, active, done, content FROM todo WHERE id = ?`
	row := db.QueryRow(query, id)
	var todo structs.Todo
	if err := row.Scan(&todo.ID, &todo.Active, &todo.Done, &todo.Content); err != nil {
		return structs.Todo{}, err
	}
	return todo, nil
}

func (app *App) createTodo(c *gin.Context) {
	var todo structs.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		log.Errorf("Could not bind JSON : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Could not bind JSON",
		})
		return
	}

	todoId, err := createTodoInDb(app.db, todo.Content)
	if err != nil {
		log.Errorf("Could not create todo in db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"todo_id": todoId,
	})
}

func createTodoInDb(db *sqlx.DB, content string) (int, error) {
	query := `INSERT INTO todo (content) VALUES (?)`
	result := db.MustExec(query, content)
	lastInsertId, _ := result.LastInsertId()
	return int(lastInsertId), nil
}

func (app *App) toggleTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Illegal todo ID : %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Illegal todo ID",
		})
		return
	}

	if err := toggleTodoInDb(app.db, id); err != nil {
		log.Errorf("Could not update todo in db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": id,
	})
}

func toggleTodoInDb(db *sqlx.DB, id int) error {
	query := `UPDATE todo SET done = NOT done WHERE id = ?`
	db.MustExec(query, id)
	return nil
}

func (app *App) updateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Illegal todo ID : %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Illegal todo ID",
		})
		return
	}

	var todo structs.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		log.Errorf("Could not bind JSON : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Could not bind JSON",
		})
		return
	}

	if err := updateTodoInDb(app.db, id, todo.Content); err != nil {
		log.Errorf("Could not update todo in db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": id,
	})
}

func updateTodoInDb(db *sqlx.DB, id int, content string) error {
	query := `UPDATE todo SET content = ? WHERE id = ?`
	db.MustExec(query, content, id)
	return nil
}

func (app *App) deleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("Illegal todo ID : %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Illegal todo ID",
		})
		return
	}

	if err := deleteTodoFromDb(app.db, id); err != nil {
		log.Errorf("Could not delete todo from db : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unknown error, sorry",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": id,
	})
}

func deleteTodoFromDb(db *sqlx.DB, id int) error {
	query := `UPDATE todo SET active = 0 WHERE id = ?`
	db.MustExec(query, id)
	return nil
}
