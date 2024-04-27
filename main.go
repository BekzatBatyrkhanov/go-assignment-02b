package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	username = "bekzatbatyrkhanov"
	password = ""
	hostname = "localhost"
	port     = 5432
	db       = "postgres"
)

type Task struct {
	ID        int64
	Name      string
	Completed bool
}

func dbInsertTask(db *sql.DB, task *Task) error {
	_, err := db.Exec("INSERT INTO tasks(name) VALUES ($1)", task.Name)

	return err
}

func dbGetAllTasks(db *sql.DB) ([]*Task, error) {
	rows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Name, &task.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func dbUpdateTask(db *sql.DB, taskID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE tasks SET completed = true WHERE id = $1", taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func dbDeleteTask(db *sql.DB, taskID int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func main() {
	DSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", username, password, hostname, port, db)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	// INSERT a task
	err = dbInsertTask(db, &Task{
		Name:      "Damir",
		Completed: true,
	})
	if err != nil {
		fmt.Println("Error creating task:", err)
	}

	// UPDATE a task
	err = dbUpdateTask(db, 1)
	if err != nil {
		fmt.Println("Error updating task:", err)
	}

	// DELETE a task
	err = dbDeleteTask(db, 1)
	if err != nil {
		fmt.Println("Error deleting task:", err)
	}
}
