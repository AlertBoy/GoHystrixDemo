package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	fmt.Println("serve started")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	http.HandleFunc("/", greet)

	fmt.Println("serve : ", ":9090")
	http.ListenAndServe(":9090", nil)

}

func greet(w http.ResponseWriter, r *http.Request) {

	i := rand.Intn(1000)
	fmt.Println("delay =", i)
	time.Sleep(time.Duration(i) * time.Millisecond)

	//fmt.Fprintf(w, "Hello World! %s", time.Now())
}
