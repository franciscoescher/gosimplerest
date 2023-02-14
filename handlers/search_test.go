package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestSearchHandler(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testResource, Logger: logrus.New(), DB: testDB, Validate: validator.New()}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	data := map[string]interface{}{
		"uuid":       "4c017ccf-0749-4744-a5a6-9c92725411b9",
		"first_name": "Fulano Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}
	data2 := map[string]interface{}{
		"uuid":       "6b548c12-5cac-42e9-aaf1-465c31fafd63",
		"first_name": "Fulano Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1.Add(-time.Hour * 48),
		"created_at": t1.Add(-time.Hour * 72),
	}
	data3 := map[string]interface{}{
		"uuid":       "aed2737d-c105-4ee1-ab41-2341871ecd1a",
		"first_name": "John Search Test",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1.Add(-time.Hour * 80),
		"created_at": t1.Add(-time.Hour * 90),
	}

	err := insertDBUserTestRow(data)
	if err != nil {
		t.Fatal(err)
	}
	err = insertDBUserTestRow(data2)
	if err != nil {
		t.Fatal(err)
	}
	err = insertDBUserTestRow(data3)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s", strcase.KebabCase(testResource.GetTable()))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	q := request.URL.Query()
	q.Add("first_name", "Fulano Search Test")
	request.URL.RawQuery = q.Encode()
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(SearchHandler(base))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	dataJson, err := json.Marshal([]map[string]interface{}{data, data2})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(dataJson), strings.TrimSpace(response.Body.String()))
}

func insertDBUserTestRow(data map[string]interface{}) error {
	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, first_name, phone, created_at, deleted_at) VALUES (?,?,?,?,?)", testResource.GetTable()),
		data["uuid"], data["first_name"], data["phone"], data["created_at"], data["deleted_at"],
	)
	return err
}
