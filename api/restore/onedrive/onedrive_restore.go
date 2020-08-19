package onedrive

import (
	"fmt"
	http_local "go-onedrive-upload/graph/net/http"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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

// SimpleUpload allows you to provide the contents of a new file or update the
// contents of an existing file in a single API call. This method only supports
// files up to 4MB in size. For larger files use ResumableUpload().
//@userId will be extracted as sent from the restore input xml
//@bearerToken will be extracted as sent from the restore input xml
//@parentFolder and @file will be extracted from the file hierarchy the needs to be restored
func (rs *RestoreService) SimpleUploadToOriginalLoc(userId string, bearerToken string, conflictOption string, filePath string, file *os.File) (*http.Response, error) {
	uploadPath := fmt.Sprintf("/users/%s/drive/root:/%s:/content", userId, filePath)
	req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), file)
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

func (rs *RestoreService) SimpleUploadToAlternateLoc(altUserId string, bearerToken string, conflictOption string, filePath string, file *os.File) (*http.Response, error) {
	uploadPath := fmt.Sprintf("/users/%s/drive/root:/%s:/content", altUserId, filePath)
	req, err := rs.NewRequest("PUT", uploadPath, getSimpleUploadHeader(bearerToken), file)
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
