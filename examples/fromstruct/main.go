package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	mysqlRepo "github.com/franciscoescher/gosimplerest/repository/mysql"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"

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
	`domain_id` varchar(191) DEFAULT NULL,
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
	type Users struct {
		UUID      string    `json:"uuid" validate:"required" pk:"true"`
		FirstName string    `json:"first_name" validate:"required"`
		Phone     string    `json:"phone" validate:"required"`
		CreatedAt time.Time `json:"created_at" created_at:"true"`
		DeletedAt null.Time `json:"deleted_at" soft_delete:"true"`
	}
	err := user.FromStruct(Users{})
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("%+v\n", user)

	// create routes for rest api
	r := gin.Default()
	resources := []resource.Resource{user}
	params := gosimplerest.AddHandlersBaseParams{Logger: logger, Resources: resources, Respository: mysqlRepo.NewRepository(db)}
	gosimplerest.AddGinHandlers(r, params)

	logrus.Fatal(r.Run(":3333"))
}

func getDB() *sql.DB {
	c := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
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
