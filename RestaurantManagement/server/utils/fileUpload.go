package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func FileUpload(c*gin.Context) (string, error){
	// set file name 
	file, err := c.FormFile("file")
	if err != nil {
		return  "", err
	}

	// create upload folders
    err = os.MkdirAll("Uploads", os.ModePerm)
	if err != nil {
      return  "", err
	}

	// set filename and file path
	fileName := fmt.Sprint(time.Now().Unix())+ "_" + file.Filename
	filePath := filepath.Join("Uploads", fileName)
	
	// save the changes 
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return "", err
	}

	return filePath, err
}