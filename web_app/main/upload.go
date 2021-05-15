package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, req *http.Request) {

	// hiddenのフォームに入ってる
	userId := req.FormValue("userid")

	// フォームからファイル持ってくる
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	avatarDir := PROJECT_ROOT + "avatars"
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		os.Mkdir(avatarDir, os.ModePerm)
	}

	filename := filepath.Join(avatarDir, userId+filepath.Ext(header.Filename))

	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "成功")
}
