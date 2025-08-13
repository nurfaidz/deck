package helpers

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
)

func FormatCurrency(amount uint) string {
	return fmt.Sprintf("Rp %d", amount)
}

func GenerateUniqueFilename(file *multipart.FileHeader) string {
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	return newFileName
}
