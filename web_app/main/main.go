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

	"github.com/stretchr/gomniauth"
	_ "github.com/stretchr/gomniauth/providers/facebook"
	_ "github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var PROJECT_ROOT = "C:\\Users\\s.mine\\dev\\oreilly\\web_app\\"

var avatars = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	if err := t.templ.Execute(w, data); err != nil {
		log.Fatal("TemplateErr:", err)
	}

}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	var authInfo = AuthInfo{}

	err := readJson("client_secret.json", &authInfo)

	if err != nil {
		log.Fatal("read client secret err:", err)
	}

	gomniauth.SetSecurityKey(authInfo.SecretKey)

	gomniauth.WithProviders(
		google.New(authInfo.ClientSecret.ClientID, authInfo.ClientSecret.ClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// MustAuthで認証制御してる
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", logoutHandler)
	http.Handle("/upload", MustAuth(&templateHandler{filename: "upload.html"}))
	http.HandleFunc("/uploader", uploadHandler)
	http.Handle("/avatars/", http.FileServer(http.Dir(PROJECT_ROOT)))

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
