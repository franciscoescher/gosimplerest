package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// Mysql Schema
/*
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
	`uuid` varchar(191) NOT NULL,
	`created_at` datetime(3) DEFAULT NULL,
	`deleted_at` datetime(3) DEFAULT NULL,
	`first_name` varchar(255) DEFAULT NULL,
	`phone` varchar(255) DEFAULT NULL,
	PRIMARY KEY (`uuid`)
);
*/

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// load resource from json file
	user := resource.Resource{}
	err := user.FromJSON("./examples/fromjson/user.json")
	if err != nil {
		logrus.Fatal(err)
	}

	// create routes for rest api
	r := gin.Default()
	gosimplerest.AddGinHandlers(
		r,
		db,
		logger,
		nil,
		[]resource.Resource{user})

	logrus.Fatal(r.Run(":3333"))
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
