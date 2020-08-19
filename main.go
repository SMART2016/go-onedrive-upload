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
	bearer_token = `eyJ0eXAiOiJKV1QiLCJub25jZSI6IjdPTXU5dnBjVjNldm82WVdBWkhWSlRiUXZobFA0RmNqc0JZcllReDVYaG8iLCJhbGciOiJSUzI1NiIsIng1dCI6Imh1Tjk1SXZQZmVocTM0R3pCRFoxR1hHaXJuTSIsImtpZCI6Imh1Tjk1SXZQZmVocTM0R3pCRFoxR1hHaXJuTSJ9.eyJhdWQiOiJodHRwczovL2dyYXBoLm1pY3Jvc29mdC5jb20iLCJpc3MiOiJodHRwczovL3N0cy53aW5kb3dzLm5ldC9lNTliMzRjMC0wMDU1LTQ1N2QtYjhhZC1mYjM1NjRiODg4NTQvIiwiaWF0IjoxNTk3NDI4NTY0LCJuYmYiOjE1OTc0Mjg1NjQsImV4cCI6MTU5NzQzMjQ2NCwiYWlvIjoiRTJCZ1lMaXYvR0J5NGgzdlp4Y1YxZTVNVTN3WkNBQT0iLCJhcHBfZGlzcGxheW5hbWUiOiJ0ZXN0LW9uZWRyaXZlLWFwcCIsImFwcGlkIjoiODQyMjFjZjAtZGQwZC00NDkxLWEyMjItN2MxYjI1ZTQzZjIxIiwiYXBwaWRhY3IiOiIxIiwiaWRwIjoiaHR0cHM6Ly9zdHMud2luZG93cy5uZXQvZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0LyIsIm9pZCI6IjUwMjRmZjM5LTRlNjMtNDBmZS1hNGU1LTdhODhlZmZiZWI2NCIsInJvbGVzIjpbIlNpdGVzLlJlYWRXcml0ZS5BbGwiLCJGaWxlcy5SZWFkV3JpdGUuQWxsIiwiRmlsZXMuUmVhZC5BbGwiXSwic3ViIjoiNTAyNGZmMzktNGU2My00MGZlLWE0ZTUtN2E4OGVmZmJlYjY0IiwidGVuYW50X3JlZ2lvbl9zY29wZSI6IkFTIiwidGlkIjoiZTU5YjM0YzAtMDA1NS00NTdkLWI4YWQtZmIzNTY0Yjg4ODU0IiwidXRpIjoiMUVnalBmeUFsa1NUTUF4cDFQSURBQSIsInZlciI6IjEuMCIsInhtc190Y2R0IjoxNTkzMTY3Mjg3fQ.5wydUjTtnv1BUe8Fr1sAmodm9zulCOCxbP8GzTCgZp_3PsOOCGil_6PTxsMqIw6KeoOzsO3u67QY7YBdcXh6akwqPuYm1Mwony1a9_FXXzdCvwkmLeT6XNY5QrvT0cWX4mfAL2IyxwMeAX61Jm2NRY7PTz6yla_9j1L3menuaodR33doxPIljWB16OKaPZX0ZL3teAcY3lnA_1zBleHTMhmwxSFuQF4-Iu6mVL23ia529alw5cgN6uw-Z5ktv9X-dSABVtrNDUEpNANjgUfGwZYpm7ZUtfgoxsodBgld6F_QHm7h1Xa7bqPmmwx8V_8bG7C0Rixo9gWfJVzNvtShWw`
)

func main() {
	restoreOption := "alt"

	//Initialize the onedrive restore service
	restoreSrvc := onedrive.GetRestoreService(http.DefaultClient)

	//Get the list of files that needs to be restore with the actual backed up path.
	fileInfoToUpload, err := fileutil.GetAllUploadItemsFrmSource("./fileutil")
	if err != nil {
		log.Fatalf("Failed to Restore :%v", err)
	}

	if restoreOption == "alt" {
		restoreToAltLoc(restoreSrvc, fileInfoToUpload)
	} else {
		restore(restoreSrvc, fileInfoToUpload)
	}
}

//Restore to original location
func restore(restoreSrvc *onedrive.RestoreService, filesToRestore map[string]fileutil.FileInfo) {
	for filePath, fileInfo := range filesToRestore {
		respStr, err := restoreSrvc.SimpleUploadToOriginalLoc(user_id, bearer_token, "rename", filePath, fileInfo)
		if err != nil {
			log.Fatalf("Failed to Restore :%v", err)
			break
		}
		fmt.Println(respStr)
	}
}

//Restore to Alternate location
func restoreToAltLoc(restoreSrvc *onedrive.RestoreService, filesToRestore map[string]fileutil.FileInfo) {
	rootFolder := fileutil.GetAlternateRootFolder()
	for filePath, fileItem := range filesToRestore {
		rootFilePath := fmt.Sprintf("%s/%s", rootFolder, filePath)
		respStr, err := restoreSrvc.SimpleUploadToAlternateLoc(user_id, bearer_token, "rename", rootFilePath, fileItem)
		if err != nil {
			log.Fatalf("Failed to Restore :%v", err)
			break
		}
		fmt.Println(respStr)
	}
}
