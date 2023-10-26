package format

import (
	"bytes"
	"errors"
	"os"
	"text/template"

	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

func GetTemplateMapfile(product string) (string, error) {

	filename := "./data/" + product + ".map"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		filename = "./data/default.map"
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetMapfile(layer *types.MapLayerData) ([]byte, error) {
	templ := template.New("MapserverLayer")
	tmpldata, err := GetTemplateMapfile(layer.Product)
	if err != nil {
		return nil, err
	}
	_, err = templ.Parse(tmpldata)
	if err != nil {

		return nil, err
	}
	var buff bytes.Buffer
	if err := templ.Execute(&buff, layer); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
