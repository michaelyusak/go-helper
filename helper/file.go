package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	defaultAllowedType = map[string]bool{
		// Images
		"image/jpeg":    true, // .jpg, .jpeg
		"image/png":     true, // .png
		"image/gif":     true, // .gif
		"image/webp":    true, // .webp
		"image/bmp":     true, // .bmp
		"image/tiff":    true, // .tif, .tiff
		"image/svg+xml": true, // .svg (vector)

		// PDF
		"application/pdf": true,
	}
)

func FileTypeAllowed(fileHeader *multipart.FileHeader, allowed map[string]bool) (bool, string, error) {
	if len(allowed) == 0 {
		allowed = defaultAllowedType
	}

	file, err := fileHeader.Open()
	if err != nil {
		return false, "", fmt.Errorf("[go-helper][FileTypeAllowed][fileHeader.Open] Failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return false, "", fmt.Errorf("[go-helper][FileTypeAllowed][file.Read] Failed to read from file: %w", err)
	}

	// Detect content type
	contentType := http.DetectContentType(buf)

	if !allowed[contentType] {
		return false, contentType, nil
	}

	return true, contentType, nil
}

func CopySourceToFile(fileName string, source io.Reader) error {
	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("[go-helper][SaveFile][os.Create] Failed to create dest file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, source)
	if err != nil {
		return fmt.Errorf("[go-helper][SaveFile][io.Copy] Failed to copy source to file: %w", err)
	}

	return  nil
}
