package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

func parseVariables(byteData []byte) (map[string]interface{}, error) {
	var variablesMap map[string]interface{}
	err := yaml.Unmarshal(byteData, &variablesMap)

	return variablesMap, err
}

func replaceTemplatePlaceholders(template string, variables map[string]interface{}) (string, error) {
	for key, value := range variables {
		matched, err := regexp.MatchString("{{ "+key+" }}", template)
		if err != nil {
			return "", err
		}

		if !matched {
			continue
		}
		template = strings.ReplaceAll(template, "{{ "+key+" }}", fmt.Sprint(value))
	}

	return template, nil
}

func main() {
	templateFilename := os.Args[1]
	variablesFilename := os.Args[2]
	outputFilename := "output.txt"
	if len(os.Args) > 3 {
		outputFilename = os.Args[3]
	}

	// read a file of variables definition
	variablesByteData, err := os.ReadFile(variablesFilename)
	if err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}

	// parse variables definition
	variablesMap, err := parseVariables(variablesByteData)
	if err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}

	// open template file
	templateByteData, err := os.ReadFile(templateFilename)
	if err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}

	// replace template placeholders
	replacedString, err := replaceTemplatePlaceholders(string(templateByteData), variablesMap)
	if err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}

	// output a new file
	_, err = os.Stat(outputFilename)
	if err == nil {
		log.Fatalln(outputFilename, "is already exists.")
		os.Exit(1)
	}

	err = os.WriteFile(outputFilename, []byte(replacedString), os.ModePerm)
	if err != nil {
		log.Fatalf("error: %s\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
