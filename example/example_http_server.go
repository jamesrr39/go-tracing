//+build example

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	tracing "github.com/jamesrr39/go-tracing"
	"github.com/jamesrr39/goutil/must"
	"github.com/jamesrr39/goutil/streamtostorage"
)

func main() {

	// set up server
	addr := "localhost:9001"
	apiPath := "/api/v1/myroute"

	tempdir, err := ioutil.TempDir("", "")
	must.NoError(err)

	profilePath := filepath.Join(tempdir, "profile.pbf")
	log.Printf("writing profile to %s\n", profilePath)

	f, err := os.Create(profilePath)
	must.NoError(err)
	defer f.Close()

	profileWriter, err := streamtostorage.NewWriter(f, streamtostorage.MessageSizeBufferLenDefault)
	must.NoError(err)

	tracer := tracing.NewTracer(profileWriter)

	router := chi.NewRouter()
	router.Use(tracing.Middleware(tracer))
	router.Get(apiPath, exampleHandler)

	errChan := make(chan error)
	go func() {
		// start server
		err = http.ListenAndServe(addr, router)
		if err != nil {
			errChan <- err
		}
	}()

	doneChan := make(chan struct{})
	go func() {
		const maxAttempts = 20
		for i := 0; i < maxAttempts; i++ {
			time.Sleep(time.Second)
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s%s", addr, apiPath), nil)
			if err != nil {
				errChan <- err
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errChan <- err
				return
			}

			if resp.StatusCode == 200 {
				doneChan <- struct{}{}
				return
			}
		}
		errChan <- fmt.Errorf("reached max attempts: %d", maxAttempts)
	}()

	select {
	case err := <-errChan:
		log.Fatalln(err)
	case <-doneChan:
		log.Println("finished successfully")
	}
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	// handler. This won't do anything, just mark some events, and sleep between them so we can see the timeline in the profileviz summary
	ctx := r.Context()

	span, err := tracing.StartSpan(ctx, "print hello and wait 1 second")
	must.NoError(err)

	println("hello")

	time.Sleep(time.Second)

	span.End(ctx)

	span, err = tracing.StartSpan(ctx, "print world and wait 0.5 seconds")
	must.NoError(err)

	println("world")

	time.Sleep(time.Millisecond * 500)
	span.End(ctx)
}
