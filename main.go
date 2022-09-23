package main

import (
	"fmt"
	. "goplayground/utils"
	templates "goplayground/views"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type File struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Size        int64
	ContentType string
}

var db *gorm.DB

// handle index page
func handlerIndex(w http.ResponseWriter, _ *http.Request) {
	// filePath := path.Join("views", "index.html")
	_template, err := template.New("").Parse(templates.Index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data = map[string]interface{}{
		"title": "File Uploader",
		"name":  "Upload your file",
	}
	err = _template.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// get auth token from environment variable
func handlerSecret(w http.ResponseWriter, r *http.Request) {
	environmentVariable := os.Getenv("secret_code")
	if r.Method == http.MethodGet {
		// set secret environment variable as payload
		payload := make(map[string]interface{})
		payload["secret"] = environmentVariable

		_, err := SendJsonPayload(w, payload, http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// handle file upload
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		authToken := r.FormValue("auth")
		environmentVariable := os.Getenv("secret_code")

		// validate auth token
		if authToken != environmentVariable {
			payload := make(map[string]interface{})
			payload["message"] = "auth token is not valid"
			fmt.Println(authToken, environmentVariable)
			_, err := SendJsonPayload(w, payload, http.StatusForbidden)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// read uploaded file
		uploadedFile, fileHandler, err := r.FormFile("data")
		defer func() {
			err = uploadedFile.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		// get filename
		filename := fileHandler.Filename

		contentType, contentTypeErr := GetContentType(&uploadedFile)
		if contentTypeErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// validate uploaded file
		isValidImage := func() bool {
			validContentTypes := []string{"image/png", "image/jpg", "image/jpeg", "image/gif"}
			for _, v := range validContentTypes {
				if contentType == v {
					return true
				}
			}
			return false
		}()

		// check if uploaded file is an image
		if !isValidImage {
			payload := make(map[string]interface{})
			payload["message"] = "the uploaded file must be an image"

			_, err := SendJsonPayload(w, payload, http.StatusForbidden)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// validate image size
		fileSize := fileHandler.Size
		maxSize := 1024 * 1024 * 8
		if fileSize > int64(maxSize) {
			payload := make(map[string]interface{})
			payload["message"] = "the uploaded file must be an image"

			_, err := SendJsonPayload(w, payload, http.StatusForbidden)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// get current directory name
		dir, getCurrentDirErr := os.Getwd()
		if getCurrentDirErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get location where the file is uploaded to
		fileLocation := filepath.Join(dir, "files", filename)
		targetFile, targetFileErr := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		fmt.Println(fileLocation)
		defer func() {
			err = targetFile.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()

		// insert metadata to database
		file := File{Name: filename, Size: fileSize, ContentType: contentType}
		result := db.Create(&file)
		fmt.Println(result)

		if targetFileErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// copy file
		// note: I did not make temporary file, instead I upload it to filesystem directly
		// because why would someone save metadata to the database while the real data is not stored ?
		if _, err := io.Copy(targetFile, uploadedFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		payload := make(map[string]interface{})
		payload["message"] = "Method is not allowed"

		_, err := SendJsonPayload(w, payload, http.StatusMethodNotAllowed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "data.db")

	if err != nil {
		panic("Failed to open the SQLite database.")
	}
	defer db.Close()

	db.AutoMigrate(&File{})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// route handling
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/secret", handlerSecret)
	http.HandleFunc("/upload", handleUpload)

	// static files handling
	http.Handle("/static/",
		http.StripPrefix("/static",
			http.FileServer(http.Dir("assets"))),
	)

	// create http server
	var addr = "0.0.0.0:8000"

	// uncomment this line instead to run on goland
	//var addr = "localhost:8000"
	server := new(http.Server)
	server.Addr = addr
	fmt.Printf("Server started at %s\n", addr)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}
