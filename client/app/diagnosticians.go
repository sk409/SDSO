package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func detectedVulnerability(request *http.Request, response *http.Response, r request, signature string) (bool, string, string, error) {
	request, err := r.toHTTPRequestWithSignature(signature)
	if err != nil {
		return false, "", "", err
	}
	requestBytes, err := httputil.DumpRequest(request, true)
	if err != nil {
		return false, "", "", err
	}
	requestString := string(requestBytes)
	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		return false, "", "", err
	}
	responseString := string(responseBytes)
	return true, requestString, responseString, nil
}

type diagnostician interface {
	diagnose(request) (bool, string, string, error)
}

type textMatchDiagnostician struct {
	signatures []string
	matches    []string
}

func newTextMatchDiagnostician(signatures, matches []string) *textMatchDiagnostician {
	return &textMatchDiagnostician{
		signatures: signatures,
		matches:    matches,
	}
}

func (t *textMatchDiagnostician) diagnose(r request) (bool, string, string, error) {
	client := new(http.Client)
	for index, signature := range t.signatures {
		request, err := r.toHTTPRequestWithSignature(signature)
		if err != nil {
			return false, "", "", err
		}
		response, err := client.Do(request)
		if err != nil {
			return false, "", "", err
		}
		defer response.Body.Close()
		responseBodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, "", "", err
		}
		responseBodyString := string(responseBodyBytes)
		if strings.Contains(responseBodyString, t.matches[index]) {
			response.Body = ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes))
			return detectedVulnerability(request, response, r, signature)
		}
	}
	return false, "", "", nil
}

type timeDiagnostician struct {
	signatures []string
	limit      float64
}

func newTimeDiagnostician(signatures []string, limit float64) *timeDiagnostician {
	return &timeDiagnostician{
		signatures: signatures,
		limit:      limit,
	}
}

func (t *timeDiagnostician) diagnose(r request) (bool, string, string, error) {
	client := new(http.Client)
	for _, signature := range t.signatures {
		request, err := r.toHTTPRequestWithSignature(signature)
		if err != nil {
			return false, "", "", err
		}
		start := time.Now()
		response, err := client.Do(request)
		end := time.Now()
		if err != nil {
			return false, "", "", err
		}
		defer response.Body.Close()
		diff := (end.Sub(start)).Seconds()
		if t.limit < float64(diff) {
			return detectedVulnerability(request, response, r, signature)
		}
	}

	return false, "", "", nil
}
