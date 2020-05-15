package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	pb "simpletodo/list-service/proto/list"

	"github.com/go-acme/lego/log"
	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
)

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type repository interface {
	InsertNewList(string, string) (int64, error)
	InsertNewTask(int64, string) (int64, error)
	GetList(int64) (*pb.List, error)
	GetTaskLists(int64) ([]*pb.Task, error)
	ToggleTask(int64) (*pb.Task, error)
	DeleteList(int64) error
}

type Repository struct {
	db *sql.DB
}

const (
	dbport     = "5432"
	dbuser     = "postgres"
	dbpassword = "vorto"
	dbname     = "todolist"
)

func openDbConnection() (db *sql.DB, err error) {
	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = "localhost"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	log.Printf("connection to postgres with connection string %s\n", psqlInfo)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error when opening db connection: %v\n", err)
		return
	}
	//log.Println("Successful connection to database on service startup")

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error when pinging db connection: %v\n", err)
		return
	}
	log.Printf("Successfully connected to %v\n", dbname)
	return
}

func (repo *Repository) InsertNewList(listName string, listUser string) (newListId int64, err error) {
	query := `INSERT INTO lists (list_name, list_user) VALUES ($1, $2) RETURNING list_id`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal("Fatal error in prepare statement of InsertNewList: %v\n", err)
		return
	}
	err = stmt.QueryRow(listName, listUser).Scan(&newListId)
	if err != nil {
		log.Fatal("Fatal error in execute statement of InsertNewList: %v\n", err)
		return
	}

	return
}

func (repo *Repository) InsertNewTask(listId int64, taskDescription string) (newTaskId int64, err error) {
	query := `INSERT INTO tasks (list_id, task_description) VALUES ($1, $2) RETURNING task_id`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal("Fatal error in prepare statement of InsertNewTask: %v\n", err)
	}
	err = stmt.QueryRow(listId, taskDescription).Scan(&newTaskId)
	if err != nil {
		log.Fatal("Fatal error in execute statement of InsertNewList: %v\n", err)
		return
	}

	return
}

func (repo *Repository) GetList(listId int64) (list *pb.List, err error) {
	query := `
		SELECT 
			list_id, list_name, list_user, list_created 
		FROM lists where list_id = $1
	`
	list = &pb.List{}
	// Todo: maybe come up with better method for handling conversion of time stamps
	var timestamp time.Time
	err = repo.db.QueryRow(query, listId).Scan(&list.ListId, &list.ListName, &list.ListUser, &timestamp)
	if err != nil {
		log.Fatal("Fatal error in execute statement of GetList for list id %v: %v\n", listId, err)
	}
	tsproto, _ := ptypes.TimestampProto(timestamp)
	list.ListCreated = tsproto
	return
}

func (repo *Repository) GetTaskLists(listId int64) (taskList []*pb.Task, err error) {
	query := `
		SELECT 
			task_id, task_description, task_created, task_complete 
		FROM tasks where list_id = $1
	`
	rows, err := repo.db.Query(query, listId)
	if err != nil {
		log.Fatal("Fatal error in execute statement of GetTaskLists: %v\n", err)
	}
	defer rows.Close()
	for rows.Next() {
		task := &pb.Task{}
		var timestamp time.Time
		err := rows.Scan(&task.TaskId, &task.TaskDescription, &timestamp, &task.TaskComplete)
		if err != nil {
			log.Fatal("Fatal error in row scan of GetTaskLists: %v\n", err)
		}
		tsproto, _ := ptypes.TimestampProto(timestamp)
		task.TaskCreated = tsproto
		taskList = append(taskList, task)
	}
	return
}

func (repo *Repository) ToggleTask(taskId int64) (task *pb.Task, err error) {
	query := `UPDATE tasks
			  SET task_complete = NOT task_complete 
			  WHERE task_id = $1
			  RETURNING task_complete`
	stmt, err := repo.db.Prepare(query)
	if err != nil {
		log.Fatal("Fatal error in preparing statement for ToggleTask: %v\n", err)
	}
	task = &pb.Task{TaskId: taskId}
	err = stmt.QueryRow(taskId).Scan(&task.TaskComplete)
	if err != nil {
		log.Fatal("Fatal error in row scan of ToggleTask: %v\n", err)
	}
	return
}

func (repo *Repository) DeleteList(listId int64) error {
	tx, err := repo.db.Begin()
	if err != nil {
		log.Fatal("Fatal error starting transaction to delete list %v: %v\n", listId, err)
	}
	if _, err := tx.Exec("DELETE FROM tasks WHERE list_id = $1", listId); err != nil {
		tx.Rollback()
		log.Fatal("Fatal error Deleting tasks associated with list %v: %v; rolling transaction back\n", listId, err)
	}
	if _, err := tx.Exec("DELETE FROM lists WHERE list_id = $1", listId); err != nil {
		tx.Rollback()
		log.Fatal("Fatal error Deleting List %v: %v; rolling transaction back\n", listId, err)
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		log.Fatal("Fatal error committing transaction to delete List %v: %v; rolling transaction back\n", listId, err)
	}
	log.Printf("Successfully deleted list and tasks for list %v\n", listId)
	return nil
}
