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

const (
	bearer_token = `eyJ0eXAiOiJKV1QiLCJub25jZSI6ImRrUzhpRTg1MTBZcWtGWkxmcUlhUUl4cjhXeDl2cGMxM2VSQ2NWTGNLSzAiLCJhbGciOiJSUzI1NiIsIng1dCI6Imh1Tjk1SXZQZmVocTM0R3pCRFoxR1hHaXJuTSIsImtpZCI6Imh1Tjk1SXZQZmVocTM0R3pCRFoxR1hHaXJuTSJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lNTliMzRjMC0wMDU1LTQ1N2QtYjhhZC1mYjM1NjRiODg4NTQvIiwiaWF0IjoxNTk3MzMwNzM5LCJuYmYiOjE1OTczMzA3MzksImV4cCI6MTU5NzMzNDYzOSwiYWlvIjoiRTJCZ1lHamI1UFkxSTI3KzlBTjdyaW4rZmlEN0R3QT0iLCJhcHBfZGlzcGxheW5hbWUiOiJ0ZXN0LW9uZWRyaXZlLWFwcCIsImFwcGlkIjoiODQyMjFjZjAtZGQwZC00NDkxLWEyMjItN2MxYjI1ZTQzZjIxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0LyIsIm9pZCI6IjUwMjRmZjM5LTRlNjMtNDBmZS1hNGU1LTdhODhlZmZiZWI2NCIsInJvbGVzIjpbIlNpdGVzLlJlYWRXcml0ZS5BbGwiLCJGaWxlcy5SZWFkV3JpdGUuQWxsIiwiRmlsZXMuUmVhZC5BbGwiXSwic3ViIjoiNTAyNGZmMzktNGU2My00MGZlLWE0ZTUtN2E4OGVmZmJlYjY0IiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IkFTIiwidGlkIjoiZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0IiwidXRpIjoibElNSlJZQUtNRTZZMjczbjNTUUVBUSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTkzMTY3Mjg3fQ.xduPL7PysDUHk0OxSfFsTzdtjHdWYulYP0CHEBrON4BVXghObpD5SabwIK7BxE9cCJ5kWyw9pAt1I_AAMj8F0JheGRGFetObgGChX0UIvdc4j21E2yc7m1qwDmvwoo4vC0Q2n2SjZ83clP5SD4tE2m3VpNaP_AedxEOWdXyB8zUGX6gUuReqGVfg9VgrgI-2RdxRZAm8bDIPCxvYKXWLjrJTqBLh06VzAskQ3M3YXiOYIS7F3yfl_DlRDQGWC0HLlsfc7vEuoRxW4kja7UcLnw5ONj4KZDxaneTDXt0-E5AEOr-bAND1Q-JFunCy2xyAYP7vx06iIMAu4FvgVH4Pfw`
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
		fmt.Println("ok")
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		return bodyString
	}
	return ""
}

func getSimpleUploadHeader(accessToken string) map[string]string {
	//As a work around for now, ultimately this will be recived as a part of restore xml
	if accessToken == "" {
		accessToken = bearer_token
	}
	bearerToken := fmt.Sprintf("bearer %s", accessToken)
	return map[string]string{
		"Content-Type":  "application/octet-stream",
		"Authorization": bearerToken,
	}
}
