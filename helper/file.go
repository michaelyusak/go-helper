package helper

import (
	"mime/multipart"
	"net/http"
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
		return false, "", err
	}
	defer file.Close()

	// Read first 512 bytes
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		return false, "", err
	}

	// Detect content type
	contentType := http.DetectContentType(buf)

	if !allowed[contentType] {
		return false, contentType, nil
	}

	return true, contentType, nil
}
