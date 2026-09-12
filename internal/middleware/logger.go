package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ffajarpratama/pos-wash-api/pkg/util"
)

// color util
var (
	magenta = color.New(color.FgHiMagenta).SprintFunc()
	cyan    = color.New(color.FgHiCyan).SprintFunc()
	green   = color.New(color.FgHiGreen).SprintFunc()
	red     = color.New(color.FgHiRed).SprintFunc()
	yellow  = color.New(color.FgHiYellow).SprintFunc()
	blue    = color.New(color.FgHiBlue).SprintFunc()
)

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *customResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func Logger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := &customResponseWriter{ResponseWriter: w}
		t1 := time.Now()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[error:io.ReadAll()] \n%v\n", err)
		}

		defer func() {
			method := r.Method
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			addr := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
			elapsed := time.Since(t1).Abs().String()
			status := red(ww.statusCode)
			switch {
			case ww.statusCode < 200:
				status = cyan(ww.statusCode)
			case ww.statusCode < 300:
				status = green(ww.statusCode)
			case ww.statusCode < 400:
				status = blue(ww.statusCode)
			case ww.statusCode < 500:
				status = yellow(ww.statusCode)
			}

			log.Printf("%s %s - %s in %s", magenta(method), cyan(addr), status, green(elapsed))

			err = r.Body.Close()
			if err != nil {
				log.Printf("[error:body.Close()] \n%v\n", err)
			}

			token, _ := util.GetTokenFromHeader(r)
			claims := ParseWithoutVerified(token)
			if token != "" && claims != nil {
				log.Printf(`{"@auth":{"user_id":%s,"outlet_id:%s","role":%s}}`, claims.ID, claims.OutletID, claims.Role)
			}

			if len(b) > 0 && !strings.Contains(r.RequestURI, "media") {
				log.Printf(`{"@request":%s}`, string(b))
			}
		}()

		r.Body = io.NopCloser(bytes.NewBuffer(b))

		h.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}
