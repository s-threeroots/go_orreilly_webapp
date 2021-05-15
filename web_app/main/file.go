package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func readJson(filename string, obj interface{}) error {
	raw, err := ioutil.ReadFile(PROJECT_ROOT + filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	json.Unmarshal(raw, &obj)

	return err
}
