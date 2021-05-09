package main

import (
	"flag"
	"log"
	"net/http"
	"oreilly/trace"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

var PROJECT_ROOT = "C:\\Users\\s.mine\\dev\\oreilly\\"

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ =
			template.Must(template.ParseFiles(filepath.Join(PROJECT_ROOT+"templates", t.filename)))
	})
	if err := t.templ.Execute(w, r); err != nil {
		log.Fatal("TemplateErr:", err)
	}

}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()
	log.Println("Webサーバーを開始します。ポート: ", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServ:", err)
	}
}

/* 1.1.1.2
func main() {

	http.Handle("/", &templateHandler{filename: "chat.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
*/

/*
1.1
func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<html>
			<head>
				<title>chat</title>
			</head>
			<body>
				let's chat!
			</body>
		</html>
		`))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
*/
