package utils

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
)

func GetContentType(uploadedFile *multipart.File) (string, error) {
	// get first 512 bytes
	buffer := make([]byte, 512)
	_, err := (*uploadedFile).Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer.
	_, err = (*uploadedFile).Seek(0, 0)
	if err != nil {
		return "", err
	}

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType, err
}
func SendJsonPayload(w http.ResponseWriter, payload map[string]interface{}, statusCode int) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	numBytes, writeError := w.Write(jsonResp)
	if writeError != nil {
		return 0, writeError
	}
	return numBytes, nil
}
