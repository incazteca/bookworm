package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type params struct {
	filePath string
	filter   string
}

func TestSuccessFileUpload(t *testing.T) {
	tServer := httptest.NewServer(http.HandlerFunc(fileUploadHandler))
	defer tServer.Close()

	var testParams = []params{
		{"sample_files/blank.txt", ""},
		{"sample_files/fox.txt", ""},
		{"sample_files/better_nate.txt", ""},
		{"sample_files/better_nate_x_90.txt", ""},
	}

	expCode := 200

	for _, param := range testParams {
		// Retrieve text of file as our expected response
		// Also not the most efficient manner of doing this
		// But these test files are small-ish
		fileContents, err := ioutil.ReadFile(param.filePath)

		if err != nil {
			t.Fatalf("Error reading from file")
		}

		// Remove newline from fileContents

		expText := strings.TrimSpace(string(fileContents))

		// Create request
		req, err := uploadFileRequest(tServer.URL, param)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Submit request
		res, err := http.DefaultClient.Do(req)
		defer res.Body.Close()

		if err != nil {
			t.Fatalf(err.Error())
		}

		// Check response code and text of file

		if expCode != res.StatusCode {
			t.Errorf("Expected status code of: %d Received %d", expCode, res.StatusCode)
		}

		var jsonRes SuccessRes
		bodyText, _ := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(bodyText, &jsonRes)

		if err != nil {
			t.Fatalf("Error in parsing JSON response: %s, body value: %v", err.Error(), string(bodyText))
		}

		if expText != jsonRes.FileText {
			t.Errorf("File body doesn't match, expected: %s, Received: %s", expText, jsonRes.FileText)
		}
	}
}

func TestFailFileUpload(t *testing.T) {
	tServer := httptest.NewServer(http.HandlerFunc(fileUploadHandler))
	defer tServer.Close()

	// Expect Payload To Large HTTP Response COde
	expCode := 413

	failure := FailRes{"File Upload is limited to 10MB. Please submit a smaller file"}
	expResp, err := json.Marshal(failure)

	if err != nil {
		t.Fatalf("Error unable to convert failure response to JSON")
	}

	// Create request
	req, err := uploadFileRequest(tServer.URL, params{"sample_files/better_nate_x_180.txt", ""})
	if err != nil {
		t.Fatalf("Error in creating request")
	}

	// Submit request
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()

	if err != nil {
		t.Fatalf("Error in submitting request")
	}

	// Check response

	if expCode != res.StatusCode {
		t.Errorf("Expected status code of: %d Received %d", expCode, res.StatusCode)
	}

	bodyText, err := ioutil.ReadAll(res.Body)

	if string(expResp) == string(bodyText) {
		t.Errorf("Expected: %s, Received: %s", string(expResp), string(bodyText))
	}
}

// Create HTTP POST request to upload file with filter parameter
func uploadFileRequest(url string, p params) (*http.Request, error) {
	fh, err := os.Open(p.filePath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	field, err := writer.CreateFormFile("file", p.filePath)

	if err != nil {
		return nil, err
	}

	// Copy over params to form
	_, err = io.Copy(field, fh)
	_ = writer.WriteField("filter", p.filter)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
