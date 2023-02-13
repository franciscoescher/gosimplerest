package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/examples"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create router
	r := fiber.New()
	r.Use(adaptor.HTTPMiddleware(examples.LoggingHandler))

	// create routes for rest api
	resources := []resource.Resource{examples.UserResource, examples.RentEventResource, examples.VehicleResource}
	gosimplerest.AddFiberHandlers(r, db, logger, nil, resources)

	log.Fatal(r.Listen(":3333"))
}

func getDB() *sql.DB {
	c := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT")),
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
