package fetch

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetData(requestURL string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(resBody, &data)
	return nil
}
