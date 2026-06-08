package llm

import (
	"bytes"
	"fmt"
	"text/template"
)

type DigestData struct {
	StripeRevenue        float64
	StripeFailedPayments int
	StripeAnomaly        bool
	StripeDelta          float64
	CalendlyEvents       string
	GmailMessages        string
}

func BuildPrompt(data DigestData, templatePath string) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	// execute into a bytes.Buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data) // like writing the template's data to a paper
	if err != nil {
    	return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}
