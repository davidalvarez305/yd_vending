package helpers

import (
	"html/template"
	"net/http"
	"os"
)

func BuildFile(fileName, baseFilePath, footerFilePath, publicFilePath, templateFilePath string, data any) error {
	if data == nil {
		if _, err := os.Stat(publicFilePath); err == nil {
			return nil
		}
	}

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
