package view

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func Dir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get the current directory where template files are stored : ", err)
		return "."
	}
	return dir
}

var (
	theme string = "default" //Can be configured if needed
	//basefile string = Dir() + "/internal/templates/" + theme + "/*/*.tmpl"
	basefile string = "/home/scaffolding/internal/templates/" + theme + "/*/*.tmpl"
)

var RenderTemplates = &Template{
	templates: template.Must(template.ParseGlob(basefile)),
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	fmt.Println("basefile: ", basefile)

	return t.templates.ExecuteTemplate(w, name, data)
}
