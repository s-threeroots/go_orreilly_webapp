package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

var providers = []string{"google", "github"}

type authHandler struct {
	next http.Handler
}

type AuthInfo struct {
	ClientSecret ClientSecret `json:"installed"`
	SecretKey    string       `json:"secret_key"`
}

type ClientSecret struct {
	ClientID     string `json:"client_id"`
	ProjectID    string `json:"project_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
	ClientSecret string `json:"client_secret"`
}

// authHandler用のServeHTTP
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// クッキーみて、認証情報がなかったらリダイレクト
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)

	} else if err != nil {
		// error
		panic(err.Error())

	} else {
		// success
		h.next.ServeHTTP(w, r)

	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	if len(segs) != 4 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "invalid url")
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":

		if arrayContains(providers, provider) {

			provider, err := gomniauth.Provider(provider)
			if err != nil {
				log.Fatalln("認証プロバイダーの取得に失敗", provider, "-", err)
			}
			loginUrl, err := provider.GetBeginAuthURL(nil, nil)
			if err != nil {
				log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました", provider, "-", err)
			}

			w.Header().Set("Location", loginUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "invalid url")
			return
		}

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗", provider, "-", err)
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした", provider, "-", err)
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザーの取得に失敗", provider, "-", err)
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/"})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sに非対応です", action)
	}
}

func arrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
