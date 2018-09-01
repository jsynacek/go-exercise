// A simple http file browser.
// Written as an exercise to learn how net/http and html/template work.
package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

const tmpl = `<table>
  <tr><th><b>{{.Root}}</b></th><th><b>size</b></th><th><b>mod timeervices</b></th></tr>
  {{range .Rows}}<tr><td><a href="{{.Path}}">{{.Name}}</a></td><td>{{.Size}}</td><td>{{.ModTime}}</td></tr>
  {{end}}
</table>
`

type File struct {
	Path    string
	Name    string
	Size    int64
	ModTime time.Time
}

type Table struct {
	Root string
	Rows []File
}

func htmlError(w http.ResponseWriter, err error) {
	w.Write([]byte(fmt.Sprintf(`<b style="color: #ee2222">error: %v</b>`, err)))
}

func serveDir(w http.ResponseWriter, dir string) {
	t, err := template.New("table").Parse(tmpl)
	if err != nil {
		htmlError(w, err)
		return
	}
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		htmlError(w, err)
		return
	}
	// Reading and writing is done in two passes, because I couldn't figure out
	// how to use http/template to do it in one. The problem was joining the dir
	// path with the file name.
	var rows []File
	for _, fi := range fileInfos {
		name := fi.Name()
		if fi.IsDir() {
			name += "/"
		}
		row := File{path.Join(dir, fi.Name()), name, fi.Size(), fi.ModTime()}
		rows = append(rows, row)
	}
	err = t.Execute(w, Table{dir, rows})
	if err != nil {
		htmlError(w, err)
		return
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	root := r.RequestURI
	fi, err := os.Stat(root)
	if err != nil {
		htmlError(w, err)
		return
	}
	if fi.IsDir() {
		serveDir(w, root)
	} else {
		http.ServeFile(w, r, root)
	}
}

func main() {
	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
