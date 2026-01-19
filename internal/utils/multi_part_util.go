package utils

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/imkarthi24/sf-backend/internal/model/models"
	"github.com/loop-kar/pixie/errs"
)

func ExtractFile(fileHeader *multipart.FileHeader) (*models.FileUpload, error) {

	if fileHeader == nil {
		return nil, nil
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, errs.NewXError(errs.INVALID_REQUEST, "Error opening offer letter file", err)
	}
	defer file.Close()

	// Read file into buffer
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errs.NewXError(errs.IO, "Failed to read offer letter file", err)
	}

	outFile := models.FileUpload{}

	// Create a new reader that can be reset
	outFile.Content = io.NopCloser(bytes.NewReader(fileBytes))
	outFile.Metadata = models.FileMetadata{
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		Header: map[string][]string{
			"Content-Type": {fileHeader.Header.Get("Content-Type")},
		},
	}

	return &outFile, nil
}
