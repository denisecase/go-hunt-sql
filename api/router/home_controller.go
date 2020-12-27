package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// IndexPageData defines the data needed for the home (default) page
type IndexPageData struct {
	Title    string
	SubTitle string
}

// ShowHome handles request to home
func (server *Server) ShowHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Show Home")
	data := IndexPageData{
		Title:    "Hunt",
		SubTitle: "Welcome to the new Hunt App (made with Go!)",
	}
	templ.ExecuteTemplate(w, "main", data)
}

var templ = func() *template.Template {
	t := template.New("")
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			fmt.Println(path)
			_, err = t.ParseFiles(path)
			if err != nil {
				fmt.Println(err)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}
	return t
}()
