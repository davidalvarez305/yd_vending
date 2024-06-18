package helpers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/types"
)

func ServeSuccessModal(w http.ResponseWriter, r *http.Request, ctx types.SuccessModal) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	templatePath := constants.PARTIAL_TEMPLATES_DIR + "modal.html"

	template, err := BuildStringFromTemplate(templatePath, ctx.TemplateName, ctx)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error parsing template.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(template))
}

func BuildStringFromTemplate(templatePath, templateName string, data any) (string, error) {
	var output string
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	tmpl, err := template.New(templateName).Parse(string(templateContent))
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	var body strings.Builder
	err = tmpl.Execute(&body, data)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return output, err
	}

	return body.String(), err
}

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

func BuildFile(outputPath string, templateFilePaths []string, data any) error {
	tmpl, err := template.ParseFiles(templateFilePaths...)

	if err != nil {
		return err
	}

	var f *os.File
	f, err = os.Create(outputPath)

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
