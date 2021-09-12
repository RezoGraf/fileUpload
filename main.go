package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// func clearConsole() {
// 	cmd := exec.Command("cmd", "/c", "cls")
// 	cmd.Stdout = os.Stdout
// 	cmd.Run()
// }

func scanFiles() []string {
	var fileNames []string
	files, err := ioutil.ReadDir("./tmp/")
	if err != nil {
		log.Fatal(err)
		return fileNames
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
		fmt.Println(file.Name(), file.IsDir())
	}
	return fileNames
}

func upload(c echo.Context) error {
	// Read form fields
	// name := c.FormValue("name")
	// email := c.FormValue("email")

	//------------
	// Read files
	//------------

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create("./tmp/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	// return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files with fields name=%s and email=%s.</p>", len(files), name, email))
	return c.HTML(http.StatusOK, fmt.Sprintf("<p>Uploaded successfully %d files with fields .</p>", len(files)))

}

func main() {
	listFiles := scanFiles()
	// clearConsole()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("TEMP DIR:", path+"/tmp/")

	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = renderer

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	// Named route "foobar"
	e.GET("/upload", func(c echo.Context) error {
		return c.Render(http.StatusOK, "template.html", map[string]interface{}{
			"FileNames": listFiles,
		})
	}).Name = "foobar"

	e.Logger.Fatal(e.Start(":8000"))

}
