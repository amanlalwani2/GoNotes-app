package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	"gofr.dev/pkg/gofr/request"
)

func TestIntegration(t *testing.T) {
	go main()
	time.Sleep(3 * time.Second)

	tests := []struct {
		desc       string
		method     string
		endpoint   string
		statusCode int
		body       []byte
	}{
		{"get notes", http.MethodGet, "notes", http.StatusOK, nil},
		{"post notes", http.MethodPost, "notes", http.StatusCreated, []byte(`{
    "note_id":23,
	"title":"testtitle123",
	"content":"testcontent123"
	}`)},
		{"update note", http.MethodPut, "notes/23", http.StatusOK, []byte(`{
		"note_id":23,
		"title":"testtitle",
		"content":"testcontent123"
		}`)},
		{"delete note", http.MethodDelete, "notes/23", http.StatusNoContent, nil},
	}

	for index, testcase := range tests {
		req, _ := request.NewMock(testcase.method, "http://localhost:3000/"+testcase.endpoint, bytes.NewBuffer(testcase.body))

		c := http.Client{}

		resp, err := c.Do(req)
		if err != nil {
			t.Errorf("TEST[%v] Failed.\t HTTP request encountered error:%v\n%s", index, err, testcase.desc)
			continue
		}

		if testcase.statusCode != resp.StatusCode {
			t.Errorf("TEST[%v] Failed.\t Expected %v \t Got %v\n%s", index, testcase.statusCode, resp.StatusCode, testcase.desc)
		}

		resp.Body.Close()
	}

}
