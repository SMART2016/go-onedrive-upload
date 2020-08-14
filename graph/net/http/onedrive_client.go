package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	baseURL                       = "https://graph.microsoft.com/v1.0"
	statusInsufficientStorage int = 507
)

// OneDrive is the entry point for the client. It manages the communication with
// Microsoft OneDrive Graph API
type OneDrive struct {
	Client  *http.Client
	BaseURL string
}

// NewOneDrive returns a new OneDrive client to enable you to communicate with
// the API
func NewOneDriveClient(c *http.Client, debug bool) *OneDrive {
	drive := OneDrive{
		Client:  c,
		BaseURL: baseURL,
	}
	return &drive
}

func createRequestBody(body interface{}) (io.ReadWriter, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

//Converts the input in byte buffer which is the expected type for file upload api.
//Expects input as a reference type os.File
func getRequestBodyForFileItem(file interface{}) (*bytes.Buffer, error) {
	if _, ok := file.(*os.File); !ok {
		return nil, fmt.Errorf("Invalid type expected type: *os.File")
	}
	body := &bytes.Buffer{}
	_, err := io.Copy(body, file.(*os.File))
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Generate request
func (od *OneDrive) NewRequest(method, uri string, requestHeaders map[string]string, body interface{}) (*http.Request, error) {
	reqBody, err := getRequestBodyForFileItem(body)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the file into Bytes  reason: %v", err)
	}
	req, err := http.NewRequest(method, od.BaseURL+uri, reqBody)
	if err != nil {
		return nil, err
	}
	//Adding default header
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", getUserAgent())

	//Adding application specific Headers
	if requestHeaders != nil {
		for header, value := range requestHeaders {
			req.Header.Set(header, value)
		}
	}

	return req, nil
}

//Execute request
func (od *OneDrive) Do(req *http.Request) (*http.Response, error) {
	resp, err := od.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode <= statusInsufficientStorage {
		newErr := new(Error)
		if err := json.NewDecoder(resp.Body).Decode(newErr); err != nil {
			return resp, err
		}
		return resp, newErr
	}
	return resp, err
}
