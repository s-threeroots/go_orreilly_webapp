package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

const otherWord = "*"

type inData struct {
	Transforms []string `json:"transforms"`
}

/*
var transforms = []string{
	otherWord,
	otherWord,
	otherWord,
	otherWord,
	otherWord + "app",
	otherWord + "site",
	otherWord + "time",
	"get" + otherWord,
	"go" + otherWord,
	"lets" + otherWord,
}
*/
func main() {

	infilepath := flag.String("infile", "", "infile path")
	flag.Parse()

	data, err := ioutil.ReadFile(*infilepath)
	if err != nil {
		log.Fatal("read infile failed")
	}

	var readdata = inData{}

	err = json.Unmarshal(data, &readdata)
	if err != nil {
		log.Fatal("json parse failed")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := readdata.Transforms[rand.Intn(len(readdata.Transforms))]
		fmt.Println(strings.Replace(t, otherWord, s.Text(), -1))
	}
}
