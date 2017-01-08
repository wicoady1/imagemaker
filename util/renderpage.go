package util

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

const (
	TemplateImageMaker string = "intools.imagemaker.html"
)

func RenderPage(w http.ResponseWriter, templateName string, templateData map[string]string) error {
	var fileTemplateName string
	var oriTemplateName string

	if templateName == "imagemaker" {
		fileTemplateName = "templates/" + TemplateImageMaker
		oriTemplateName = TemplateImageMaker
	}

	dat, err := ioutil.ReadFile(fileTemplateName)
	if err != nil {
		return err
	}

	t, err := template.New(oriTemplateName).Parse(string(dat))
	if err != nil {
		return err
	}

	if err = t.ExecuteTemplate(w, oriTemplateName, templateData); err != nil {
		return err
	}

	return nil
}
