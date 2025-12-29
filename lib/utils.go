package lib

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetValidationMessage(err validator.ValidationErrors) string {
	if len(err) == 0 {
		return ""
	}
	for _, v := range err {
		switch v.Tag() {
		case "required":
			return fmt.Sprintf("%s harus diisi.", strings.ToLower(v.Field()))
		case "min":
			return fmt.Sprintf("%s minimal %s karakter", strings.ToLower(v.Field()), v.Param())
		case "max":
			return fmt.Sprintf("%s maksimal %s karakter.", strings.ToLower(v.Field()), v.Param())
		case "email":
			return fmt.Sprintf("%s bukan email yang valid.", strings.ToLower(v.Field()))
		case "oneof":
			return fmt.Sprintf("%s bukan pilihan yang valid.", strings.ToLower(v.Field()))
		default:
			fmt.Println(v.Tag())
		}
	}
	return ""
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash)
}

func IsPasswordMatch(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func OmitFields(columns ...string) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Omit(columns...)
	}
}

func ValidateFile(fileHeader *multipart.FileHeader) error {
	config := GetUploadConfig()

	// Check file size
	if fileHeader.Size > config.MaxFileSize {
		return fmt.Errorf("ukuran file %s melebihi batas maksimum %s",
			FormatFileSize(fileHeader.Size),
			FormatFileSize(config.MaxFileSize))
	}

	// Detect file type
	fileType := DetectFileType(fileHeader.Filename)

	// Check if file type is allowed
	if !IsAllowedFileType(fileHeader.Filename, fileType) {
		return fmt.Errorf("tipe file %s tidak diperbolehkan", filepath.Ext(fileHeader.Filename))
	}

	return nil
}

func ParseSize(sizeStr string) (int64, error) {
	sizeStr = strings.ToUpper(strings.TrimSpace(sizeStr))

	if strings.HasSuffix(sizeStr, "KB") {
		size, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "KB"), 10, 64)
		return size * 1024, err
	} else if strings.HasSuffix(sizeStr, "MB") {
		size, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "MB"), 10, 64)
		return size * 1024 * 1024, err
	} else if strings.HasSuffix(sizeStr, "GB") {
		size, err := strconv.ParseInt(strings.TrimSuffix(sizeStr, "GB"), 10, 64)
		return size * 1024 * 1024 * 1024, err
	}

	// Assume bytes if no suffix
	return strconv.ParseInt(sizeStr, 10, 64)
}

func ParseInt(str string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(str))
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func DetectFileType(filename string) FileType {
	ext := strings.ToLower(filepath.Ext(filename))

	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	documentExts := []string{".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt"}

	for _, imageExt := range imageExts {
		if ext == imageExt {
			return FileTypeImage
		}
	}

	for _, docExt := range documentExts {
		if ext == docExt {
			return FileTypeDocument
		}
	}

	return FileTypeOther
}

func IsAllowedFileType(filename string, fileType FileType) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))

	var allowedTypes []string

	switch fileType {
	case FileTypeImage:
		allowedTypesStr := "jpg,jpeg,png,gif,webp"
		allowedTypes = strings.Split(allowedTypesStr, ",")
	case FileTypeDocument:
		allowedTypesStr := "pdf,doc,docx,txt"
		allowedTypes = strings.Split(allowedTypesStr, ",")
	default:
		return false
	}

	for _, allowedType := range allowedTypes {
		if strings.TrimSpace(allowedType) == ext {
			return true
		}
	}

	return false
}

func IsUrl(urlInput string) bool {
	u, err := url.Parse(urlInput)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func GenerateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	filename := strings.TrimSuffix(originalName, ext)

	// Sanitize filename - hapus karakter yang tidak diinginkan
	filename = SanitizeFilename(filename)

	// Generate UUID untuk uniqueness
	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()

	return fmt.Sprintf("%s_%d_%s%s", filename, timestamp, uniqueID[:8], ext)
}

func SanitizeFilename(filename string) string {
	// Hapus karakter yang tidak diinginkan
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	sanitized := replacer.Replace(filename)

	// Limit panjang nama file
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}
