package core

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"html/template"
)

func Root(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// 捕获panic异常
		if err := recover(); err != nil{
			log.Println(err)
		}
	}()
	r.ParseForm() //解析参数, 默认是不会解析的
	if len(r.Form) > 0 && r.Form["url"] != nil && r.Form["token"] != nil {
		genShortURL(w, r)
	} else {
		getOriginalURL(w, r)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	if Conf.ErrorPage == "" {
		http.NotFound(w, r)
	} else {
		http.Redirect(w, r, Conf.ErrorPage, http.StatusFound)
	}
}


var t = func() *template.Template {
	var tmpl = `
<html>
  <body>
    <script type="text/javascript">
      window.location.href="{{ . }}"
      </script>
  </body>
</html>
`
	var t = template.New("layout.html")
	t, _ = t.Parse(tmpl)
	return t
}()

func getOriginalURL(w http.ResponseWriter, r *http.Request) {
	suffix := strings.TrimPrefix(r.URL.Path, `/`)
	if len(suffix) == 0 {
		NotFound(w, r)
	} else {
		originalURL, err := find(suffix)
		if err != nil {
			NotFound(w, r)
		} else if !strings.Contains(originalURL, "#") {
			http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
		}else {
			t.Execute(w, originalURL)
		}
	}
}

func genShortURL(w http.ResponseWriter, r *http.Request) {
	url := r.Form["url"]
	token := r.Form["token"]

	w.Header().Set("Content-type", "application/json")
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
