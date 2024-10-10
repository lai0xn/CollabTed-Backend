package utils

import (
	"fmt"
	"math/rand"

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
