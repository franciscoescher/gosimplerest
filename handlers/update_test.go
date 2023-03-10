package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/repository/local"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestUpdateHandlerPatchOK(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "683143f8-e262-409c-b0a7-3df3ef296e2b",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}
	_, _ = base.Repository.Insert(&testResource, data)

	// Prepare the request
	dataUpdate := map[string]interface{}{
		"uuid":       data["uuid"],
		"first_name": "John",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(dataUpdate)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodPatch, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	data["first_name"] = dataUpdate["first_name"]
	data["phone"] = dataUpdate["phone"]
	dataDB, _ := base.Repository.Find(&testResource, data["uuid"])

	assert.Equal(t, data, dataDB)
}

func TestUpdateHandlerPutOK(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	var data = map[string]interface{}{
		"uuid":       "bec64968-663e-4e9e-9598-f3c139106bc4",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}
	_, _ = base.Repository.Insert(&testResource, data)

	// Prepare the request
	dataUpdate := map[string]interface{}{
		"uuid":       data["uuid"],
		"first_name": "John",
	}
	jsonData, err := json.Marshal(dataUpdate)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
	request, err := http.NewRequest(http.MethodPut, route, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	data["first_name"] = dataUpdate["first_name"]
	data["phone"] = dataUpdate["phone"]
	dataDB, _ := base.Repository.Find(&testResource, data["uuid"])

	assert.Equal(t, dataUpdate["first_name"], dataDB["first_name"])
	assert.Equal(t, dataUpdate["phone"], dataDB["phone"])

	// value of phone should be null
	assert.Nil(t, dataDB["phone"])
}

func TestUpdateNotFound(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":       "cec5c034-533a-4b97-8b86-6620bffb4242",
		"first_name": "John",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
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
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":  "1421b0f0-9b75-4220-b211-5356c93a8147",
		"phone": "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
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
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":       "78a63434-6e69-4a64-9138-01462f1c9721",
		"first_name": "A",
		"phone":      "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
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
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"phone": "+55 (11) 99999-9999",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
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

func TestUpdateImmutable(t *testing.T) {
	// Prepare the test
	base := &GetHandlerFuncParams{Resource: &testResource, Logger: logrus.New(), Repository: local.NewRepository(), Validate: validator.New()}

	// Prepare the request
	data := map[string]interface{}{
		"uuid":       "c660dfe1-0c4e-4c8a-b263-9c9cd79f3550",
		"created_at": time.Now(),
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := "/" + strcase.KebabCase(testResource.Table())
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
