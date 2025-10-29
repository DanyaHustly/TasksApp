package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type TaskBody struct {
	ID   int    `json:"id"`
	Task string `json:"task"`
}

var task string

var db *gorm.DB

func initBD() {
	dns := "host=localhost user=postgres password=task dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&TaskBody{})
}

func GetTask(c echo.Context) error {
	var tasks []TaskBody

	if err := db.Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Не может найти Task"})
	}

	return c.JSON(http.StatusOK, &tasks)
}

func PostTask(c echo.Context) error {
	var body TaskBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if err := db.Create(&body).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "не удалось создать task"})
	}

	return c.JSON(http.StatusOK, map[string]string{"Задачи на сегодня: ": task})
}

func PatchTask(c echo.Context) error {
	idPram := c.Param("id")
	id, err := strconv.Atoi(idPram)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var updatedTask TaskBody
	if err := c.Bind(&updatedTask); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	if err := db.Model(&TaskBody{}).Where("id = ?", id).Update("task", updatedTask.Task).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Неудалось обновить task"})
	}

	return c.JSON(http.StatusOK, map[string]string{"Задача обновлена: ": task})
}

func deleteTask(c echo.Context) error {
	idPram := c.Param("id")
	id, err := strconv.Atoi(idPram)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := db.Delete(&TaskBody{}, id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Не удалось удалить task"})
	}
	return c.JSON(http.StatusOK, task)
}

func main() {
	initBD()
	e := echo.New()
	e.GET("/tasks", GetTask)
	e.POST("/tasks", PostTask)
	e.PATCH("/tasks/:id", PatchTask)
	e.DELETE("/tasks/:id", deleteTask)
	e.Logger.Fatal(e.Start(":8080"))
}
