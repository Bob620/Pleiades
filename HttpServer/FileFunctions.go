package HttpServer

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func ReadFile(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "404"
	}
	return string(buf)
}

func ServeDir(dirname string, URL string) func(res http.ResponseWriter, req http.Request) {
	return func(res http.ResponseWriter, req http.Request) {
		requestArray := strings.Split(strings.Split(req.RequestURI, "?")[0], "/")
		filename := strings.Join(requestArray[find(URL, requestArray)+1:],"/")

		res.Write([]byte(ReadFile(dirname+"/"+filename)))
	}
}

func find(value string, slice []string) int {
	for p, v := range slice {
		if (v == value) {
			return p
		}
	}
	return -1
}