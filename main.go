package main

import (
	"fmt"
	"os"
	"restapi-echo-gorm/database"
	"strconv"
	"time"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/line/line-bot-sdk-go/linebot"
)

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	LineUserId string    `json:"line_user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Invoices   []Invoice `json:"invoices"`
}

type Invoice struct {
	Id                 int        `json:"id"`
	Name               string     `json:"name"`
	Price              int        `json:"price"`
	TargetYearAndMonth string     `json:"target_year_and_month"`
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
	month := c.Param("target_year_and_month")
	fmt.Println("TargetYearAndMonth", month)
	invoice := Invoice{Name: i.Name, Price: i.Price, UserID: i.UserID, TargetYearAndMonth: i.TargetYearAndMonth}
	database.DB.Create(&invoice)
	return c.JSON(http.StatusCreated, invoice)
}

func chargeInvoice(c echo.Context) error {
	// /:idからInvoiceを取得
	invoiceID := c.Param("id")
	invoice := Invoice{}
	database.DB.First(&invoice, invoiceID)
	t := time.Now()
	invoice.SentAt = &t
	database.DB.Save(&invoice)

	user := User{}
	database.DB.First(&user, invoice.UserID)

	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))

	if err != nil {
		fmt.Println(err)
	}

	bot.PushMessage(user.LineUserId, linebot.NewTextMessage(invoice.TargetYearAndMonth+"月分の請求\n請求内容:"+invoice.Name+"\n"+"金額: "+strconv.Itoa(invoice.Price)+"円")).Do()

	return c.JSON(http.StatusCreated, invoice)
}

func getUserWithInvoices(c echo.Context) error {
	users := []User{}
	database.DB.Preload("Invoices").Find(&users)
	return c.JSON(http.StatusOK, users)
}

func lineWebHook(c echo.Context) error {
	req := c.Request()
	bot, err := linebot.New(os.Getenv("LINE_CHANNEL_SECRET"), os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	events, err := bot.ParseRequest(req)
	if err != nil {
		fmt.Println(err)
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "振込先を追加" {
					user := User{}
					result := database.DB.Where(&User{LineUserId: event.Source.UserID}).First(&user)
					if result.RowsAffected > 0 {
						_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("既に登録済みです")).Do()
						if err != nil {
							// Do something when some bad happened
							fmt.Println(err)
						}
					} else {
						lineUserId := event.Source.UserID
						res, err := bot.GetProfile(lineUserId).Do()
						if err != nil {
							// Do something when some bad happened
							fmt.Println(err)
						}
						newUser := User{Name: res.DisplayName, LineUserId: lineUserId}
						database.DB.Create(&newUser)
						_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(res.DisplayName+"で登録しました")).Do()
						if err != nil {
							// Do something when some bad happened
							fmt.Println(err)
						}
					}
				} else {
					_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("そのメッセージには応答できません")).Do()
					if err != nil {
						// Do something when some bad happened
						fmt.Println(err)
					}
				}
			}
		}
	}
	return c.JSON(http.StatusOK, "OK")
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
	e.POST("/invoices/:id/charge", chargeInvoice)
	e.POST("/invoices", createInvoice)
	e.POST("/line_webhook", lineWebHook)

	e.Logger.Fatal(e.Start(":3000"))
}
