package provision

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// RenderTemplate loads template file, executes with data, writes to dest with 0640.
func RenderTemplate(templatesDir, tmplName string, data any, dest string) error {
	src := filepath.Join(templatesDir, tmplName)
	t, err := template.ParseFiles(src)
	if err != nil {
		return fmt.Errorf("parse template %s: %w", src, err)
	}
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	// write temp then rename to avoid partial writes
	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, b.Bytes(), 0640); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		return fmt.Errorf("rename tmp: %w", err)
	}
	return nil
}
