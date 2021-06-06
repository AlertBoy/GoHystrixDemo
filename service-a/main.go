package main

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	hystrix.DefaultVolumeThreshold = 3
	hystrix.DefaultErrorPercentThreshold = 75
	hystrix.DefaultTimeout = 500
	hystrix.DefaultSleepWindow = 3500

	streamHandler := hystrix.NewStreamHandler()
	streamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8081"), streamHandler)

	http.HandleFunc("/", HandelSubSystem)
	http.ListenAndServe(":8080", nil)
}

func HandelSubSystem(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("hi"))
	result := make(chan []byte)
	hystrix.Go("my_command", func() error {
		resp, err := http.Get("http://localhost:9090")
		if err != nil {
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		result <- all
		return nil
	}, func(err error) error {
		result <- []byte("error")
		return nil
	})
	ret := <-result
	w.Write(ret)
}
