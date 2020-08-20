package onedrive

import (
	"encoding/json"
	"fmt"
	"go-onedrive-upload/fileutil"
	"net/http"
)

const (
	upload_session_path = "users/%s/drive/root:/%s:/createUploadSession"
)

func (rs *RestoreService) ressumableUpload(userId string, bearerToken string, conflictOption string, filePath string, fileInfo fileutil.FileInfo) (*http.Response, error) {
	return nil, nil
}

//Returns the restore session url for part file upload
func (rs *RestoreService) getUploadSession(userId string, bearerToken string, conflictOption string, filePath string, fileInfo fileutil.FileInfo) (*http.Response, error) {
	uploadSessionPath := fmt.Sprintf(upload_session_path, userId, filePath)
	body, err := getRessumableSessionBody(filePath, conflictOption)
	if err != nil {
		return nil, err
	}
	req, err := rs.NewRequest("POST", uploadSessionPath, getRessumableUploadSessionHeader(bearerToken), body)
	if err != nil {
		return nil, err
	}
	//Execute the request
	resp, err := rs.Do(req)
	if err != nil {
		//Need to return a generic object from onedrive upload instead of response directly
		return nil, err
	}
	return resp, nil
}

func getRessumableUploadSessionHeader(accessToken string) map[string]string {
	//As a work around for now, ultimately this will be recived as a part of restore xml
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": bearerToken,
	}
}

//Returns the expected body for creating file upload session to onedrive
func getRessumableSessionBody(filePath string, conflictOption string) (string, error) {
	bodyMap := map[string]string{"@microsoft.graph.conflictBehavior": conflictOption, "description": "", "name": filePath}
	jsonBody, err := json.Marshal(bodyMap)
	return string(jsonBody), err
}
