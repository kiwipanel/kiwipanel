package helpers

import "path/filepath"

func GetAbsolutePath(relativePath string) (string, error) {
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}
