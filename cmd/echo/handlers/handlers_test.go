package handlers

import (
	"database/sql"
	"encoding/json"
	"final/cmd/echo/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	task := models.Task{
		ID:        1,
		Name:      "test",
		ListID:    1,
		Completed: false,
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	want := string(taskJson)

	_, err = db.Exec("INSERT INTO tasks (id, name,user_id) values (?,?,?)", task.ID, task.Name, task.ListID, task.Completed)

	if err != nil {
		log.Println(err)
	}

	e := echo.New()

	body := strings.NewReader(`{"name": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/lists/:id/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", 1)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := CreateList(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestGetTasks(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	task := models.Task{
		ID:        1,
		Name:      "da",
		ListID:    1,
		Completed: false,
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	want := string(taskJson)

	_, err = db.Exec("INSERT INTO tasks (id, name,user_id) values (?,?,?)", task.ID, task.Name, task.ListID, task.Completed)

	if err != nil {
		log.Println(err)
	}

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/lists/:id/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Request().Header.Set("name", "da")

	handler := GetLists(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestUpdateTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	task := models.Task{
		ID:        1,
		Name:      "test",
		ListID:    1,
		Completed: false,
	}

	_, err = db.Exec("INSERT INTO tasks (id, name,user_id) values (?,?,?)", task.ID, task.Name, task.ListID, task.Completed)

	if err != nil {
		log.Println(err)
	}

	taskJson, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}
	want := string(taskJson)

	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := UpdateTask(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestDeleteTask(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	task := models.Task{
		ID:        1,
		Name:      "test",
		ListID:    1,
		Completed: false,
	}

	_, err = db.Exec("INSERT INTO tasks (id, name,user_id) values (?,?,?)", task.ID, task.Name, task.ListID, task.Completed)

	if err != nil {
		log.Println(err)
	}

	want := string("{}")

	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := DeleteTask(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestCreateList(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	list := models.List{
		ID:     1,
		Name:   "test",
		UserID: 1,
	}

	_, err = db.Exec("INSERT INTO lists (id, name,user_id) values (?,?,?)", list.ID, list.Name, list.UserID)

	if err != nil {
		log.Println(err)
	}

	listJson, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	want := string(listJson)
	body := strings.NewReader(`{"name": "test"}`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/lists/:id/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", 1)

	handler := CreateList(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestGetLists(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	lists := []models.List{
		{
			ID:     1,
			Name:   "test",
			UserID: 1,
		},
	}

	_, err = db.Exec("INSERT INTO lists (id, name,user_id) values (?,?,?)", lists[0].ID, lists[0].Name, lists[0].UserID)

	if err != nil {
		log.Println(err)
	}

	listsJson, err := json.Marshal(lists)
	if err != nil {
		t.Fatal(err)
	}
	want := string(listsJson)

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/api/lists", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user_id", 1)

	handler := GetLists(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestDeleteList(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		log.Println(err)
	}

	migrate(db)

	lists := []models.List{
		{
			ID:     1,
			Name:   "test",
			UserID: 1,
		},
	}

	_, err = db.Exec("INSERT INTO lists (id, name,user_id) values (?,?,?)", lists[0].ID, lists[0].Name, lists[0].UserID)

	if err != nil {
		log.Println(err)
	}

	want := string("{}")

	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/lists/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	handler := DeleteList(db)(c)
	got := rec.Body.String()

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, want, got)
	}
}

func TestGetWeather(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/api/lists/:id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("lat", "lon")
	c.SetParamValues("41.99", "21")

	handler := GetWeather()(c)

	if assert.NoError(t, handler) {
		assert.Equal(t, http.StatusOK, rec.Code)
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
