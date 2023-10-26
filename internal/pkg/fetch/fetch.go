package fetch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gitlab.com/EysteinnSig/stackmap-mapserver/data"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/format"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/logger"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

func FetchAllProducts(outdir string, apihost string, project string, sqldata types.SQLData) error {
	prods := types.UniqueProducts{}
	//err := GetData(apihost+"/api/v1/products?project="+project, &prods)
	logger.GetLogger().Debug("Fetching from: ", apihost+"/api/v1/projects/"+project+"/products")
	err := GetData(apihost+"/api/v1/projects/"+project+"/products", &prods)
	if err != nil {
		return err
	}

	for _, prod := range prods.Products {
		logger.GetLogger().Debugf("Fetching product: %s\n", prod)
		err := FetchProduct(project, prod, outdir, apihost, sqldata)
		if err != nil {
			return err
		}
	}

	prodstr := ""
	for _, prod := range prods.Products {
		prodstr += fmt.Sprintf("INCLUDE \"product_%s.map\"\n", prod)

	}
	err = os.MkdirAll(outdir, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(outdir, "products.map"), []byte(prodstr), 0644)
	if err != nil {
		return err
	}

	return nil
}

func FetchProduct(project string, product string, outdir string, apihost string, sqldata types.SQLData) error {

	//product := "hrit-ash"
	//pt := ProductTimes{}
	pt := types.MapLayerData{SQLData: sqldata}
	//GetData("http://localhost:3000/api/v1/times?product="+product, &pt)
	//err := GetData(apihost+"/api/v1/times?product="+product+"&project="+project, &pt)
	err := GetData(apihost+"/api/v1/projects/"+project+"/products/"+product+"/times", &pt)
	if err != nil {
		return err
	}
	sort.SliceStable(pt.Times, func(i, j int) bool {
		return pt.Times[i].Before(pt.Times[j])
	})
	pt.StartRange = pt.Times[0]
	pt.EndRange = pt.Times[len(pt.Times)-1]
	pt.DefaultTime = pt.EndRange
	pt.Project = project

	/*pt.SQLDB = sqldata.SQLDB
	pt.SQLHost = sqldata.SQLHost
	pt.SQLPass = sqldata.SQLPass
	pt.SQLUser = sqldata.SQLUser*/

	//fmt.Println(pt.AllTimesString())
	//txt, err := Replace(product, timesstr, mintime, maxtime, mintime)
	bufr, err := format.GetMapfile(&pt)
	if err != nil {
		/*fmt.Println(err)
		os.Exit(1)*/
		return err
	}
	err = os.MkdirAll(outdir, 0655) // 0700)
	if err != nil {
		logger.GetLogger().Error(err)
	}

	if _, err := os.Stat(filepath.Join(outdir, "rasters.map")); os.IsNotExist(err) {
		err = os.WriteFile(filepath.Join(outdir, "rasters.map"), data.RastersMap, 0655)
	}
	outpath := filepath.Join(outdir, "product_"+product+".map")
	fmt.Println("Writing to: ", outpath)
	err = os.WriteFile(outpath, bufr, 0644)
	if err != nil {
		return err
	}
	return nil
}
