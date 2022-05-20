package models

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type List struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int
}

type ListCollection struct {
	Lists []List
}

type Task struct {
	ID        int    `json:"id"`
	Name      string `json:"text"`
	ListID    int    `json:"listId"`
	Completed bool   `json:"completed"`
}

type TaskCollection struct {
	Tasks []Task
}

type User struct {
	ID       int
	Username string
	Password string
}

type Weather struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather,omitempty"`

	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main,omitempty"`

	City string `json:"name,omitempty"`
}

type WeatherInfo struct {
	FormatedTemp string `json:"formatedTemp"`
	Description  string `json:"description"`
	City         string `json:"city"`
}

func GetTasks(db *sql.DB, listID int) []Task {
	rows, err := db.Query("SELECT * FROM tasks WHERE list_id = ?", listID)

	if err != nil {
		panic(err)
	}
	r := TaskCollection{}
	for rows.Next() {
		task := Task{}
		err2 := rows.Scan(&task.ID, &task.Name, &task.ListID, &task.Completed)

		if err2 != nil {
			panic(err2)
		}
		r.Tasks = append(r.Tasks, task)
	}
	if r.Tasks == nil {
		emptyTaskList := []Task{}
		return emptyTaskList
	}
	return r.Tasks
}

func CreateTask(db *sql.DB, name string, listID int) (Task, error) {
	result, err := db.Exec("INSERT INTO tasks(name, list_id,completed) VALUES(?,?,?);", name, listID, 0)

	if err != nil {
		panic(err)
	}

	taskID, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM tasks WHERE id = ?", taskID)
	if err != nil {
		panic(err)
	}

	task := Task{}
	for rows.Next() {
		err2 := rows.Scan(&task.ID, &task.Name, &task.ListID, &task.Completed)

		if err2 != nil {
			panic(err2)
		}
	}
	return task, nil
}

func UpdateTask(db *sql.DB, id int) (Task, error) {
	query := "UPDATE TASKS SET completed = NOT completed WHERE id = (?)"
	_, err := db.Exec(query, id)

	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM tasks WHERE id = ?", id)
	if err != nil {
		panic(err)
	}

	task := Task{}
	for rows.Next() {
		err2 := rows.Scan(&task.ID, &task.Name, &task.ListID, &task.Completed)

		if err2 != nil {
			panic(err2)
		}
	}
	return task, nil
}

func DeleteTask(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)

	if err != nil {
		panic(err)
	}

	return result.RowsAffected()
}

func ExportTasksCSV(db *sql.DB, loggedUser User) error {
	rows, err := db.Query("SELECT t.* FROM tasks AS t LEFT JOIN lists AS l ON l.id = t.list_id WHERE l.user_id = ?", loggedUser.ID)

	if err != nil {
		panic(err)
	}

	r := TaskCollection{}

	for rows.Next() {
		task := Task{}
		err2 := rows.Scan(&task.ID, &task.Name, &task.ListID, &task.Completed)

		if err2 != nil {
			panic(err2)
		}
		r.Tasks = append(r.Tasks, task)
	}

	file, err := os.Create("result.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var tasksReadyForImport []string
	for _, task := range r.Tasks {
		tasksReadyForImport = append(tasksReadyForImport, task.Name)
	}

	if err := writer.Write(tasksReadyForImport); err != nil {
		return err
	}

	return nil
}

func GetLists(db *sql.DB, loggedUser User) []List {
	rows, err := db.Query("SELECT * FROM lists WHERE user_id = ?", loggedUser.ID)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result := ListCollection{}
	for rows.Next() {
		list := List{}
		err2 := rows.Scan(&list.ID, &list.Name, &list.UserID)

		if err2 != nil {
			panic(err2)
		}
		result.Lists = append(result.Lists, list)
	}

	if result.Lists == nil {
		emptyList := []List{}
		return emptyList
	}
	return result.Lists
}

func CreateList(db *sql.DB, name string, loggedUser User) (int64, error) {
	result, err := db.Exec("INSERT INTO lists(name, user_id) VALUES(?,?);", name, loggedUser.ID)

	if err != nil {
		panic(err)
	}

	return result.LastInsertId()
}

func DeleteList(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM lists WHERE id = ?", id)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM tasks WHERE list_id = ?", id)

	if err != nil {
		panic(err)
	}

	return nil
}

func GetWeather(lat string, lon string) WeatherInfo {
	weather := Weather{}
	apiKey := "8376a04ba3fd6b44983c55f31b28d93a"

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&appid=%s", lat, lon, apiKey)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("weather api response not ok")
	}

	json.NewDecoder(resp.Body).Decode(&weather)

	weatherResponse := WeatherInfo{
		FormatedTemp: fmt.Sprintf("%f Celsius", weather.Main.Temp-273.15),
		Description:  weather.Weather[0].Description,
		City:         weather.City,
	}
	return weatherResponse
}
