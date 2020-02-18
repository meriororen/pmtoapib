package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

type CollectionItem struct {
	Name      string                  `json:"name"`
	Items     []CollectionItem        `json:"item"`
	Request   CollectionItemRequest   `json:"request"`
	Responses CollectionItemResponses `json:"response"`
}

func applyEnvVars(url string) string {
	if DefaultCollectionEnv.Name == "" {
		//log.Println("No Env")
		return url
	}

	str := url
	for _, e := range DefaultCollectionEnv.Values {
		if e.Enabled {
			str = strings.Replace(str,
				fmt.Sprintf("{{%v}}", e.Key), e.Value, -1)
			//	log.Println("applying => ", url, str)
		}
	}
	return str
}

func (i CollectionItem) Markup() template.HTML {
	tpl := `
{{ $length := len .Items }} {{ if eq $length 0 }}
### {{ .Name }} [{{ .Request.Method }} {{ .Request.ShortUrl }}{{ if .Request.UrlParameterListString }}{?{{ .Request.UrlParameterListString }}}{{ end }}]

{{ if .Request.Description }}{{ .Request.Description }}{{ else }}No Description{{ end }}
{{ if .Request.UrlParameterList }}
+ Parameters

	{{ range .Request.UrlParameterList }}+ {{ .Key }}: {{ .BacktickedValue }} (string, required) - {{ if .Description }}{{ .Description }}{{ else }}No Description{{ end }}
    {{ end }}{{ end }}
+ Request

    + Headers
            {{ range .Request.Header }}{{ if not .Disabled }}
            {{ .Key }}: {{ .Value }}{{ end }}{{ end }}
    {{ if .Request.Body.Raw }}
    + Body 

    	    {{ .Request.Body.RawString }}
    {{ end }}
{{ .ResponseSectionMarkup }}
{{ else }}
{{ range .Items }} 
	{{ .Markup }}
{{ end }}
{{ end }}
`

	i.Request.Url.Raw = applyEnvVars(i.Request.Url.Raw)

	if len(i.Items) > 0 {
		// This current item is probably a folder
		// Rename nested items with current folder name prepended
		for idx, item := range i.Items {
			i.Items[idx].Name = fmt.Sprint(i.Name, " > ", item.Name)
		}

		// Create all response files for nested items
		for _, file := range getResponseFiles(i.Items) {
			path := applyEnvVars(file["path"])
			// in case path is a full fledge url now
			if strings.Index(path, "://") >= 0 {
				u, _ := url.Parse(path)
				path = fmt.Sprintf("responses%v", u.Path)
			}
			//log.Println("Applied vars for", file["path"], "to", path)
			writeToFile(
				fmt.Sprintf("%v/%v", filepath.Clean(DefaultConfig.DestinationPath), path),
				file["body"],
				DefaultConfig.ForceResponsesCreation,
			)
		}
	}

	log.Println(i.Name, "=> ", i.Request.ShortUrl())

	t := template.New("Item Template")
	t, _ = t.Parse(tpl)

	var doc bytes.Buffer
	t.Execute(&doc, i)
	s := doc.String()

	return template.HTML(s)
}

func (i CollectionItem) ResponseSectionMarkup() template.HTML {
	tpl :=
		`{{ range .ResponseList }}
+ Response {{ .Code }}{{ if .ContentType }} ({{ .ContentType }}){{ end }}

    + Headers
            {{ range .Header }}{{ if not .Hidden }}
            {{ .Key }}: {{ .Value }}{{ end }}{{ end }}

    + Body

            {{ .BodyIncludeString $.Request }}
{{ end }}`

	t := template.New("Response Section Template")
	t, _ = t.Parse(tpl)

	var doc bytes.Buffer
	t.Execute(&doc, i)
	s := doc.String()

	return template.HTML(s)
}

func (i CollectionItem) ResponseList() CollectionItemResponses {
	responses := CollectionItemResponses{}

	dummyTwoHundredResponse := CollectionItemResponse{
		Id:     "1",
		Name:   "200",
		Status: "OK",
		Code:   200,
		Header: []ResponseHeader{
			{
				Key:         "Content-Type",
				Value:       "application/json",
				Name:        "Content-Typ",
				Description: "The mime type of this content",
			},
			{
				Key:         "NAME",
				Value:       "VALUE",
				Name:        "NAME",
				Description: "Dummy Header",
			},
		},
		Body: "{}",
	}

	if len(i.Responses) == 0 {
		responses = append(responses, dummyTwoHundredResponse)
		return responses
	}

	responses = i.Responses

	hasTwoHundred := false

	for _, response := range responses {
		if response.Code == 200 {
			hasTwoHundred = true
		}
	}

	if !hasTwoHundred {
		responses = append(responses, dummyTwoHundredResponse)
	}

	sort.Sort(responses)

	return responses
}
