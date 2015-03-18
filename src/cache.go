package main

import (
	"container/list"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var count = make(map[string]int)
var cache = make(map[string]*list.Element)
var max = 1000
var ll = list.New()
var db *sql.DB

type ele struct {
	key   string
	value []int
}

func main() {
	Init()
	http.HandleFunc("/", do)
	http.ListenAndServe(":33333", nil)
}

func Init() {
	f, err := os.Open("/home/zhu/test/psd.in")
	if err != nil {
		panic(err)
	}
	psd, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	temp := string(psd)
	n := len(psd) - 1
	for temp[n] == '\n' {
		n--
	}
	defer f.Close()

	db, err = sql.Open("mysql", "root:"+string(psd[:n+1])+"@tcp(adserver1:3306)/adserver-staging")
	db.SetMaxOpenConns(2000000)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect!")
}

func do(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var s string
	var now = 0
	var body1 = make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		if string(b[i]) != " " && string(b[i]) != "\"" {
			body1[now] = b[i]
			now++
		}
	}
	body := string(body1[:now]) //"英孚4" and "英孚  4" are the same
	//    body:=string(body1)//"英孚4" and "英孚  4" are not the same

	idList, ok := cache[body]
	if !ok {
		rows, err := db.Query("SELECT bannerid from qtad_banners where description = ?", body)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		var id int
		var idlist ele
		for rows.Next() {
			ok = true
			rows.Scan(&id)
			idlist.value = append(idlist.value, id)
		}
		if ok {
			idlist.key = body
			e := ll.PushFront(&idlist)
			cache[body] = e
			count[body] = 0
			if ll.Len() > max {
				delete(cache, ll.Back().Value.(*ele).key)
				delete(count, ll.Back().Value.(*ele).key)
				ll.Remove(ll.Back())
			}
		}

	}
	idList, ok = cache[body]
	if ok {
		count[body]++
		ll.MoveToFront(idList)
		s = s + "bannerid: " + strconv.Itoa(cache[body].Value.(*ele).value[0])
		for i := 1; i < len(cache[body].Value.(*ele).value); i++ {
			s = s + ", " + strconv.Itoa(cache[body].Value.(*ele).value[i])
		}
		s = s + "  count:" + strconv.Itoa(count[body]) + "\n"

	} else {
		s = s + "Not Found \"" + body + "\"\n"
	}

	http.Error(w, s, 200)
}
