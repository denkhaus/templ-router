package main

import (
	"fmt"
	"log"
	"github.com/denkhaus/templ-router/pkg/shared"
)

func main() {
	found, config, err := shared.ParseYAMLMetadata("test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found: %v\n", found)
	fmt.Printf("Title: %v\n", config.Metadata.Title)
	fmt.Printf("Company: %v\n", config.Metadata.Custom["company_name"])

	i18n := config.GetMultiLocaleI18n()
	fmt.Printf("EN welcome: %v\n", i18n["en"]["welcome"])
	fmt.Printf("DE welcome: %v\n", i18n["de"]["welcome"])

	fmt.Printf("Auth type: %v\n", config.Auth.Type)
}