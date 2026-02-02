package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// MediaResult resultado do processamento de midia
type MediaResult struct {
	Data     []byte
	MimeType string
	FileName string
}

// ProcessMedia processa midia de diferentes fontes: base64, URL ou form-data
// Retorna os bytes da midia, mimetype detectado e nome do arquivo (se disponivel)
func ProcessMedia(media string, providedMimeType string) (*MediaResult, error) {
	if media == "" {
		return nil, fmt.Errorf("media is empty")
	}

	// Verifica se é uma URL (http:// ou https://)
	if strings.HasPrefix(media, "http://") || strings.HasPrefix(media, "https://") {
		return downloadFromURL(media, providedMimeType)
	}

	// Verifica se é data URL (data:mime/type;base64,...)
	if strings.HasPrefix(media, "data:") {
		return decodeDataURL(media)
	}

	// Assume que é base64 puro
	return decodeBase64(media, providedMimeType)
}

// ProcessFormFile processa arquivo enviado via multipart/form-data
func ProcessFormFile(r *http.Request, fieldName string) (*MediaResult, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, fmt.Errorf("failed to get form file: %w", err)
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Detecta mimetype
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = detectMimeType(data, header.Filename)
	}

	return &MediaResult{
		Data:     data,
		MimeType: mimeType,
		FileName: header.Filename,
	}, nil
}

// downloadFromURL baixa midia de uma URL publica
func downloadFromURL(url string, providedMimeType string) (*MediaResult, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download from URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download from URL: status %d", resp.StatusCode)
	}

	// Limita o tamanho do download (50MB)
	limitReader := io.LimitReader(resp.Body, 50*1024*1024)
	data, err := io.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Detecta mimetype
	mimeType := providedMimeType
	if mimeType == "" {
		mimeType = resp.Header.Get("Content-Type")
		if mimeType == "" || mimeType == "application/octet-stream" {
			mimeType = detectMimeType(data, url)
		}
	}

	// Extrai nome do arquivo da URL
	fileName := filepath.Base(url)
	if idx := strings.Index(fileName, "?"); idx != -1 {
		fileName = fileName[:idx]
	}

	return &MediaResult{
		Data:     data,
		MimeType: mimeType,
		FileName: fileName,
	}, nil
}

// decodeDataURL decodifica data URL (data:mime/type;base64,...)
func decodeDataURL(dataURL string) (*MediaResult, error) {
	// Formato: data:mime/type;base64,XXXXX
	if !strings.HasPrefix(dataURL, "data:") {
		return nil, fmt.Errorf("invalid data URL format")
	}

	// Remove "data:"
	dataURL = dataURL[5:]

	// Separa mimetype do conteudo
	idx := strings.Index(dataURL, ",")
	if idx == -1 {
		return nil, fmt.Errorf("invalid data URL format: missing comma")
	}

	header := dataURL[:idx]
	content := dataURL[idx+1:]

	// Extrai mimetype
	mimeType := ""
	parts := strings.Split(header, ";")
	if len(parts) > 0 {
		mimeType = parts[0]
	}

	// Decodifica base64
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return &MediaResult{
		Data:     data,
		MimeType: mimeType,
	}, nil
}

// decodeBase64 decodifica base64 puro
func decodeBase64(b64 string, providedMimeType string) (*MediaResult, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	mimeType := providedMimeType
	if mimeType == "" {
		mimeType = detectMimeType(data, "")
	}

	return &MediaResult{
		Data:     data,
		MimeType: mimeType,
	}, nil
}

// detectMimeType detecta o mimetype baseado no conteudo ou extensao
func detectMimeType(data []byte, filename string) string {
	// Tenta detectar pelo conteudo
	mimeType := http.DetectContentType(data)

	// Se for generico, tenta pela extensao
	if mimeType == "application/octet-stream" && filename != "" {
		ext := filepath.Ext(filename)
		if ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				mimeType = mt
			}
		}
	}

	return mimeType
}

// IsURL verifica se a string é uma URL
func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// IsDataURL verifica se a string é uma data URL
func IsDataURL(s string) bool {
	return strings.HasPrefix(s, "data:")
}
