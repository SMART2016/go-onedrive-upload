package onedrive

import (
	"fmt"
	"go-onedrive-upload/fileutil"
	http_local "go-onedrive-upload/graph/net/http"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	simple_upload_path = "/users/%s/drive/root:/%s:/content"
)

func GetRestoreService(c *http.Client) *RestoreService {
	return &RestoreService{
		http_local.NewOneDriveClient(c, false),
	}
}

// ItemService manages the communication with Item related API endpoints
type RestoreService struct {
	*http_local.OneDrive
}

// SimpleUploadToOriginalLoc allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 4MB in size. For larger files use ResumableUpload().
//@userId will be extracted as sent from the restore input xml
//@bearerToken will be extracted as sent from the restore input xml
//@filePath will be extracted from the file hierarchy the needs to be restored
//@fileInfo it is the file info struct that contains the actual file reference and the size_type
func (rs *RestoreService) SimpleUploadToOriginalLoc(userId string, bearerToken string, conflictOption string, filePath string, fileInfo fileutil.FileInfo) (*http.Response, error) {
	if fileInfo.SizeType == fileutil.SIZE_TYPE_LARGE {
		//For Large file type use resummable onedrive upload API
		fmt.Printf("\nInside Large File Processing: %s", filePath)
		return rs.ressumableUpload(userId, bearerToken, conflictOption, filePath, fileInfo)
	} else {
		uploadPath := fmt.Sprintf(simple_upload_path, userId, filePath)
		req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), fileInfo.FileData)
		if err != nil {
			return nil, err
		}

		//Handle query parameter for conflict resolution
		//The different values for @microsoft.graph.conflictBehavior= rename|replace|fail
		q := url.Values{}
		q.Add("@microsoft.graph.conflictBehavior", conflictOption)
		req.URL.RawQuery = q.Encode()

		//Execute the request
		resp, err := rs.Do(req)
		if err != nil {
			//Need to return a generic object from onedrive upload instead of response directly
			return nil, err
		}
		return resp, nil
	}

}

// SimpleUploadToAlternateLoc allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 4MB in size. For larger files use ResumableUpload().
//@userId will be extracted as sent from the restore input xml
//@filePath will be extracted from the file hierarchy the needs to be restored
//@fileInfo it is the file info struct that contains the actual file reference and the size_type
func (rs *RestoreService) SimpleUploadToAlternateLoc(altUserId string, bearerToken string, conflictOption string, filePath string, fileInfo fileutil.FileInfo) (*http.Response, error) {
	if fileInfo.SizeType == fileutil.SIZE_TYPE_LARGE {
		//For Large file type use resummable onedrive upload API
		return rs.ressumableUpload(altUserId, bearerToken, conflictOption, filePath, fileInfo)
	} else {

		uploadPath := fmt.Sprintf(simple_upload_path, altUserId, filePath)
		req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), fileInfo.FileData)
		if err != nil {
			return nil, err
		}

		//Handle query parameter for conflict resolution
		//The different values for @microsoft.graph.conflictBehavior= rename|replace|fail
		q := url.Values{}
		q.Add("@microsoft.graph.conflictBehavior", conflictOption)
		req.URL.RawQuery = q.Encode()

		//Execute the request
		resp, err := rs.Do(req)
		if err != nil {
			//Need to return a generic object from onedrive upload instead of response directly
			return nil, err
		}
		return resp, nil
	}
}

//Get response as string
func readRespAsString(resp *http.Response) string {
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return bodyString
	}
	return ""
}

func getSimpleUploadHeader(accessToken string) map[string]string {
	//As a work around for now, ultimately this will be recived as a part of restore xml
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Type":  "application/octet-stream",
		"Authorization": bearerToken,
	}
}
