package helpers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func ServeContent(w http.ResponseWriter, templateFilePaths []string, data any) error {
	tmpl, err := template.ParseFiles(templateFilePaths...)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building template.", http.StatusInternalServerError)
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error executing template.", http.StatusInternalServerError)
		return err
	}

	return nil
}

func BuildFile(fileName, publicFilePath string, templateFilePaths []string, data any) error {
	tmpl, err := template.ParseFiles(templateFilePaths...)

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

func GetUserIPFromRequest(r *http.Request) string {
	// Try to get the IP from the X-Forwarded-For header (used by proxies and load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Try to get the IP from the X-Real-IP header (used by some proxies)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to using the RemoteAddr from the request
	ip := r.RemoteAddr
	return ip
}
