package main

import "testing"

func TestAuthAvater(t *testing.T) {
	var authAvatar AuthAvatar

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
}

func TestGravatarAvatar(t *testing.T) {

	var gravatarAvatar GravatarAvatar

	client := new(client)

	client.userData = map[string]interface{}{"email": "MyEmailAddress@example.com"}
	url, err := gravatarAvatar.GetAvatarURL(client)

	if err != nil {
		t.Error("gravatarAvatarはGetAvatarURLでエラーを返すべきではない")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
	}

}
