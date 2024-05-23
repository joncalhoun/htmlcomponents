package htmlcomponents

import (
	"bytes"
	"fmt"
	"html/template"
)

// DataGrid can be used to generate any HTML table. It accepts a slice of
// columns used as headers, then a slice of rows which is a map of column name
// to value. This function will return the HTML for the table.
func DataGrid(columns []any, rows []map[any]any) (template.HTML, error) {
	var data struct {
		Columns []any
		Rows    []map[any]any
	}
	data.Columns = columns
	data.Rows = rows
	var buffer bytes.Buffer
	err := dataGridTpl.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("creating data grid: %w", err)
	}
	return template.HTML(buffer.String()), nil
}

var dataGridTpl = template.Must(template.New("datagrid").Parse(dataGridHTML))

const dataGridHTML = `
<table>
	<thead>
		<tr>
			{{range .Columns}}
				<th>{{.}}</th>
			{{end}}
		</tr>
	</thead>
	<tbody>
		{{$columns := .Columns}}
		{{range $row := .Rows}}
			<tr>
				{{range $col := $columns}}
					<!-- This is like doing $row[$col] in Go since it is a map -->
					<td>{{index $row $col}}</td>
				{{end}}
			</tr>
		{{end}}
	</tbody>
</table>
`
