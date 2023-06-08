package utils

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

func IsImageFile(file *multipart.FileHeader) bool {
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)

	buf := make([]byte, 512)
	_, err = io.ReadFull(src, buf)
	if err != nil {
		return false
	}

	_, _ = src.Seek(0, io.SeekStart)

	mimeType := http.DetectContentType(buf)

	return isImageMIMEType(mimeType)
}

func isImageMIMEType(mimeType string) bool {
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err != nil {
		return false
	}

	return mediaType == "image/jpeg" || mediaType == "image/png"
}
