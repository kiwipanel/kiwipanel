package view

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func getAbsolutePath(relativePath string) (string, error) {
	absolutePath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}

func Dir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get the current directory where template files are stored : ", err)
		return "."
	}
	return dir
}

// var indexHTML embed.FS

//https://forum.golangbridge.org/t/how-to-handle-paths-for-supporting-files-in-a-package-in-go/14651

var (
	theme    string = "default" //Can be configured if needed
	basefile string = Dir() + "/internal/templates/" + theme + "/*/*.tmpl"
	//basefile string = "internal/templates/" + theme + "/*/*.tmpl"
)

// Get the absolute path
var absoluteFilePath, _ = getAbsolutePath(basefile)

var RenderTemplates = &Template{
	templates: template.Must(template.ParseGlob(absoluteFilePath)),
	//templates: template.Must(template.ParseFS(indexHTML, absoluteFilePath)),
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
