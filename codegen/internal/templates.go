package internal

const ClientTemplate = `
package {{.PackageName}}

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var url = "{{.BaseURL}}"

type {{.PackageName | ToCamelCase}} struct {
	baseURL string
}

func New{{.PackageName | ToCamelCase}}() *{{.PackageName | ToCamelCase}} {
	return &{{.PackageName | ToCamelCase}}{baseURL: url}
}

{{range $key, $value := .Structs}}
{{$value}}
{{end}}

{{range .Interfaces}}
type {{.Path | ToCamelCase}}Response struct {
	{{range .Returns}}{{ToCamelCase .Name}} {{paramType .Type}} ` + "`json:\"{{.Name}}\"`" + `
	{{end}}
}

func (c *{{$.PackageName | ToCamelCase}}) {{.Method}}_{{.Path | ToCamelCase}}(req {{(index .Params 0).Name | ToCamelCase}}) ({{.Path | ToCamelCase}}Response, error) {
	var resp {{.Path | ToCamelCase}}Response

	url := fmt.Sprintf("%s{{.Path}}", c.baseURL)
	body, err := json.Marshal(req)
	if err != nil {
		return resp, err
	}

	httpResp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return resp, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("unexpected status: %v", httpResp.Status)
	}

	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
{{end}}
`
