package validator

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"strings"
)

const (
	ProfileMinWidth  = 320
	ProfileMinHeight = 320
	PostMinDimension = 1080

	MaxProfileSize   = 5 << 20  // 5 MB
	MaxPostImageSize = 10 << 20 // 10 MB
	MaxCoverSize     = 10 << 20 // 10 MB 

	CoverMinWidth  = 1080
	CoverMinHeight = 540
)

func ValidateProfileImage(file *multipart.FileHeader) ([]byte, string, error) {
	if file == nil {
		return nil, "", fmt.Errorf("file is required")
	}
	if file.Size > MaxProfileSize {
		return nil, "", fmt.Errorf("profile image too large (max 5MB)")
	}
	return validateImage(file, ProfileMinWidth, ProfileMinHeight)
}

func ValidatePostImage(file *multipart.FileHeader) ([]byte, string, error) {
	if file == nil {
		return nil, "", fmt.Errorf("file is required")
	}
	if file.Size > MaxPostImageSize {
		return nil, "", fmt.Errorf("post image too large (max 10MB)")
	}
	return validateImage(file, PostMinDimension, PostMinDimension)
}

func ValidateChallengeCover(file *multipart.FileHeader) ([]byte, string, error) {
	if file == nil {
		return nil, "", fmt.Errorf("file is required")
	}
	if file.Size > MaxCoverSize {
		return nil, "", fmt.Errorf("cover image too large (max 10MB)")
	}
	return validateImage(file, CoverMinWidth, CoverMinHeight)
}

func validateImage(file *multipart.FileHeader, minW, minH int) ([]byte, string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	if ext != ".jpg" && ext != ".png" {
		return nil, "", fmt.Errorf("unsupported image type (only jpg/png)")
	}

	src, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("cannot open file")
	}
	defer src.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, "", fmt.Errorf("failed to read file")
	}
	data := buf.Bytes()

	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, "", fmt.Errorf("invalid or unsupported image")
	}

	mime := ""
	switch format {
	case "jpeg":
		mime = "image/jpeg"
	case "png":
		mime = "image/png"
	default:
		return nil, "", fmt.Errorf("unsupported image type (only jpg/png)")
	}

	if cfg.Width < minW || cfg.Height < minH {
		return nil, "", fmt.Errorf("image must be at least %dx%d pixels", minW, minH)
	}

	return data, mime, nil
}
