package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/joncalhoun/htmlcomponents"
)

type Gallery struct {
	ID    int
	Title string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Pretend we get this from a database or something like that.
		galleries := []Gallery{
			{ID: 1, Title: "Nature"},
			{ID: 2, Title: "City"},
			{ID: 3, Title: "People"},
		}

		html, err := galleriesDataGrid(r, galleries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(html))
	})
	http.ListenAndServe(":8080", mux)
}

// galleriesDataGrid is repsonsible for translating a slice of Galleries into
// the format expected for a DataGrid, and then returning the resulting HTML.
func galleriesDataGrid(r *http.Request, galleries []Gallery) (template.HTML, error) {
	// Build the data for the data grid. I am opting to do this manually
	// instead of using reflect so that I have more control, but we could also
	// use reflect. https://github.com/joncalhoun/form has some examples of
	// using reflect and struct tags to build forms, which is similar. I also don't see a great way to generate actions without doing those manually to some degree.
	var columns []any
	columns = append(columns, "ID")
	columns = append(columns, "Title")
	columns = append(columns, "Actions")

	var rows []map[any]any
	for _, gallery := range galleries {
		row := make(map[any]any)
		row["ID"] = gallery.ID
		row["Title"] = gallery.Title
		action, err := galleryActionHTML(r, gallery)
		if err != nil {
			return "", fmt.Errorf("creating data grid: %w", err)
		}
		row["Actions"] = action
		rows = append(rows, row)
	}
	return htmlcomponents.DataGrid(columns, rows)
}

func galleryActionHTML(r *http.Request, gallery Gallery) (template.HTML, error) {
	var data struct {
		ID        int
		CSRFField template.HTML
	}
	data.ID = gallery.ID
	data.CSRFField = csrf.TemplateField(r)
	var buffer bytes.Buffer
	err := actionsTpl.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("creating gallery actions: %w", err)
	}
	return template.HTML(buffer.String()), nil
}

var actionsTpl = template.Must(template.New("actions").Parse(actionsHTML))

// Note: This could be in a separate file.
const actionsHTML = `
<a href="/galleries/{{.ID}}" class="add-these-here">View</a>
<a href="/galleries/{{.ID}}/edit" class="add-these-here">Edit</a>
<form action="/galleries/{{.ID}}/delete" method="POST" class="add-these-here" onsubmit="return confirm('Are you sure you want to delete this gallery?')">
	<div class="hidden">{{.CSRFField}}</div>
	<button type="submit">Delete</button>
</form>
`
