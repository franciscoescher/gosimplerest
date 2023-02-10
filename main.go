package gosimplerest

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var db *sql.DB

/*
HandleResources creates a REST API for the given models.
It creates the following routes (models are table names of the resources, converted to kebab case):

GET /model/{id}
POST /model
PUT /model
DELETE /model/{id}
GET /model

Also, for each belongs to relation, it creates the following routes:

GET /belongs-to/{id}/model

The handlers parameter is a function to wrap the handlers with, for example, authentication and logging
*/
