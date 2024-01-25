package view

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/kiwipanel/scaffolding/config"
	"github.com/labstack/echo/v4"
)

var app config.AppConfig

type mode struct {
	condition string
}

var test string

func LoadConfig(appConfig config.AppConfig) {
	app = appConfig
}

func Loadmode(getmode string) {
	fmt.Println("run load mode in render.go, getmode: ", getmode)
	test = getmode
}

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

//https://forum.golangbridge.org/t/how-to-handle-paths-for-supporting-files-in-a-package-in-go/14651

var (
	theme                string = "default" //Can be configured if needed
	basefile_development string = Dir() + "/internal/templates/" + theme + "/*/*.tmpl"
	basefile_production  string = "/home/scaffolding/internal/templates/" + theme + "/*/*.tmpl"
	runtime              string = "production"
)

func loadTemplateConditionally(condition string) string {
	if condition == "development" {
		return basefile_development
	}
	return basefile_production

}

var RenderTemplates = &Template{
	templates: template.Must(template.ParseGlob(loadTemplateConditionally(runtime))),
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
