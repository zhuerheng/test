package main

import (
	"encoding/json"
	_ "fmt"
	_ "io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func testHandle(w http.ResponseWriter, r *http.Request) {
	type Time struct {
		Year   int
		Mon    int
		Day    int
		Hour   int
		Minu   int
		Second int
		Zone   int
	}
	var s Time
	now := time.Now()

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	var ok bool
	var zone []string

	if r.Method == "GET" {
		zone, ok = r.Form["zone"]
	} else if r.Method == "POST" {
		zone, ok = r.PostForm["zone"]
	}

	if ok {
		z, _ := strconv.Atoi(zone[0])
		loc := time.FixedZone("a", z*3600)
		now = now.In(loc)
		temp, _ := strconv.Atoi(now.Format("1"))
		s = Time{
			Year:   now.Year(),
			Mon:    temp,
			Day:    now.Day(),
			Hour:   now.Hour(),
			Minu:   now.Minute(),
			Second: now.Second(),
			Zone:   z,
		}
		out, _ := json.Marshal(s)
		w.Write(out)
	} else {
		http.Error(w, "缺少参数zone", http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/", testHandle)
	panic(http.ListenAndServe(":22222", nil))
}
