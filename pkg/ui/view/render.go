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
	theme string = "default" //Can be configured if needed
	//basefile string = Dir() + "/internal/templates/" + theme + "/*/*.tmpl"
	basefile string = "/home/scaffolding/internal/templates/" + theme + "/*/*.tmpl"
)

// folder, err := filepath.Abs(filepath.Dir(os.Args[0]))
// if err != nil {
// 	log.Fatalln(err)
// }
// if _, err := templ.ParseGlob(filepath.Join(folder, "tmpl", "*.html")); err != nil {
// 	log.Fatalln(err)
// }

var RenderTemplates = &Template{
	templates: template.Must(template.ParseGlob(basefile)),
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	fmt.Println("basefile: ", basefile)
	fmt.Println(filepath.Abs("/home/scaffolding"))

	return t.templates.ExecuteTemplate(w, name, data)
}

// Add is our function that sums two integers
func Add(x, y int) (res int) {
	return x + y
}

// Subtract subtracts two integers
func Subtract(x, y int) (res int) {
	return x - y
}
