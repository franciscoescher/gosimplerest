package gosimplerest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

var db *sql.DB

var testResource = Resource{
	Table:      "users_test",
	PrimaryKey: "uuid",
	Fields: map[string]Field{
		"uuid":       {},
		"first_name": {},
		"phone":      {},
		"created_at": {},
		"deleted_at": {},
	},
	SoftDeleteField: null.NewString("deleted_at", true),
}

func TestGetHandler(t *testing.T) {
	// Prepare the test
	base := Base{Resource: &testResource, Logger: logrus.New(), DB: db}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "644159a8-0b21-4250-8184-9f06457435c8",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}

	_, err := db.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.Table),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	if err != nil {
		t.Fatal(err)
	}

	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": data["uuid"].(string)})

	response := httptest.NewRecorder()

	// Call the handleras
	handler := http.HandlerFunc(GetHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)
	var bodyJson map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &bodyJson)
	if err != nil {
		t.Fatal(err)
	}

	// Convert time fields
	bodyJson["created_at"], err = time.Parse(time.RFC3339, bodyJson["created_at"].(string))
	if err != nil {
		t.Fatal(err)
	}
	bodyJson["deleted_at"], err = time.Parse(time.RFC3339, bodyJson["deleted_at"].(string))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, data, bodyJson)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	db = getDB()
	db.Exec(fmt.Sprintf("DELETE FROM %s", testResource.Table))
}

func shutdown() {
	db.Close()
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
