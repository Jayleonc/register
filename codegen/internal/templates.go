package internal

const ClientTemplate = `
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

{{range $key, $value := .Structs}}
{{$value}}
{{end}}

{{range .Interfaces}}
type {{.ResponseType}} struct {
	{{range .Returns}}{{ToCamelCase .Name}} {{paramType .Type}} ` + "`json:\"{{.Name}}\"`" + `
	{{end}}
}

func (c *Client) {{.Method}}_{{.Path | ToCamelCase}}(req {{.RequestType}}) ({{.ResponseType}}, error) {
	var resp {{.ResponseType}}

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
