package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/examples"
	mysqlRepo "github.com/franciscoescher/gosimplerest/repository/mysql"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/labstack/echo/v4"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create router
	r := echo.New()
	r.Use(echo.WrapMiddleware(examples.LoggingHandler))

	// create routes for rest api
	resources := []resource.Resource{examples.UserResource, examples.RentEventResource, examples.VehicleResource}
	params := gosimplerest.AddHandlersBaseParams{Logger: logger, Resources: resources, Respository: mysqlRepo.NewRepository(db)}
	gosimplerest.AddEchoHandlers(r, params)

	log.Fatal(r.Start(":3333"))
}

func getDB() *sql.DB {
	c := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_SCHEMA"),
		ParseTime:            true,
		Timeout:              5 * time.Second,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         5 * time.Second,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}
