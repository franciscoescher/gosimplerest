package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/examples"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-chi/chi"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create router
	r := chi.NewRouter()
	r.Use(examples.LoggingHandler)

	// create routes for rest api
	resources := []resource.Resource{examples.UserResource, examples.RentEventResource, examples.VehicleResource}
	r = gosimplerest.AddChiHandlers(r, db, logger, nil, resources)

	// iterates over routes and logs them
	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.WithFields(logrus.Fields{
			"method": method,
			"path":   route,
		}).Info("route registered")
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// start server
	srv := &http.Server{
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
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
