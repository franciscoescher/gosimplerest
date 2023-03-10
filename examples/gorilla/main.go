package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/examples"
	mysqlRepo "github.com/franciscoescher/gosimplerest/repository/mysql"
	"github.com/franciscoescher/gosimplerest/resource"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create router
	r := mux.NewRouter()

	// create routes for rest api
	resources := []resource.Resource{examples.UserResource, examples.RentEventResource, examples.VehicleResource}
	params := gosimplerest.AddHandlersBaseParams{Logger: logger, Resources: resources, Respository: mysqlRepo.NewRepository(db)}
	r = gosimplerest.AddGorillaMuxHandlers(r, params, examples.LoggingHandlerFunc)

	// iterates over routes and logs them
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, nil := route.GetMethods()
		if err != nil {
			return err
		}
		for _, method := range methods {
			logrus.WithFields(logrus.Fields{
				"method": method,
				"path":   tpl,
			}).Info("route registered")
		}
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
