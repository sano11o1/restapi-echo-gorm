package main

import (
	"restapi-echo-gorm/database"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	LineUserId string    `json:"line_user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Invoices   []Invoice
}

type Invoice struct {
	Id                 int        `json:"id"`
	Name               string     `json:"name"`
	Price              int        `json:price`
	TargetYearAndMonth string     `json:target_year_and_month`
	SentAt             *time.Time `json:"sent_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	UserID             uint       `json:"user_id"`
}

func getUsers(c echo.Context) error {
	users := []User{}
	database.DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}
func createUser(c echo.Context) error {
	name := c.Param("name")
	user := User{Name: name}
	if err := c.Bind(&user); err != nil {
		return err
	}
	database.DB.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

func getInvoices(c echo.Context) error {
	invoices := []Invoice{}
	database.DB.Find(&invoices)
	return c.JSON(http.StatusOK, invoices)
}

func createInvoice(c echo.Context) error {
	i := new(Invoice)
	if err := c.Bind(i); err != nil {
		return err
	}
	invoice := Invoice{Name: i.Name, Price: i.Price, UserID: i.UserID, TargetYearAndMonth: i.TargetYearAndMonth}
	database.DB.Create(&invoice)
	return c.JSON(http.StatusCreated, invoice)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	database.Connect()
	sqlDB, _ := database.DB.DB()
	defer sqlDB.Close()

	e.GET("/users", getUsers)
	e.POST("/users", createUser)
	e.GET("/invoices", getInvoices)
	e.POST("/invoices", createInvoice)

	e.Logger.Fatal(e.Start(":3000"))
}
