package http_server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ffajarpratama/pos-wash-api/config"
)

func NewHTTPServer(addr string, handler http.Handler, cnf *config.Config) (err error) {
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 90 * time.Second,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// implement graceful shutdown
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Printf("[graceful-shutdown-time-out] \n%v\n", err.Error())
			}
		}()

		defer cancel()

		log.Println("graceful shutdown.....")

		// trigger graceful shutdown
		err = httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("[graceful-shutdown-error] \n%v\n", err.Error())
		}

		serverStopCtx()
	}()

	// run server
	log.Printf("[http-server-online] %v\n", cnf.App.URL)
	err = httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("[http-server-failed] \n%v\n", err.Error())
		return err
	}

	<-serverCtx.Done()

	return
}
