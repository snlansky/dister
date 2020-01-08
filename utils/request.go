package utils

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func HttpGet(urlStr string) (string, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return "", errors.WithMessage(err, "http get error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "read response body error")
	}
	return string(body), err
}
