package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	dsn = "host=pglsb user=rodrigo password=root dbname=events port=5432"
)

type webService struct{}

func NewWebServer() *webService {
	return &webService{}
}

type Event struct {
	gorm.Model
	TypeEvent string
	Context   string
	Price     float64
}

func (w webService) Serve() {
	e := echo.New()
	e.GET("/events", func(c echo.Context) error {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("failed to connect database")
		}

		var events []Event

		db.Find(&events)

		return c.JSON(http.StatusOK, &events)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
