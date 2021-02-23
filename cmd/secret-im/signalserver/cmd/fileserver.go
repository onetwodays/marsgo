package main

import (

	"net/http"

)

//处理函数,传入文件地址
func filehandler(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, filepath)
	}
}

func dirhandler(patern, filedir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(patern, http.FileServer(http.Dir(filedir)))
		handler.ServeHTTP(w, req)
	}
}


