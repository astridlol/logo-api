package caching

import (
	"fmt"
	. "logo-api/structs"
	"os"
)

// exists - Returns whether the given file or directory
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetName - Get the file name for the logo
func GetName(logo Logo) string {
	return fmt.Sprintf("%s-%s-%s", logo.Emoji, logo.Color, logo.Platform)
}

// IsCached - Checks to see if a logo is already cached
func IsCached(logo Logo) bool {
	fileName := GetName(logo)
	exists, _ := exists(fmt.Sprintf("cache/%s.png", fileName))
	return exists
}
