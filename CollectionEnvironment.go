package main

var DefaultCollectionEnv CollectionEnvironment

type CollectionEnvironmentVar struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

type CollectionEnvironment struct {
	ID     string                     `json:"id"`
	Name   string                     `json:"name"`
	Values []CollectionEnvironmentVar `json:"values"`
}
