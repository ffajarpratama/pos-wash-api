package app

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ffajarpratama/pos-wash-api/config"
	"github.com/ffajarpratama/pos-wash-api/internal/repository"
	"github.com/ffajarpratama/pos-wash-api/internal/router"
	"github.com/ffajarpratama/pos-wash-api/internal/usecase"
	"github.com/ffajarpratama/pos-wash-api/pkg/cloudinary"
	"github.com/ffajarpratama/pos-wash-api/pkg/postgres"
	"github.com/ffajarpratama/pos-wash-api/pkg/util"
)

func Exec() (err error) {
	cnf := config.New()
	rand.New(rand.NewSource(time.Now().UnixNano()))
	util.SetTimeZone("UTC")

	db, err := postgres.NewDBClient(cnf)
	if err != nil {
		return err
	}

	cld, err := cloudinary.NewClient(cnf)
	if err != nil {
		return err
	}

	repo := repository.New(db)
	uc := usecase.New(&usecase.Usecase{
		Cnf:      cnf,
		Repo:     repo,
		DB:       db,
		Uploader: cld,
	})

	handler := router.NewHTTPHandler(cnf, uc)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cnf.App.Port),
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
