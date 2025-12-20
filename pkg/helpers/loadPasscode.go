package helpers

import (
	"fmt"
	"os"
	"strings"
)

var (
	production = "/opt/kiwipanel/meta/passcode"
	dev        = "kiwipanel/meta/passcode"
)

func LoadGatePasscode() (string, error) {
	var filepath string
	if IsInstalled() {
		filepath = production
	} else {
		filepath = dev
	}

	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read passcode file (%s): %w", filepath, err)
	}
	kiwipanelPasscode := string(b)
	kiwipanelPasscode = strings.TrimSpace(kiwipanelPasscode)
	if kiwipanelPasscode == "" {
		return "", fmt.Errorf("passcode file is empty")
	}

	return kiwipanelPasscode, nil
}
