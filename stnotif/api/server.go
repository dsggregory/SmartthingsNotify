package api

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"code.dsg.com/smartthings_notif/stnotif/conf"
	"code.dsg.com/smartthings_notif/stnotif/dao"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type server struct {
	config *conf.Conf
	router *mux.Router
	db     *dao.DbHandle
}

// debug
func (s *server) dumpRoutes() {
	_ = s.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
}

// middleware to record the response status
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// logs the request before passing on to the mux router
func (s *server) wrapRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"RemoteAddr": r.RemoteAddr,
			"Method":     r.Method,
			"URL":        r.URL,
			"state":      "begin",
		}).Info()
		if s.config.AllowsHost(r.Host) {
			log.WithFields(log.Fields{
				"RemoteAddr": r.RemoteAddr,
			}).
				Error("refused")
		} else {
			// Initialize the status to 200 in case WriteHeader is not called
			rec := statusRecorder{w, 200}
			handler.ServeHTTP(w, r)
			log.WithFields(log.Fields{
				"RemoteAddr": r.RemoteAddr,
				"Method":     r.Method,
				"URL":        r.URL,
				"Status":     rec.status,
				"state":      "complete",
			}).Info()
		}
		// have to wrap the ResponseWriter if we want to log the status
	})
}

// StartServer starts the web server
func StartServer(config *conf.Conf) {
	db, err := dao.NewDbHandler(config)
	if err != nil {
		log.WithError(err).Fatalln("unable to open database")
	}

	s := server{config: config, router: mux.NewRouter(), db: db}
	s.initRoutes()
	if log.GetLevel() == log.DebugLevel {
		s.dumpRoutes()
	}

	svcPort := 8080
	if config.ServerPort > 0 {
		svcPort = config.ServerPort
	}
	laddr := fmt.Sprintf(":%d", svcPort)

	srv := &http.Server{
		Addr: laddr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.wrapRequest(s.router), // Pass our instance of gorilla/mux in
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Infof("listening on %s\n", laddr)
		if err := srv.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("Unable to start server")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, and SIGTERM (`kill(1)`).
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
