package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-sql-driver/mysql"
	null "gopkg.in/guregu/null.v3"
)

var testDB *sql.DB

var testResource = resource.Resource{
	Table:      "users_test",
	PrimaryKey: "uuid",
	Fields: map[string]resource.Field{
		"uuid":       {Validator: "uuid4"},
		"first_name": {Validator: "required,min=4,max=20"},
		"phone":      {Unsearchable: true},
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

func insertDBUserTestRow(data map[string]interface{}) error {
	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.Table),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	return err
}
