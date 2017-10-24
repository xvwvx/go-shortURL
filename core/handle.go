package core

import (
	"fmt"
	"net/http"
	"strings"
)

func Root(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数, 默认是不会解析的
	if len(r.Form) > 0 {
		genShortURL(w, r)
	} else {
		getOriginalURL(w, r)
	}
}

func getOriginalURL(w http.ResponseWriter, r *http.Request) {
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

func genShortURL(w http.ResponseWriter, r *http.Request) {
	url := r.Form["url"]
	token := r.Form["token"]
	if len(url) == 0 || len(token) == 0 || token[0] != Conf.Token {
		fmt.Fprintf(w, `{"success":0}`)
		return
	}

	urlLong := url[0]
	suffix, err := insert(urlLong)
	if err != nil {
		fmt.Fprintf(w, `{"success":0}`)
		return
	}

	urlShort := fmt.Sprintf("%s%s", Conf.RootURL, suffix)
	fmt.Fprintf(w, `{"success":1,"result":{"url_short":"%s","url_long":"%s"}}`, urlShort, urlLong) //输出到客户端的信息
}
