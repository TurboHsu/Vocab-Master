package automatic

import (
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

// Do a post request using stored cookies and headers
func doPOST(url string, data string) string {
	req, _ := http.NewRequest("POST", url, strings.NewReader(data))
	//Add headers and cookies
	req.Header = dataset.RequestInfo.Header
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprint(len(data)))
	req.Header.Set("Origin", "https://app.vocabgo.com")

	for i := 0; i < len(dataset.RequestInfo.Cookies); i++ {
		req.AddCookie(dataset.RequestInfo.Cookies[i])
	}
	//Do request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[E]" + err.Error())
	}
	defer response.Body.Close()
	read, _, _ := switchContentEncoding(response)
	raw, _ := io.ReadAll(read)
	return string(raw)
}

// Do a get request using stored cookies and headers
func doGET(url string) string {
	req, _ := http.NewRequest("GET", url, nil)
	//Add headers and cookies
	req.Header = dataset.RequestInfo.Header
	for i := 0; i < len(dataset.RequestInfo.Cookies); i++ {
		req.AddCookie(dataset.RequestInfo.Cookies[i])
	}
	//Do request
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[E]" + err.Error())
	}
	defer response.Body.Close()
	read, _, _ := switchContentEncoding(response)
	raw, _ := io.ReadAll(read)
	return string(raw)
}

func switchContentEncoding(res *http.Response) (bodyReader io.Reader, encoder string, err error) {
	encoder = res.Header.Get("Content-Encoding")
	switch encoder {
	case "gzip":
		bodyReader, err = gzip.NewReader(res.Body)
	case "deflate":
		bodyReader = flate.NewReader(res.Body)
	case "br":
		bodyReader = brotli.NewReader(res.Body)
	default:
		bodyReader = res.Body
	}
	return
}

// Splits the salt and the valid JSON string
func splitSalt(raw string) (salt string, validJSON string) {
	salt = raw[0:32]
	validJSON = raw[32:]
	return
}
