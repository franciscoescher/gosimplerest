package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestUpdateHandlerPatchOK(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe3",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}

	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.Table),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the request
	data = map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe3",
		"first_name": "John",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPatch, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	sqlResult := base.DB.QueryRow(fmt.Sprintf(`SELECT uuid, first_name, phone FROM %s WHERE uuid = ? LIMIT 1`,
		testResource.Table), data["uuid"])
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

func TestUpdateHandlerPutOK(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe2",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}

	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.Table),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	if err != nil {
		t.Fatal(err)
	}

	// Prepare the request
	data = map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe2",
		"first_name": "John",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPut, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	sqlResult := base.DB.QueryRow(fmt.Sprintf(`SELECT uuid, first_name, phone FROM %s WHERE uuid = ? LIMIT 1`,
		testResource.Table), data["uuid"])
	dataDB := make([]null.String, 3)
	err = sqlResult.Scan(&dataDB[0], &dataDB[1], &dataDB[2])
	if err != nil {
		t.Fatal(err)
	}
	var dataMapDB = map[string]interface{}{
		"uuid":       dataDB[0].String,
		"first_name": dataDB[1].String,
	}

	assert.Equal(t, data, dataMapDB)

	// value of phone should be null
	assert.False(t, dataDB[2].Valid)
}

func TestUpdateNotFound(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe4",
		"first_name": "John",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPatch, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusNotFound, response.Code)
}

// Validates if all fields are validated for put (not only the ones present in the request)
func TestUpdatePutBadRequest(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":  "df51f94b-9061-4b07-9a23-c1d493804fe4",
		"phone": "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPut, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUpdateBadRequest(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":       "df51f94b-9061-4b07-9a23-c1d493804fe4",
		"first_name": "A",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPatch, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestUpdateNoPrimaryKey(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"phone": "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.Table))
	request, err := http.NewRequest(http.MethodPatch, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusBadRequest, response.Code)
}
