package main

import (
	"encoding/json"
	"fmt"
	"go-onedrive-upload/api/restore/onedrive"
	"go-onedrive-upload/fileutil"
	"log"
	"net/http"
)

const (
	user_id      = "ea49dc8a-cf17-49a4-8585-098155eaa5fa"
	bearer_token = `eyJ0eXAiOiJKV1QiLCJub25jZSI6IjVadEYxT2hRbE0tbEZweFJMMzh2a08zNE9tTm1uanlSOHZwZEZNN3FrV1UiLCJhbGciOiJSUzI1NiIsIng1dCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyIsImtpZCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lNTliMzRjMC0wMDU1LTQ1N2QtYjhhZC1mYjM1NjRiODg4NTQvIiwiaWF0IjoxNTk4MDE0NzE3LCJuYmYiOjE1OTgwMTQ3MTcsImV4cCI6MTU5ODAxODYxNywiYWlvIjoiRTJCZ1lOQTV3Qi9idCtWVXgrR0Z5NXN6L0JVOUFRPT0iLCJhcHBfZGlzcGxheW5hbWUiOiJ0ZXN0LW9uZWRyaXZlLWFwcCIsImFwcGlkIjoiODQyMjFjZjAtZGQwZC00NDkxLWEyMjItN2MxYjI1ZTQzZjIxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0LyIsIm9pZCI6IjUwMjRmZjM5LTRlNjMtNDBmZS1hNGU1LTdhODhlZmZiZWI2NCIsInJvbGVzIjpbIlNpdGVzLlJlYWRXcml0ZS5BbGwiLCJGaWxlcy5SZWFkV3JpdGUuQWxsIiwiRmlsZXMuUmVhZC5BbGwiXSwic3ViIjoiNTAyNGZmMzktNGU2My00MGZlLWE0ZTUtN2E4OGVmZmJlYjY0IiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IkFTIiwidGlkIjoiZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0IiwidXRpIjoicW5BWmhmZ1pUa2V2SmZDQlAtMVdBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTkzMTY3Mjg3fQ.nmXRR7xcqVGcTVre3BnU3yVfhDaRVJSsajpaPzoBjOTIypX6caSuB_w7MZXw3DJjxmghnvOAACAvweNCVWOl9MdhFbxRmz0llOqm-oSQr9yZL8XaOrUDR53jXN3cZa2Db7wfgGa8ngFPiCal9B4ms4cr5BJQJdfwwNGpNndjSTXM7iHut_JXucEj2s4K-ayxXIcUwdtk65Pe95qt7DDUaxVfYh-9zEWm6monFFk4BifYY2SF6VnTKSqR0iGNcjFlNOp2zvCnyaw76fBvUFL6Z90kf6VbKP7phj-TLN1SOtsaHQ9wO9DsRRNqmGJ7JDVA93nMIjNXUY4rV7YQs5q6Xw`
)

func main() {
	restoreOption := "orig"

	//Initialize the onedrive restore service
	restoreSrvc := onedrive.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path.
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource("./fileutil/Test")
	if err != nil {
		log.Fatalf("Failed to Load Files from source :%v", err)
	}

	//Call restore process based on alternate or original location
	if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload)
	}
}

//Restore to original location
func restore(restoreSrvc *onedrive.RestoreService, filesToRestore map[string]fileutil.FileInfo) {
	for filePath, fileInfo := range filesToRestore {
		respStr, err := restoreSrvc.SimpleUploadToOriginalLoc(user_id, bearer_token, "replace", filePath, fileInfo)
		if err != nil {
			log.Fatalf("Failed to Restore :%v", err)
			break
		}
		PrintResp(respStr)
	}
}

func PrintResp(resp interface{}) {
	switch resp.(type) {
	case *http.Response:
		respMap := make(map[string]interface{})
		rs := resp.(*http.Response)
		err := json.NewDecoder(rs.Body).Decode(&respMap)
		if err != nil {
			fmt.Printf("Error: --> %v", err)
		}
		if rs.Body != nil {
			defer rs.Body.Close()
		}
		fmt.Println(respMap)
		break
	case []*http.Response:
		for _, rs := range resp.([]*http.Response) {
			respMap := make(map[string]interface{})

			err := json.NewDecoder(rs.Body).Decode(&respMap)
			if err != nil {
				fmt.Printf("Error: --> %v", err)
			}
			if rs.Body != nil {
				defer rs.Body.Close()
			}
			fmt.Println(respMap)
			break
		}
	}
}

//Restore to Alternate location
func restoreToAltLoc(restoreSrvc *onedrive.RestoreService, filesToRestore map[string]fileutil.FileInfo) {
	rootFolder := fileutil.GetAlternateRootFolder()
	for filePath, fileItem := range filesToRestore {
		rootFilePath := fmt.Sprintf("%s/%s", rootFolder, filePath)
		respStr, err := restoreSrvc.SimpleUploadToAlternateLoc(user_id, bearer_token, "replace", rootFilePath, fileItem)
		if err != nil {
			log.Fatalf("Failed to Restore :%v", err)
			break
		}
		fmt.Println(respStr)
	}
}
