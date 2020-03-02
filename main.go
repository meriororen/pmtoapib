package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func getApibFileContent(c Collection) string {
	tpl := `
# Group {{ .Info.Name }}

{{ .Info.Description }} 
 
{{ range .Items }}
	{{ .Markup }}
{{ end }}
`

	t := template.New("Template")
	t, _ = t.Parse(tpl)

	var doc bytes.Buffer
	t.Execute(&doc, c)
	s := doc.String()

	return s
}

func getResponseFiles(items []CollectionItem) []map[string]string {
	var files []map[string]string

	for _, item := range items {
		responses := item.ResponseList()
		for _, response := range responses {
			m := map[string]string{}
			m["path"] = response.BodyIncludePath(item.Request)
			m["body"] = response.FormattedBody()
			files = append(files, m)
		}
	}

	return files
}

func writeToFile(path string, content string, force bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) || force {
		os.MkdirAll(filepath.Dir(path), os.ModePerm)
		err := ioutil.WriteFile(path, []byte(content), 0644)
		if err == nil {
			fmt.Printf("Created %v\n", path)
		}
	}
}

func shouldWriteFiles(c Config) bool {
	return c.DumpRequest == ""
}

func main() {

	config := Config{}
	config.Init()

	if config.CollectionPath == "" {
		fmt.Println("No collection file defined!")
		return
	}

	file, _ := ioutil.ReadFile(config.CollectionPath)
	var c Collection
	err := json.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Error in unmarshaling postman collection: ", err)
	}

	log.Println(c.Items[1].Name)

	apibFileName := strings.Replace(c.Info.Name, " ", "-", -1)

	if config.ApibFileName != "" {
		apibFileName = config.ApibFileName
	}

	if config.EnvironmentPath != "" {
		file, _ := ioutil.ReadFile(config.EnvironmentPath)
		json.Unmarshal(file, &DefaultCollectionEnv)
	}

	apibFile := getApibFileContent(c)

	if shouldWriteFiles(config) {
		writeToFile(
			fmt.Sprintf("%v/%v.apib", filepath.Clean(config.DestinationPath), apibFileName),
			apibFile,
			config.ForceApibCreation,
		)

		for _, file := range getResponseFiles(c.Items) {
			writeToFile(
				fmt.Sprintf("%v/%v", filepath.Clean(config.DestinationPath), file["path"]),
				file["body"],
				config.ForceResponsesCreation,
			)
		}
	}

	if config.DumpRequest != "" {
		for _, request := range c.Items {
			if request.Name == config.DumpRequest {
				fmt.Println(request.Markup())
			}
		}
	}
}
