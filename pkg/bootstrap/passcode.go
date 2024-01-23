package bootstrap

import (
	"fmt"
	"os"
	"strings"
)

func readFileIfExists(filePath string) (string, error) {
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("File does not exist: %s", filePath)
	} else if err != nil {
		return "", err
	}

	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %s", err)
	}

	return string(content), nil
}

func ReadPasscode() (string, error) {
	filePath := "/home/state/passcode.txt"
	content, err := readFileIfExists(filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return strings.TrimSpace(content), nil
}
