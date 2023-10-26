package apirequest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func GetData(requestURL string, data interface{}) error {
	if !strings.HasPrefix(requestURL, "http") {
		requestURL = "http://" + requestURL
	}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resBody, &data)
	return err
}
