package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/annaworks/chatservice/pkg/api/handlers"
	Conf "github.com/annaworks/chatservice/pkg/conf"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Api struct {
	logger *zap.Logger
	router *mux.Router
	conf   Conf.Conf
}

func NewApi(logger *zap.Logger, conf Conf.Conf) Api {
	return Api{
		logger: logger,
		router: mux.NewRouter().StrictSlash(true),
		conf:   conf,
	}
}

func (a Api) Init() {
	err := a.load_routes()
	if err != nil {
		a.logger.Fatal("Could not load routes", zap.Error(err))
	}

	a.logger.Info("Loaded routes")
}

func (a Api) load_routes() error {
	v1 := a.router.PathPrefix("/api/v1").Subrouter()

	// health
	health_handler := handlers.Health_handler{
		Logger: a.logger.Named("health_handler"),
	}
	v1.HandleFunc("/health", health_handler.Ping).Methods(http.MethodGet)

	return nil
}

func (a Api) Serve() {
	server := http.Server{
		Addr:         fmt.Sprintf(":%s", a.conf.API_PORT),
		Handler:      a.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	a.logger.Info(fmt.Sprintf("Server serving on port %s", a.conf.API_PORT))

	log.Fatal(server.ListenAndServe())
}
