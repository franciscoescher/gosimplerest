package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

var testDB *sql.DB

var testResource = resource.Resource{
	Table:      "users_test",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":       {},
		"first_name": {},
		"phone":      {},
		"created_at": {},
		"deleted_at": {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
}

var testBelongsResource = resource.Resource{
	Table:      "rent_events_test",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":          {},
		"user_id":       {},
		"starting_time": {},
		"hours":         {},
		"created_at":    {},
		"deleted_at":    {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
	BelongsToFields: []resource.BelongsTo{{Table: "users_test", Field: "user_id"}},
}

func TestCreateHandler(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	var data = map[string]interface{}{
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	var bodyJson map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &bodyJson)
	if err != nil {
		t.Fatal(err)
	}
	data["uuid"] = bodyJson["uuid"]

	assert.Equal(t, http.StatusOK, response.Code)

	sqlResult := base.DB.QueryRow(fmt.Sprintf(`SELECT uuid, first_name, phone FROM %s WHERE uuid = ? LIMIT 1`,
		testResource.Table), bodyJson["uuid"])
	dataDB := make([]string, 3)
	err = sqlResult.Scan(&dataDB[0], &dataDB[1], &dataDB[2])
	if err != nil {
		t.Fatal(err)
	}
	var dataMapDB = map[string]interface{}{
		"uuid":       dataDB[0],
		"first_name": dataDB[1],
		"phone":      dataDB[2],
	}

	assert.Equal(t, data, dataMapDB)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	testDB = getDB()
	testDB.Exec(fmt.Sprintf("DELETE FROM %s", testResource.Table))
	testDB.Exec(fmt.Sprintf("DELETE FROM %s", testBelongsResource.Table))
}

func shutdown() {
	testDB.Close()
}

func getDB() *sql.DB {
	c := mysql.Config{
		User:                 os.Getenv("DB_USER_TEST"),
		Passwd:               os.Getenv("DB_PASSWORD_TEST"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("DB_HOSTNAME_TEST"), os.Getenv("DB_PORT_TEST")),
		DBName:               os.Getenv("DB_SCHEMA_TEST"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return db
}
