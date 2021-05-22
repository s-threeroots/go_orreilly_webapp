package meander

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const PROJECT_ROOT = "C:\\Users\\s.mine\\dev\\oreilly\\meander\\"

func AuthData(filename string, obj interface{}) error {
	raw, err := ioutil.ReadFile(PROJECT_ROOT + filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	json.Unmarshal(raw, &obj)

	return err
}

type AuthInfo struct {
	ClientSecret ClientSecret `json:"installed"`
	SecretKey    string       `json:"secret_key"`
	APIKey       string       `json:"api_key"`
}

type ClientSecret struct {
	ClientID     string `json:"client_id"`
	ProjectID    string `json:"project_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
	ClientSecret string `json:"client_secret"`
}
