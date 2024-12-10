package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/CollabTED/CollabTed-Backend/pkg/types"
)

func GenerateUniqueName(name string, existingUsers []types.UserWorkspace) string {
	highestSuffix := 0
	existingNameMap := make(map[string]bool)

	for _, userWorkspace := range existingUsers {
		existingNameMap[userWorkspace.Name] = true

		var currentSuffix int
		_, err := fmt.Sscanf(userWorkspace.Name, name+"(%d)", &currentSuffix)
		if err == nil && currentSuffix > highestSuffix {
			highestSuffix = currentSuffix
		}
	}

	if !existingNameMap[name] {
		return name
	}

	return fmt.Sprintf("%s(%d)", name, highestSuffix+1)
}

func RandomHexColor() string {
	// #nosec G404
	return fmt.Sprintf("#%06x", rand.Intn(0xFFFFFF+1))
}

// FetchAndEncodeImageToBase64 fetches an image from a URL and encodes it to Base64
func FetchAndEncodeImageToBase64(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	contentType := http.DetectContentType(buf.Bytes())
	if contentType != "image/png" && contentType != "image/jpeg" {
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}

	base64Image := base64.StdEncoding.EncodeToString(buf.Bytes())
	base64String := fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)

	return base64String, nil
}

func GenerateResetToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
