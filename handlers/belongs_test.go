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
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/testify/assert"
)

func TestBelongsToHandler(t *testing.T) {
	// Prepare the test
	base := &resource.Base{Resource: &testBelongsResource, Logger: logrus.New(), DB: testDB}

	t1 := time.Now()
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), t1.Hour(), t1.Minute(), t1.Second(), 0, time.UTC)
	data := map[string]interface{}{
		"uuid":       "c2a3495f-093f-47b2-9b66-0d1b18420a16",
		"first_name": "Fulano",
		"phone":      "+55 (11) 99999-9999",
		"deleted_at": t1,
		"created_at": t1.Add(-time.Hour * 24),
	}

	err := insertDBUserTestRow(data)
	if err != nil {
		t.Fatal(err)
	}

	data2 := map[string]interface{}{
		"uuid":          "48ee19da-8bf2-40ce-a84b-3bcc9a1c916f",
		"user_id":       data["uuid"],
		"starting_time": t1,
		"hours":         4,
		"deleted_at":    t1,
		"created_at":    t1.Add(-time.Hour * 24),
	}
	data3 := map[string]interface{}{
		"uuid":          "9c1a073d-e662-4a88-b7cc-105709cf427b",
		"user_id":       data["uuid"],
		"starting_time": t1,
		"hours":         4,
		"deleted_at":    t1,
		"created_at":    t1.Add(-time.Hour * 24),
	}
	err = insertDBEventTestRow(data2)
	if err != nil {
		t.Fatal(err)
	}

	err = insertDBEventTestRow(data3)
	if err != nil {
		t.Fatal(err)
	}

	// Make the request
	route := fmt.Sprintf("/%s/%s/%s", strcase.KebabCase(testBelongsResource.BelongsToFields[0].Table),
		data["uuid"], strcase.KebabCase(testBelongsResource.Table))
	request, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		t.Fatal(err)
	}
	request = GetRequestWithParams(request, map[string]string{"id": data["uuid"].(string)})
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBelongsToHandler(base, base.Resource.BelongsToFields[0]))
	handler.ServeHTTP(response, request)

	// Make assertions
	assert.Equal(t, http.StatusOK, response.Code)

	dataJson, err := json.Marshal([]map[string]interface{}{data2, data3})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(dataJson), strings.TrimSpace(response.Body.String()))
}

func insertDBEventTestRow(data map[string]interface{}) error {
	_, err := testDB.Exec(fmt.Sprintf("INSERT INTO %s (uuid, user_id, starting_time, hours, created_at, deleted_at) VALUES (?,?,?,?,?,?)",
		testBelongsResource.Table),
		data["uuid"], data["user_id"], data["starting_time"], data["hours"], data["created_at"], data["deleted_at"],
	)
	return err
}
