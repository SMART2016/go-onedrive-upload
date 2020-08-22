package main

import (
	"fmt"
	"go-onedrive-upload/api/restore/onedrive"
	"go-onedrive-upload/fileutil"
	"log"
	"net/http"
)

const (
	user_id      = "ea49dc8a-cf17-49a4-8585-098155eaa5fa"
	bearer_token = `eyJ0eXAiOiJKV1QiLCJub25jZSI6InBkNmpOM0lhY1J5eU1KNUxQYjNzbHBHMTNzTlZaX3lMZnE2OGxsSzBHbGsiLCJhbGciOiJSUzI1NiIsIng1dCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyIsImtpZCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lNTliMzRjMC0wMDU1LTQ1N2QtYjhhZC1mYjM1NjRiODg4NTQvIiwiaWF0IjoxNTk4MTI1OTEyLCJuYmYiOjE1OTgxMjU5MTIsImV4cCI6MTU5ODEyOTgxMiwiYWlvIjoiRTJCZ1lOaFNrUFY4NFJObFRwSGxPKy8xZFNidEJnQT0iLCJhcHBfZGlzcGxheW5hbWUiOiJ0ZXN0LW9uZWRyaXZlLWFwcCIsImFwcGlkIjoiODQyMjFjZjAtZGQwZC00NDkxLWEyMjItN2MxYjI1ZTQzZjIxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0LyIsIm9pZCI6IjUwMjRmZjM5LTRlNjMtNDBmZS1hNGU1LTdhODhlZmZiZWI2NCIsInJvbGVzIjpbIlNpdGVzLlJlYWRXcml0ZS5BbGwiLCJGaWxlcy5SZWFkV3JpdGUuQWxsIiwiRmlsZXMuUmVhZC5BbGwiXSwic3ViIjoiNTAyNGZmMzktNGU2My00MGZlLWE0ZTUtN2E4OGVmZmJlYjY0IiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IkFTIiwidGlkIjoiZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0IiwidXRpIjoienN5VGlaR21mVXV4TjQ0ay11RkNBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTkzMTY3Mjg3fQ.oybnuwDoSbzaxLp3xEvoawHCCj9et15ISSza5KmC-VkJMz9zLEOhMW2Yk8Rn2HG_bgNl-c5yG6wFTfeER5eh_bnVtcupZyoYi02QMvS1NjOz2ZToTCQ9Ntg3pFH-yD_FXFzKm-YJJrl4SqxGRWkuFFx08VAvvekBJRnAAMBJSmdwx8QZkgT6uK9mFMlP9Mb06d-v9r2anODFEKV0-Ysn3U_DCXTvysljV6p2GSPjZBkZxwZ82k6STIktirdCm1olbwJOPrMcm67J_5Yy6c7zBmhLOMKGko7wCtS5t9Z2-LnyyeCJ4e3ZO7Jb4-HgR-Sd-C4q7WqFbeKnvICBBybVcA`
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
		resp, err := restoreSrvc.SimpleUploadToOriginalLoc(user_id, bearer_token, "replace", filePath, fileInfo)
		if err != nil {
			log.Fatalf("Failed to Restore :%v", err)
			break
		}
		PrintResp(resp)
	}
}

func PrintResp(resp interface{}) {
	switch resp.(type) {
	case map[string]interface{}:
		fmt.Printf("\n%+v\n", resp)
		break
	case []map[string]interface{}:
		for _, rs := range resp.([]map[string]interface{}) {
			fmt.Printf("\n%+v\n", rs)
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
