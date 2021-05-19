package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nsq "github.com/nsqio/go-nsq"
	mgo "gopkg.in/mgo.v2"
)

func main() {

	var stoplock sync.Mutex
	stop := false
	stopChan := make(chan struct{}, 1)
	signalChan := make(chan os.Signal, 1)

	go func() {
		<-signalChan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Print("停止します...")
		stopChan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := dialdb(); err != nil {
		log.Fatalln("MongoDBへのダイアルに失敗しました: ", err)
	}
	defer closedb()

	votes := make(chan string)
	publisherStoppedChan := publishVotes(votes)
	twitterStopedChan := startTwitterStream(stopChan, votes)

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<-twitterStopedChan
	close(votes)
	<-publisherStoppedChan
}

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("Dialing to MongoDB: chapter5_mongo_1")

	db, err = mgo.Dial("chapter5_mongo_1")
	return err
}

func closedb() {
	db.Close()
	log.Println("DB was closed.")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	defer iter.Close()
	var p poll
	for iter.Next(&p) {
		options = append(options, p.Options...)
	}

	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("chapter5_nsqd_1:4150", nsq.NewConfig())

	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote))
		}
		log.Println("Publisher: stopping.")
		pub.Stop()
		log.Println("Publisher: stopped.")
		stopchan <- struct{}{}
	}()
	return stopchan
}
