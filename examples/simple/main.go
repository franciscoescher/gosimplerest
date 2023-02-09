package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

// Mysql Schema
/*
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
	`id` varchar(191) NOT NULL,
	`created_at` datetime(3) DEFAULT NULL,
	`name` varchar(255) DEFAULT NULL,
	`contact` varchar(255) DEFAULT NULL,
	PRIMARY KEY (`uuid`)
);
*/

var UserResource = gosimplerest.Resource{
	Table:      "users",
	PrimaryKey: "id",
	Fields: []gosimplerest.Field{
		{Name: "id"},
		{Name: "name"},
		{Name: "contact"},
		{Name: "created_at"},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
}

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create routes for rest api
	r := mux.NewRouter()
	r = gosimplerest.AddHandlers(
		db,
		logger,
		r,
		[]gosimplerest.Resource{UserResource},
		nil)

	// start server
	srv := &http.Server{
		Addr:         "127.0.0.1:3333",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      r,
	}
	logrus.Fatal(srv.ListenAndServe())
}

func getDB() *sql.DB {
	c := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT")),
		DBName:               os.Getenv("DB_SCHEMA"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}
