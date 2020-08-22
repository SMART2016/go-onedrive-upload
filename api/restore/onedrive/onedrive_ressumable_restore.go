package onedrive

import (
	"encoding/json"
	"fmt"
	"go-onedrive-upload/fileutil"
	"net/http"
)

const (
	upload_session_path = "/users/%s/drive/root:/%s:/createUploadSession"
	upload_url_key      = "uploadUrl"
)

func (rs *RestoreService) ressumableUpload(userId string, bearerToken string, conflictOption string, filePath string, fileInfo fileutil.FileInfo) ([]map[string]interface{}, error) {
	//1. Get ressumable upload session for the current file path
	uploadSessionData, err := rs.getUploadSession(userId, bearerToken, conflictOption, filePath)
	if err != nil {
		return nil, err
	}

	//2. Get the upload url returned as a response from the ressumable upload session above.
	uploadUrl := uploadSessionData[upload_url_key].(string)

	//3. Get the startOffset list for the file
	startOfsetLst, err := fileutil.GetFileOffsetStash(filePath)
	if err != nil {
		return nil, err
	}

	//4. Loop over the file start offset list to read files in chunk and upload in onedrive
	var uploadResp []map[string]interface{}
	lastChunkIndex := len(startOfsetLst) - 1
	var isLastChunk bool
	for i, sOffset := range startOfsetLst {
		if i == lastChunkIndex {
			isLastChunk = true
		}
		//4a. Get the bytes for the file based on the offset
		filePartInBytes, err := fileutil.GetFilePartInBytes(filePath, sOffset, isLastChunk)
		if err != nil {
			return nil, err
		}
		fmt.Printf("\nUploading Part --> %d --> offset: %d", i, sOffset)
		//3b. make a call to the upload url with the file part based on the offset.
		resp, err := rs.uploadFilePart(uploadUrl, filePath, bearerToken, filePartInBytes, sOffset, isLastChunk)

		if err != nil {
			return nil, err
		}
		respMap := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&respMap)
		if err != nil {
			fmt.Println(err)
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		//fmt.Printf("%+v, status code: %s", respMap, resp.Status)
		uploadResp = append(uploadResp, respMap)
	}
	return uploadResp, nil
}

//Returns the restore session url for part file upload
func (rs *RestoreService) getUploadSession(userId string, bearerToken string, conflictOption string, filePath string) (map[string]interface{}, error) {
	uploadSessionPath := fmt.Sprintf(upload_session_path, userId, filePath)
	uploadSessionData := make(map[string]interface{})
	//Get the body for ressumable upload session call.
	body, err := getRessumableSessionBody(filePath, conflictOption)
	if err != nil {
		return nil, err
	}

	//Create request instance
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

	//convert http.Response to map
	err = json.NewDecoder(resp.Body).Decode(&uploadSessionData)
	if err != nil {
		return nil, err
	}
	return uploadSessionData, nil
}

//Uploads the file part to Onedrive
func (rs *RestoreService) uploadFilePart(uploadUrl string, filePath string, bearerToken string, filePart []byte, startOffset int64, isLastPart bool) (*http.Response, error) {
	//This is required for Content-Range header key
	fileSizeInBytes, err := fileutil.GetFileSize(filePath)
	if err != nil {
		return nil, err
	}

	//Fetch Last chunklength -- will be needed in Content_length header
	lastChunkLength, err := fileutil.GetLatsChunkSizeInBytes(filePath)
	if err != nil {
		return nil, err
	}

	//Create upload part file request
	req, err := rs.NewRequest("PUT", uploadUrl, getRessumableUploadHeader(fileSizeInBytes, bearerToken, startOffset, isLastPart, lastChunkLength), filePart)
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

//Returns header for upload session API
func getRessumableUploadSessionHeader(accessToken string) map[string]string {
	//As a work around for now, ultimately this will be recived as a part of restore xml
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Type":  "application/json",
		"Authorization": bearerToken,
	}
}

//Returns headers for ressumable actual upload as file parts
func getRessumableUploadHeader(fileSizeInBytes int64, accessToken string, startOffset int64, isLastChunk bool, lastChunkSize int64) map[string]string {
	var cRange string
	var cLength string

	if isLastChunk {
		cRange = fmt.Sprintf("bytes %d-%d/%d", startOffset, fileSizeInBytes-2, fileSizeInBytes-1)
		cLength = fmt.Sprintf("%d", lastChunkSize)
	} else {
		cRange = fmt.Sprintf("bytes %d-%d/%d", startOffset, startOffset+fileutil.GetDefaultChunkSize()-1, fileSizeInBytes-1)
		cLength = fmt.Sprintf("%d", fileutil.GetDefaultChunkSize())
	}

	fmt.Printf("\nCLength: %s , cRange: %s\n", cLength, cRange)
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Length": cLength,
		"Content-Range":  cRange,
		"Authorization":  bearerToken,
	}
}

//Returns the expected body for creating file upload session to onedrive
func getRessumableSessionBody(filePath string, conflictOption string) (string, error) {
	bodyMap := map[string]string{"@microsoft.graph.conflictBehavior": conflictOption, "description": "", "name": filePath}
	jsonBody, err := json.Marshal(bodyMap)
	return string(jsonBody), err
}
