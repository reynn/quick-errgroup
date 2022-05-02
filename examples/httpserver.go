package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/reynn/quickerr"
)

func main() {
	group, _ := quickerr.New(context.Background())

	group.Go(func() error {
		defer func() {
			fmt.Println("HTTP server closing")
		}()

		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		fmt.Println("HTTP Listening on port 8080")
		return http.ListenAndServe(":8080", nil)
	})

	spawn(15, group)

	group.Go(func() error {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Kill, os.Interrupt)

		sig := <-ch
		return fmt.Errorf("received %s signal", sig)
	})

	fmt.Println("Waiting for routines to complete")
	if err := group.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "App failed: %v\n", err)
	}

	fmt.Println("-----COMPLETE------")
}

func spawn(i int, group *quickerr.Group) {
	for j := 0; j < i; j++ {
		num := j
		group.Go(func() error {
			defer func(h int) {
				fmt.Printf("Routine %d complete\n", h)
			}(num)
			time.Sleep(30 * time.Second)
			if num == 5 {
				return fmt.Errorf("random err")
			}
			return nil
		})
	}
}
