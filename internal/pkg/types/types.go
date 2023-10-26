package types

import (
	"errors"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/apirequest"
)

type UniqueProducts struct {
	Products []string `json:"product_names"`
}

func GetProducts(host string, project string) (UniqueProducts, error) {
	up := UniqueProducts{}

	url := host + "/api/v1/projects/" + project + "/products"
	err := apirequest.GetData(url, &up)
	if err != nil {
		return up, err
	}
	log.Println("Found products:", up.Products)
	return up, nil
}

type UniqueProjects struct {
	Projects []string `json:"project_names"`
}

type ProductTimes struct {
	Product string      `json:"product"`
	Times   []time.Time `json:"times"`
}

/*func (pt *ProductTimes) GetExtent(start time.Time, stop time.Time) {

}*/

func GetTimes(host string, project string, product string) (ProductTimes, error) {
	pt := ProductTimes{}
	if host == "" {
		return pt, errors.New("Api host is empty")
	}

	log.Println("Getting times for", project, "/", product, "on endpoint", host)
	url := host + "/api/v1/projects/" + project + "/products/" + product + "/times"
	log.Println("Url:", url)
	err := apirequest.GetData(url, &pt)
	if err != nil {
		return pt, err
	}

	sort.SliceStable(pt.Times, func(i, j int) bool {
		return pt.Times[i].Before(pt.Times[j])
	})

	log.Println("Successfully got data for product: ", pt.Product)
	log.Println("Times: ", pt.Times)
	return pt, nil
}

type SQLData struct {
	SQLHost string
	SQLDB   string
	SQLUser string
	SQLPass string
}

func (s *SQLData) FillFromEnviron() {
	s.SQLHost = os.Getenv("SQLHOST")
	s.SQLDB = os.Getenv("SQLDB")
	s.SQLUser = os.Getenv("SQLUSER")
	s.SQLPass = os.Getenv("SQLPASS")
}

func (s *SQLData) Validate() error {
	if s.SQLHost == "" {
		return errors.New("Missing SQLHOST")
	}
	if s.SQLDB == "" {
		return errors.New("Missing SQLDB")
	}
	if s.SQLUser == "" {
		return errors.New("Missing SQLUSER")
	}
	if s.SQLPass == "" {
		return errors.New("Missing SQLPASS")
	}
	return nil
}

type MapLayerData struct {
	ProductTimes
	SQLData
	Project     string
	StartRange  time.Time
	EndRange    time.Time
	DefaultTime time.Time
}

func (self *MapLayerData) TimeRangeString() string {
	startRange := self.ProductTimes.Times[0]
	endRange := self.ProductTimes.Times[len(self.ProductTimes.Times)-1]
	//return self.StartRange.Format(time.RFC3339) + "/" + self.EndRange.Format(time.RFC3339)
	return startRange.Format(time.RFC3339) + "/" + endRange.Format(time.RFC3339)
}

func (self *MapLayerData) AllTimesString() string {
	strtimes := make([]string, len(self.Times))
	for idx, tmp := range self.Times {
		//strtimes =/usr/share/code/resources/app/out/vs/code/electron-sandbox/workbench/workbench.html append(strtimes, tmp.Format(time.RFC3339))
		strtimes[idx] = tmp.Format(time.RFC3339)
	}
	return strings.Join(strtimes, ",")
}

func (self *MapLayerData) DefaultTimeString() string {
	endRange := self.ProductTimes.Times[len(self.ProductTimes.Times)-1]
	return endRange.Format(time.RFC3339)
	//return self.DefaultTime.Format(time.RFC3339)
}

type ProductModified struct {
	Product   string `yaml:"product"`
	Timestamp string `yaml:"timestamp"`
	Project   string `yaml:"project"`
}
