package helper

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

	if contentType == "image/svg+xml" {
		if bytes.Contains(buf, []byte("<script")) {
			return false, contentType, nil
		}
	}

	if contentType == "text/plain; charset=utf-8" && LooksLikeCSV(buf) {
		contentType = "text/csv"
	}

	if !allowed[contentType] {
		return false, contentType, nil
	}

	return true, contentType, nil
}

func LooksLikeCSV(b []byte) bool {
	s := string(b)
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return false
	}
	return strings.Contains(lines[0], ",") || strings.Contains(lines[0], ";")
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

func NewCSVReader(file multipart.File) (*csv.Reader, error) {
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	line := string(buf[:n])
	firstLine := strings.SplitN(line, "\n", 2)[0]

	if strings.TrimSpace(firstLine) == "" {
		return nil, fmt.Errorf("empty or invalid CSV file")
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true

	switch {
	case strings.Count(firstLine, ";") > strings.Count(firstLine, ","):
		reader.Comma = ';'
	default:
		reader.Comma = ','
	}

	return reader, nil
}

func ReadCSVFromUpload(file multipart.File) ([]string, [][]string, error) {
	reader, err := NewCSVReader(file)
	if err != nil {
		return nil, nil, err
	}

	header := []string{}
	lines := [][]string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		if len(header) == 0 {
			header = append([]string{}, record...)
			continue
		}

		if len(record) != len(header) {
			return nil, nil, fmt.Errorf(
				"invalid column count: got %d, want %d",
				len(record), len(header),
			)
		}

		lines = append(lines, record)
	}

	if len(header) == 0 {
		return nil, nil, errors.New("csv has no header")
	}

	return header, lines, nil
}

func MultipartFromFilePath(fieldName, filePath string) (*multipart.FileHeader, multipart.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, nil, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		file.Close()
		return nil, nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		file.Close()
		return nil, nil, err
	}

	file.Close()
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())

	form, err := reader.ReadForm(stat.Size())
	if err != nil {
		return nil, nil, err
	}

	headers := form.File[fieldName]
	if len(headers) == 0 {
		return nil, nil, errors.New("no file in multipart form")
	}

	fh := headers[0]

	f, err := fh.Open()
	if err != nil {
		return nil, nil, err
	}

	return fh, f, nil
}
