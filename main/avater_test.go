package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvater(t *testing.T) {

	var authAvatar AuthAvatar

	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)

	testChatUser := &chatUser{User: testUser}

	url, err := authAvatar.GetAvatarURL(testChatUser)

	if err != ErrNoAvatarURL {
		t.Error("authAvatar.GetAvatarURL should return ErrNoAvatarURL if AvatarURL is empty")
	}

	testUrl := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, ErrNoAvatarURL)
	url, err = authAvatar.GetAvatarURL(testChatUser)

	if err != nil {
		t.Error("authAvatar.GetAvatarURL should not ErrNoAvatarURL if AvatarURL is not empty")
	} else {
		if url != testUrl {
			t.Errorf("wrong url: %s", url)
		}
	}
	/*
		client := new(client)

		_, err := authAvatar.GetAvatarURL(client)
		if err != ErrNoAvatarURL {
			t.Error("値が存在しない場合、AuthAvatar.GetAvatarURLは、ErrNoAvotarURLを返すべき")
		}

		testUrl := "http://url-to-avatar/"

		client.userData = map[string]interface{}{"avatar_url": testUrl}

		url, err := authAvatar.GetAvatarURL(client)
		if err == ErrNoAvatarURL {
			t.Error("値が存在する場合、AuthAvatar.GetAvatarURLは、ErrNoAvotarURLを返すべきでない")
		} else {
			if url != testUrl {
				t.Error("AuthAvatar.GetAvatarURLは正しいURLを返すべき")
			}
		}
	*/
}

func TestGravatarAvatar(t *testing.T) {

	var gravatarAvatar GravatarAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("gravatarAvatar should not return err")
	}

	if url != "//www/gravatar/com/avatar/abc" {
		t.Errorf("wrong url: %s", url)
	}

	/*
		client := new(client)

		client.userData = map[string]interface{}{"email": "MyEmailAddress@example.com"}
		url, err := gravatarAvatar.GetAvatarURL(client)

		if err != nil {
			t.Error("gravatarAvatarはGetAvatarURLでエラーを返すべきではない")
		}
		if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
			t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
		}
	*/
}

func TestFileSystemAvater(t *testing.T) {

	filename := filepath.Join(PROJECT_ROOT+"avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("fileSystemAvatar should not return err.")
	}

	/*
		client := new(client)

		client.userData = map[string]interface{}{"userid": "abc"}

		url, err := fileSystemAvatar.GetAvatarURL(client)

		if err != nil {
			t.Error("FileSystemAvatar.GetAvatarURLはエラーを返すべきではない")
		}

	*/

	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatarが%sという誤った値を返した", url)
	}
}
