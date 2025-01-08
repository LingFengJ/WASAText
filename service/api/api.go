package api

import (
	"errors"
	"github.com/LingFengJ/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Logger   *logrus.Logger
	Database database.AppDatabase
}

type _router struct {
	router     *httprouter.Router
	baseLogger *logrus.Logger
	db         database.AppDatabase
}

// // New returns a new instance of the router
// func NewRouter(logger *logrus.Logger) *_router {
//     return &_router{
//         router:     httprouter.New(),
//         baseLogger: logger,
//     }
// }

func New(cfg Config) (*_router, error) {
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if cfg.Database == nil {
		return nil, errors.New("database is required")
	}

	rt := &_router{
		router:     httprouter.New(),
		baseLogger: cfg.Logger,
		db:         cfg.Database,
	}

	if err := cfg.Database.Ping(); err != nil {
		return nil, errors.New("could not connect to database")
	}

	return rt, nil
}
