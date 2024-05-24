package helpers

import (
	"html/template"
	"os"
)

func BuildFile(fileName, baseFilePath, footerFilePath, publicFilePath, templateFilePath string, data any) error {
	tmpl, err := template.ParseFiles(baseFilePath, footerFilePath, templateFilePath)

	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(publicFilePath)

	if err != nil {
		return err
	}

	err = tmpl.Execute(f, data)

	if err != nil {
		return err
	}

	err = f.Close()

	if err != nil {
		return err
	}

	return err
}
