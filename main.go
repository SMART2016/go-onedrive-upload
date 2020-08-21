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
	bearer_token = `eyJ0eXAiOiJKV1QiLCJub25jZSI6Imc0YXlyVFYzOHNjeGJXMFBJQmV5TThOVlFPajc4Xy0ybHlaZGZXeERhSWsiLCJhbGciOiJSUzI1NiIsIng1dCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyIsImtpZCI6ImppYk5ia0ZTU2JteFBZck45Q0ZxUms0SzRndyJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lNTliMzRjMC0wMDU1LTQ1N2QtYjhhZC1mYjM1NjRiODg4NTQvIiwiaWF0IjoxNTk4MDEwNjI0LCJuYmYiOjE1OTgwMTA2MjQsImV4cCI6MTU5ODAxNDUyNCwiYWlvIjoiRTJCZ1lKaTk5NlBkcll6KzQwK2V4a3pkOFBQMklnQT0iLCJhcHBfZGlzcGxheW5hbWUiOiJ0ZXN0LW9uZWRyaXZlLWFwcCIsImFwcGlkIjoiODQyMjFjZjAtZGQwZC00NDkxLWEyMjItN2MxYjI1ZTQzZjIxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0LyIsIm9pZCI6IjUwMjRmZjM5LTRlNjMtNDBmZS1hNGU1LTdhODhlZmZiZWI2NCIsInJvbGVzIjpbIlNpdGVzLlJlYWRXcml0ZS5BbGwiLCJGaWxlcy5SZWFkV3JpdGUuQWxsIiwiRmlsZXMuUmVhZC5BbGwiXSwic3ViIjoiNTAyNGZmMzktNGU2My00MGZlLWE0ZTUtN2E4OGVmZmJlYjY0IiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IkFTIiwidGlkIjoiZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0IiwidXRpIjoiS0ZyeGFyUGZ6MGVEaHhaeDE3R1JBUSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTkzMTY3Mjg3fQ.x-53FEJKJkLQkjZROTUfawTDf4BuSLT4EuBkwCnGY78jucsWR1TlFc7qo7ubEQPZKbfo0-FuuYO23KbPNR_yQGngkhYoKjrZ_-kFdqJgHbAoBpPnvziLlsTnmf2ucsLng9sN2ueot2VYy6H0Pw-hdaS4muvcfmxEH5KHyzf0N8F1Kvz6Yum7bnsVqIIwN23w3-3dI3E7VrBlH5x4dJYTQyjI6YtHqU9dZWhvS2_HpPguJvewNOhEsEVVmS_3cTP1jT_NaTIBCDxhKZb1j-B9vQ-qOrSrfb-KFzmRvnLftjGhXw2JQ8o5oCtpGh6WrhF_Fj9KdsSDE6DwNLolxAzz5A`
)

func main() {
	restoreOption := "orig"

	//Initialize the onedrive restore service
	restoreSrvc := onedrive.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path.
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource("./fileutil")
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
		defer respStr.Body.Close()
		//fmt.Println(respStr)
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
