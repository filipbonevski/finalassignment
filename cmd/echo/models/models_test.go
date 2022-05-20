package models

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestCreateTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)

	if err != nil {
		panic(err)
	}

	CreateTask(db, "prv", 1)
	tasks := GetTasks(db, 1)

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestGetTasks(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)

	if err != nil {
		panic(err)
	}

	CreateTask(db, "prv", 1)
	CreateTask(db, "dva", 1)
	tasks := GetTasks(db, 1)

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)

	if err != nil {
		panic(err)
	}

	task, _ := CreateTask(db, "prv", 1)
	task, _ = UpdateTask(db, task.ID)

	if task.Completed != true {
		t.Fatalf("expected true after update, got %t", task.Completed)
	}
}

func TestDeleteTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)

	if err != nil {
		panic(err)
	}

	CreateTask(db, "prv", 1)
	tasks := GetTasks(db, 1)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task after creation, got %d", len(tasks))
	}

	DeleteTask(db, 1)
	tasks = GetTasks(db, 1)

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks after deletion, got %d", len(tasks))
	}
}

func TestExportTasksCSV(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)

	if err != nil {
		panic(err)
	}

	user := User{
		ID:       1,
		Username: "user",
		Password: "pass",
	}

	CreateList(db, "nova", user)
	CreateTask(db, "prv", 1)
	CreateTask(db, "prv", 1)
	ExportTasksCSV(db, user)
	records := readCsvFile("result.csv")
	for _, strings := range records {
		for _, v := range strings {
			if v != "prv" {
				t.Fatalf("expected prv list, got %s", v)
			}

		}
	}
}

func TestCreateList(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)
	var user User
	if err != nil {
		panic(err)
	}

	CreateList(db, "nova", user)
	lists := GetLists(db, user)

	if len(lists) != 1 {
		t.Fatalf("expected 1 list, got %d", len(lists))
	}
}

func TestGetLists(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)
	var user User
	if err != nil {
		panic(err)
	}

	CreateList(db, "nova", user)
	CreateList(db, "vtor", user)
	lists := GetLists(db, user)

	if len(lists) != 2 {
		t.Fatalf("expected 2 lists, got %d", len(lists))
	}
}

func TestDeleteList(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	migrate(db)
	var user User
	if err != nil {
		panic(err)
	}

	id, _ := CreateList(db, "nova", user)
	lists := GetLists(db, user)
	if len(lists) != 1 {
		t.Errorf("expected 1 list after creation, got %d", len(lists))
	}

	DeleteList(db, int(id))
	lists = GetLists(db, user)

	if len(lists) != 0 {
		t.Fatalf("expected 0 lists after deletion, got %d", len(lists))
	}
}

func TestGetWeather(t *testing.T) {
	want := "Skopje"
	check := GetWeather("42", "21,41")
	got := check.City

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got: %v, want %v", got, want)
	}
}

func migrate(db *sql.DB) {
	sql := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		username VARCHAR NOT NULL UNIQUE,
		password VARCHAR NOT NULL
	);
    CREATE TABLE IF NOT EXISTS lists(
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        name VARCHAR NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
    );
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER NOT NULL,
		name VARCHAR NOT NULL,
		list_id INTEGER NOT NULL,
		completed INTEGER,
		PRIMARY KEY (id),
		FOREIGN KEY(list_id) REFERENCES lists(id)
	);
    `
	_, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
