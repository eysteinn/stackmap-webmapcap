package mapfile

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"gitlab.com/EysteinnSig/stackmap-mapserver/data"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

type WMS struct {
	DstDir       string
	MapfilesRoot string
	SQLData      types.SQLData
	ApiHost      string
}

/*
func (w *WMS) Init(ApiHost string) error {

		return nil
	}
*/
func (w *WMS) CreateOutputDir(project string) (string, error) {
	outdir := filepath.Join(w.MapfilesRoot, project)

	_, err := os.Stat(filepath.Join(outdir, "rasters.map"))
	if os.IsNotExist(err) {
		err := os.MkdirAll(outdir, 0700) // 0700)
		if err != nil {
			return outdir, err
		}
	} else {
		return outdir, err
	}

	return outdir, nil
}
func (w *WMS) GenerateProductsIndexMapfile(project string) error {

	products, err := types.GetProducts(w.ApiHost, project)
	if err != nil {
		return err
	}
	prodstr := ""
	for _, prod := range products.Products {
		prodstr += fmt.Sprintf("INCLUDE \"product_%s.map\"\n", prod)
	}
	outdir, err := w.CreateOutputDir(project)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(outdir, "products.map"), []byte(prodstr), 0655)
	return err
}
func (w *WMS) GenerateProductMapfile(project string, product string) error {
	if err := w.SQLData.Validate(); err != nil {
		return err
	}
	prodTimes, err := types.GetTimes(w.ApiHost, project, product)
	if err != nil {
		return err
	}
	layerData := types.MapLayerData{
		SQLData:      w.SQLData,
		Project:      project,
		ProductTimes: prodTimes,
	}
	log.Println("Creating mapfiles")
	b, err := CreateLayerFile(&layerData)
	if err != nil {
		return err
	}

	/*outdir := filepath.Join(w.MapfilesRoot, project)
	err = os.MkdirAll(outdir, 0700) // 0700)
	if err != nil {
		return err
	}*/
	outdir, err := w.CreateOutputDir(project)
	if err != nil {
		return err
	}

	prod_filepath := filepath.Join(outdir, "product_"+product+".map")

	//log.Println(string(b))
	log.Println("Writing to:", prod_filepath)
	err = os.WriteFile(prod_filepath, b, 0655)
	if err != nil {
		return err
	}
	return nil
}
func (w *WMS) GenerateMapfile(project string, product string) error {
	/*if err := w.SQLData.Validate(); err != nil {
		return err
	}
	prodTimes, err := types.GetTimes(w.ApiHost, project, product)
	if err != nil {
		return err
	}
	layerData := types.MapLayerData{
		SQLData:      w.SQLData,
		Project:      project,
		ProductTimes: prodTimes,
	}
	log.Println("Creating mapfiles")
	b, err := CreateLayerFile(&layerData)
	if err != nil {
		return err
	}

	outdir := filepath.Join(w.MapfilesRoot, project)
	err = os.MkdirAll(outdir, 0700) // 0700)
	if err != nil {
		return err
	}

	prod_filepath := filepath.Join(outdir, "product_"+product+".map")

	log.Println(string(b))
	log.Println("Writing to:", prod_filepath)
	err = os.WriteFile(prod_filepath, b, 0655)
	if err != nil {
		return err
	}

	products, err := types.GetProducts(w.ApiHost, project)
	if err != nil {
		return err
	}
	prodstr := ""
	for _, prod := range products.Products {
		prodstr += fmt.Sprintf("INCLUDE \"product_%s.map\"\n", prod)
	}
	err = os.WriteFile(filepath.Join(outdir, "products.map"), []byte(prodstr), 0655)
	*/
	outdir, err := w.CreateOutputDir(project)
	if err != nil {
		return nil
	}
	err = w.GenerateProductMapfile(project, product)
	if err != nil {
		return err
	}

	err = w.GenerateProductsIndexMapfile(project)
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(outdir, "rasters.map")); os.IsNotExist(err) {
		err = os.WriteFile(filepath.Join(outdir, "rasters.map"), data.RastersMap, 0655)
	}

	return nil
}

func (w *WMS) OnReceived(prod *types.ProductModified) error {
	log.Println("Got a message:", prod.Project, "-", prod.Product)
	err := w.GenerateMapfile(prod.Project, prod.Product)
	return err
	/*if err := w.SQLData.Validate(); err != nil {
		return err
	}
	prodTimes, err := types.GetTimes(w.ApiHost, prod.Project, prod.Product)
	if err != nil {
		return err
	}
	layerData := types.MapLayerData{
		SQLData:      w.SQLData,
		Project:      prod.Project,
		ProductTimes: prodTimes,
	}
	log.Println("Creating mapfiles")
	b, err := CreateLayerFile(&layerData)
	if err != nil {
		return err
	}

	outdir := filepath.Join(w.MapfilesRoot, prod.Project)
	err = os.MkdirAll(outdir, 0700) // 0700)
	if err != nil {
		return err
	}

	prod_filepath := filepath.Join(outdir, "product_"+prod.Product+".map")

	log.Println(string(b))
	log.Println("Writing to:", prod_filepath)
	err = os.WriteFile(prod_filepath, b, 0655)
	if err != nil {
		return err
	}

	products, err := types.GetProducts(w.ApiHost, prod.Project)
	if err != nil {
		return err
	}
	prodstr := ""
	for _, prod := range products.Products {
		prodstr += fmt.Sprintf("INCLUDE \"product_%s.map\"\n", prod)
	}
	err = os.WriteFile(filepath.Join(outdir, "products.map"), []byte(prodstr), 0655)

	if _, err := os.Stat(filepath.Join(outdir, "rasters.map")); os.IsNotExist(err) {
		err = os.WriteFile(filepath.Join(outdir, "rasters.map"), data.RastersMap, 0655)
	}

	return nil*/
}

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

/*func CreateProductsFile() ([]byte, error) {

}*/

func CreateLayerFile(layer *types.MapLayerData) ([]byte, error) {
	if err := layer.SQLData.Validate(); err != nil {
		return nil, err
	}
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
