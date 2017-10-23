package core

import (
	"fmt"
	"net/http"
	"strings"
)

func GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数, 默认是不会解析的

	suffix := strings.TrimPrefix(r.URL.Path, `/`)
	if len(suffix) == 0 {
		w.WriteHeader(404)
		w.Write([]byte("<h1>404</h1>"))
	} else {
		originalURL, err := find(suffix)
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("<h1>404</h1>"))
		} else {
			http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
		}
	}
}

func GenShortURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数, 默认是不会解析的

	url := r.Form["url"]
	token := r.Form["token"]
	if len(url) == 0 || len(token) == 0 || token[0] != Conf.Token {
		fmt.Fprintf(w, `{"success":0}`)
		return
	}

	suffix, err := insert(url[0])
	if err != nil {
		fmt.Fprintf(w, `{"success":0}`)
		return
	}

	shortURL := fmt.Sprintf("%s%s", Conf.RootURL, suffix)
	fmt.Fprintf(w, `{"success":1,"url":"%s"}`, shortURL) //输出到客户端的信息
}
