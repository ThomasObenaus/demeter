package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/thomasobenaus/sokar/api"
)

// # Notes:
// - Durations timeouts are optional, since the test itself has a defined deadline
// ## Procedure
// Sokar is at http://127.0.0.1:11000"
// 1. Send Request
// http.POST("http://127.0.0.1:11000/api/alerts",JSON{Alert{firing}})
// 2. Wait for Request to nomad, expect a certain body and respond with suitable data
// expect(timeout time.Duration).POST(data JSON).Return(code http.StatusCode, data JSON)

func Test_My(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mock := NewMockHTTP(mockCtrl)

	mock.EXPECT().POST("HELLO").Return(http.StatusOK, "huhu")

	receiver := api.New(18000)
	receiver.Run()
	// receiver.GET(/health,response)
	receiver.Router.HandlerFunc("GET", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		code, _ := mock.POST("HELLO")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		io.WriteString(w, "data")

		//receiver.Stop()
	}))
	time.Sleep(time.Millisecond * 100)

	res, err := http.Get("http://127.0.0.1:18000/health")
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", greeting)

	//receiver.Join()
	receiver.Stop()

}
