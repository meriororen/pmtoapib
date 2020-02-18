package main

import (
	"encoding/json"
	"fmt"
)

type UrlParameter struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type UrlVariable = UrlParameter

type CollectionItemRequestUrl struct {
	Raw       string         `json:"raw"`
	Hosts     []string       `json:"host"`
	Port      string         `json:"port"`
	Paths     []string       `json:"path"`
	Queries   []UrlParameter `json:"query"`
	Variables []UrlVariable  `json:"variable"`
}

func (p UrlParameter) BacktickedValue() string {
	return fmt.Sprintf("`%v`", p.Value)
}

func (u *CollectionItemRequestUrl) UnmarshalJSON(data []byte) error {
	var raw interface{}
	json.Unmarshal(data, &raw)

	switch raw := raw.(type) {
	case string:
		*u = CollectionItemRequestUrl{Raw: raw}
	case map[string]interface{}:
		*u = CollectionItemRequestUrl{
			Raw: raw["raw"].(string),
		}
		if hosts := raw["host"].([]interface{}); hosts != nil {
			u.Hosts = make([]string, len(hosts))
			for i, h := range hosts {
				u.Hosts[i] = h.(string)
			}
		}
		if paths := raw["path"].([]interface{}); paths != nil {
			u.Paths = make([]string, len(paths))
			for i, p := range paths {
				u.Paths[i] = p.(string)
			}
		}
		if raw["query"] != nil {
			if queries := raw["query"].([]interface{}); len(queries) > 0 {
				u.Queries = make([]UrlParameter, len(queries))
				for i, q := range queries {
					if q := q.(map[string]interface{}); q != nil {
						if q["key"] != nil {
							u.Queries[i].Key = q["key"].(string)
						}
						if q["value"] != nil {
							u.Queries[i].Value = q["value"].(string)
						}
						if q["description"] != nil {
							u.Queries[i].Description = q["description"].(string)
						}
					}
				}
			}
		}
		if raw["variable"] != nil {
			if vars := raw["variable"].([]interface{}); len(vars) > 0 {
				u.Variables = make([]UrlVariable, len(vars))
				for i, v := range vars {
					if v := v.(map[string]interface{}); v != nil {
						if v["key"] != nil {
							u.Variables[i].Key = v["key"].(string)
						}
						if v["value"] != nil {
							u.Variables[i].Value = v["value"].(string)
						}
						if v["description"] != nil {
							u.Variables[i].Description = v["description"].(string)
						}
					}
				}
			}
		}
	}

	return nil
}
