package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/tosuke/yukkuri-server/aqtalk"
)

func handler() http.Handler {
	sf := &aqtalk.SynthFactory{}

	mux := http.NewServeMux()
	mux.HandleFunc("/koe.wav", func (w http.ResponseWriter, r *http.Request) {
		if(r.Method != http.MethodGet) {
			w.WriteHeader(404)
			return
		}

		typ := "f1"
		if v := r.URL.Query().Get("type"); v != "" {
			typ = v
		}

		speed := 100
		if v := r.URL.Query().Get("speed"); v != "" {
			s, err := strconv.Atoi(v);
			if err != nil {
				panic(fmt.Errorf("invalid speed"))
			}
			speed = s
		}

		koe := "ゆっくりしていってね"
		if v := r.URL.Query().Get("koe"); v != "" {
			koe = v
		}

		synth, err := sf.Get(typ)
		if err != nil {
			panic(err)
		}
		
		data, err := synth.Synthe(koe, uint32(speed))
		if err != nil {
			panic(err)
		} 

		w.Header().Add("content-type", "audio/wav")
		w.WriteHeader(200)
		w.Write(data)
	})

	return mux	
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	port := 3000
	if v, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		port = v
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: handler(),
	}

	// Graceful Shutdown
	go func() {
		<- ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
	}()

	fmt.Printf("Listen :%d......\n", port)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("server closed with error: ", err)
	}
}
