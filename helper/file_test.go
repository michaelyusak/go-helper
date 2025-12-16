package helper_test

import (
	"testing"

	"github.com/michaelyusak/go-helper/helper"
	"github.com/stretchr/testify/assert"
)

func TestFileTypeAllowed(t *testing.T) {
	t.Run("Should allow CSV", func(t *testing.T) {
		fh, _, err := helper.MultipartFromFilePath("file", "/Users/michaelyusaktarigan/Documents/Financial Journal/Accounts/permata.csv")
		if err != nil {
			panic(err)
		}

		csvAllowed := map[string]bool{
			"text/csv": true,
		}

		allowed, _, err := helper.FileTypeAllowed(fh, csvAllowed)
		if err != nil {
			panic(err)
		}

		assert.True(t, allowed)
	})

	t.Run("Should return fileType text/csv", func(t *testing.T) {
		fh, _, err := helper.MultipartFromFilePath("file", "/Users/michaelyusaktarigan/Documents/Financial Journal/Accounts/permata.csv")
		if err != nil {
			panic(err)
		}

		csvAllowed := map[string]bool{
			"text/csv": true,
		}

		_, fileType, err := helper.FileTypeAllowed(fh, csvAllowed)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, fileType, "text/csv")
	})
}
