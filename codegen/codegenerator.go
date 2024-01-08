package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CodegenFunc func() error

func codegenerator(filename string, data any) CodegenFunc {
	return func() error {
		funcs := template.FuncMap{
			"ToUpper": strings.ToUpper,
			"Title":   func(s string) string { return cases.Title(language.English, cases.NoLower).String(s) },
		}

		tags := template.New(fmt.Sprintf("%s.go", filename)).Funcs(funcs)
		tags, err := tags.ParseGlob("./_templates/*")
		if err != nil {
			return err
		}

		file, err := os.Create(fmt.Sprintf("../%s.go", filename))
		defer file.Close()
		if err != nil {
			return err
		}

		err = tags.Execute(file, data)
		if err != nil {
			return err
		}

		return nil
	}
}
