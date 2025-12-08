package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readFileIfExists(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var content string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "passcode:") {
			// If the line starts with "passcode:", extract the passcode
			content = strings.TrimPrefix(line, "passcode:")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}

func ReadPasscode() (string, error) {
	filePath := "/home/state/kiwipanel.conf"
	content, err := readFileIfExists(filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return strings.TrimSpace(content), nil
}
