package main

import (
	"flag"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "エンドポイントのアドレス")
		mongo = flag.String("mongo", "chapter5_mongo_1", "mongoDBのアドレス")
	)

	flag.Parse()
	log.Println("Connectiong MongDB", mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("Failed connectiong MongoDB: ", err)
	}
	defer db.Close()
	/* 	mux := http.NewServeMux()
	   	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))

	   	log.Println("Starting Web server:", mux)
	   	graceful.Run(*addr, 1*time.Second, mux)
	   	log.Println("Quiting...")
	*/
	http.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	http.ListenAndServe(*addr, nil)
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "不正なAPIキーです")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
