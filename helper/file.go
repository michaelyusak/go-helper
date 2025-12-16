package helper

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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

		// Other
		"text/plain; charset=utf-8": true,
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

	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	if !allowed[contentType] {
		return false, contentType, nil
	}

	if contentType == "image/svg+xml" {
		if bytes.Contains(buf, []byte("<script")) {
			return false, contentType, nil
		}
	}

	if contentType == "text/plain; charset=utf-8" {
		if LooksLikeCSV(buf) && allowed["text/csv"] {
			return true, "text/csv", nil
		}
	}

	return true, contentType, nil
}

func LooksLikeCSV(b []byte) bool {
	s := string(b)
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return false
	}
	return strings.Contains(lines[0], ",")
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

	return nil
}
