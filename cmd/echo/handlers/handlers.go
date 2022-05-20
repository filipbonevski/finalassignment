package handlers

import (
	"database/sql"
	"final/cmd/echo/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type H map[string]interface{}

var loggedUser models.User

func AuthenticateUser(db *sql.DB, username string, password string) (bool, error) {
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err2 := rows.Scan(&user.ID, &user.Username, &user.Password)
		if err2 != nil {
			panic(err2)
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if username == user.Username && err == nil {
			loggedUser = user
			return true, nil
		}
	}

	return false, nil
}

func GetTasks(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		listID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, models.GetTasks(db, listID))
	}
}

func CreateTask(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		task := models.Task{}
		c.Bind(&task)
		listID, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return err
		}

		myTask, err := models.CreateTask(db, task.Name, listID)

		if err == nil {
			return c.JSON(http.StatusOK, myTask)
		} else {
			return err
		}
	}
}

func UpdateTask(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return err
		}
		task, err := models.UpdateTask(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, task)
		} else {
			return err
		}
	}
}

func DeleteTask(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		_, err := models.DeleteTask(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})

		} else {
			return err
		}

	}
}

func ExportTasks(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := models.ExportTasksCSV(db, loggedUser)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusOK, "successful operation")
	}
}

func GetLists(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, models.GetLists(db, loggedUser))
	}
}

func CreateList(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var list models.List
		c.Bind(&list)

		id, err := models.CreateList(db, list.Name, loggedUser)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"id":   id,
				"name": list.Name,
			})
		} else {
			return err
		}

	}
}

func DeleteList(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		err := models.DeleteList(db, id)

		if err == nil {
			return c.JSON(http.StatusOK, H{
				"deleted": id,
			})
		} else {
			return err
		}

	}
}

func GetWeather() echo.HandlerFunc {
	return func(c echo.Context) error {
		lat := c.Request().Header.Get("lat")
		lon := c.Request().Header.Get("lon")
		return c.JSON(http.StatusOK, models.GetWeather(lat, lon))
	}
}
