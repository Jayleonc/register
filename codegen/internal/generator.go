package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Jayleonc/register/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
)

type templateData struct {
	Structs     map[string]string
	Interfaces  []registry.Api
	BaseURL     string
	PackageName string
}

func ToCamelCase(s string) string {
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.Title(s)
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func paramType(t interface{}) string {
	switch t.(type) {
	case map[string]interface{}:
		return structName(t)
	case string:
		return "string"
	case int:
		return "int"
	case bool:
		return "bool"
	default:
		return "interface{}"
	}
}

func structName(t interface{}) string {
	m := t.(map[string]interface{})
	nameParts := make([]string, 0, len(m))
	for k := range m {
		nameParts = append(nameParts, ToCamelCase(k))
	}
	return strings.Join(nameParts, "")
}

func mapToStructFields(m map[string]interface{}, structs map[string]string) string {
	fields := ""
	for k, v := range m {
		fieldType := paramType(v)
		if nestedStruct, ok := v.(map[string]interface{}); ok {
			nestedStructName := ToCamelCase(k)
			structs[nestedStructName] = fmt.Sprintf("type %s struct { %s }", nestedStructName, mapToStructFields(nestedStruct, structs))
			fieldType = nestedStructName
		}
		fields += fmt.Sprintf("%s %s `json:\"%s\"`;", ToCamelCase(k), fieldType, k)
	}
	return fields
}

func GenerateStructDefinitions(interfaces []registry.Api) map[string]string {
	structs := make(map[string]string)
	for _, api := range interfaces {
		for _, param := range api.Params {
			if t, ok := param.Type.(map[string]interface{}); ok {
				structName := ToCamelCase(param.Name)
				if _, exists := structs[structName]; !exists {
					structs[structName] = fmt.Sprintf("type %s struct { %s }", structName, mapToStructFields(t, structs))
				}
			}
		}
	}
	return structs
}

func GenerateClientCode(serviceName, outputPath string, etcdClient *clientv3.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	registryClient, err := registry.NewClient(
		registry.ClientWithResolver(registry.NewEtcdResolver(etcdClient)),
	)
	if err != nil {
		return fmt.Errorf("failed to create registry client: %v", err)
	}

	interfaces, address, err := registryClient.GetServiceInterfaces(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("failed to get service interfaces: %v", err)
	}

	baseURL := fmt.Sprintf("http://%s", address)
	structs := GenerateStructDefinitions(interfaces)

	packageName := strings.Split(strings.TrimPrefix(interfaces[0].Path, "/"), "/")[0]

	data := templateData{
		Structs:     structs,
		Interfaces:  interfaces,
		BaseURL:     baseURL,
		PackageName: packageName,
	}

	funcMap := template.FuncMap{
		"ToCamelCase": ToCamelCase,
		"paramType":   paramType,
		"structName":  structName,
	}

	tmpl, err := template.New("client").Funcs(funcMap).Parse(ClientTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// 写入文件
	filePath := fmt.Sprintf("%s/%s.go", outputPath, packageName)
	err = ioutil.WriteFile(filePath, output.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
