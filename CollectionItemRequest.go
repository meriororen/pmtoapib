package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

type CollectionItemRequest struct {
	Url         CollectionItemRequestUrl `json:"url"`
	Method      string                   `json:"method"`
	Header      []RequestHeader          `json:"header"`
	Body        RequestBody              `json:"body"`
	Description string                   `json:"description"`
}

func prependHttp(str string) string {
	if strings.Index(str, "http://") != 0 && strings.Index(str, "https://") != 0 {
		return fmt.Sprintf("https://%v", str)
	}
	return str
}

func (r CollectionItemRequest) ShortUrl() string {
	if len(r.Url.Paths) == 0 {
		u, err := url.Parse(prependHttp(r.Url.Raw))
		if err != nil {
			log.Fatal(r.Body, " -> ", err)
		}

		return u.Path
	} else {
		ps := []string{""}
		for _, p := range r.Url.Paths {
			if p[0] == ':' {
				p = fmt.Sprintf("{%v}", p[1:])
			}
			ps = append(ps, p)
		}

		return strings.Join(ps, "/")
	}

	return ""
}

func (r CollectionItemRequest) UrlParameterList() []UrlParameter {
	if len(r.Url.Paths) == 0 {
		u, _ := url.Parse(prependHttp(r.Url.Raw))
		m, _ := url.ParseQuery(u.RawQuery)

		queries := []UrlParameter{}
		for k, v := range m {
			queries = append(queries, UrlParameter{Key: k, Value: v[0]})
		}

		return queries
	} else {
		return append(r.Url.Queries, r.Url.Variables...)
	}
}

func (r CollectionItemRequest) UrlParameterListString() string {
	ps := []string{}
	for _, p := range r.UrlParameterList() {
		ps = append(ps, p.Key)
	}

	return strings.Join(ps, ",")
}

func (r CollectionItemRequest) IsExcluded() bool {
	return strings.Contains(r.Description, "pmtoapib_exclude")
}
