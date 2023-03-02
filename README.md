# Go Simple Rest

This package provides an out of the box implementation of a rest api router for go, with simple configuration of the resources (tables in the database).

Currently contains implementation using sql as database, logrus as logger and go validator v10 as param validator.

The api will create endpoints for each resource configuration provided to the Add<Router>Handlers functions.

Currently supported routers are gorilla mux, gin, chi, echo and fiber.

To resource configuration can be see in the `./resource/resource.go` file, in the `Resource` struct.

It creates the following routes (models are table names of the resources, converted to kebab case):
- GET /model/{id}
- GET /model
- POST /model
- PUT /model
- PATCH /model
- DELETE /model/{id}
- HEAD /model
- HEAD /model/{id}
  
The handlers created are standard http.HandlerFunc, so they can be used with any router.

Params from the url are passed to the handlers in the request context, using the `GetRequestWithParams` function.

## Simple usage

First, import the package:

`go get github.com/franciscoescher/gosimplerest`

Bellow, a simple example of how to use the package. For a more complete example, see the `./examples/gin` folder.

```
package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/franciscoescher/gosimplerest"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
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

var UserResource = gosimplerest.Resource{
	Name:      "users",
	PrimaryKey: "uuid",
	Fields: map[string]gosimplerest.Field{
		"uuid":       {},
		"first_name": {},
		"phone":      {},
		"created_at": {},
		"deleted_at": {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
}

func main() {
	logger := logrus.New()

	logger.Info("starting application")

	db := getDB()
	defer db.Close()

	// create routes for rest api
	r := gin.Default()
	gosimplerest.AddGinHandlers(
		r,
		db,
		logger,
		validator.New(),
		[]gosimplerest.Resource{UserResource})

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
```

## Disabling routes

Each resource can be configured with Ommit route flags, which can be used to disable a specific route for that resource

The available flags are:

```
// Ommmit<Route Type>Route are flags that omit the generation of the specific route from the router
OmitCreateRoute        bool `json:"omit_create_route"`
OmitRetrieveRoute      bool `json:"omit_retrieve_route"`
OmitUpdateRoute        bool `json:"omit_update_route"`
OmitPartialUpdateRoute bool `json:"omit_partial_update_route"`
OmitDeleteRoute        bool `json:"omit_delete_route"`
OmitSearchRoute        bool `json:"omit_search_route"`
OmitHeadRoutes         bool `json:"omit_head_routes"`
```

## Adding a new router type

To add a new router type, create a new file with the type of the router as name and that contains a function with the following signature:

`func Add<name>Handlers(router <new type>, *sql.DB, l *logrus.Logger, v *validator.Validate, resources []Resource`

This function should call the `AddHandlers` func, passing the AddRouteFunctions and AddParamFunc, which are a struct with functions that will add a route to the router, given a name and a handler (depending on the method), and a function that adds a parameter to a route url, respectively.
