package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/franciscoescher/gosimplerest"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v3"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create router
	r := gin.Default()

	// create routes for rest api
	resources := []resource.Resource{{
		Table:      "users",
		PrimaryKey: "uuid",
		Fields: map[string]resource.Field{
			"uuid":        {Validator: "uuid4"},
			"first_name":  {},
			"last_name":   {},
			"phone":       {},
			"credit_card": {Unsearchable: true},
			"created_at":  {},
			"deleted_at":  {},
			"updated_at":  {},
		},
		SoftDeleteField: null.NewString("deleted_at", true),
		CreatedAtField:  null.NewString("created_at", true),
		UpdatedAtField:  null.NewString("updated_at", true),
	}}

	// This is the function that registers the routes!!!
	gosimplerest.AddGinHandlers(r, db, logger, nil, resources)

	log.Fatal(r.Run(":3333"))
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
