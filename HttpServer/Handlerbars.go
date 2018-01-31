package HttpServer

import (
	"io/ioutil"
)

func ReadFile(filename string) string {
	buf, _ := ioutil.ReadFile(filename)
	return string(buf)
}