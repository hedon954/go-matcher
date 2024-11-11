package apihttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hedon954/go-matcher/internal/api"
	internalapi "github.com/hedon954/go-matcher/internal/api"
	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/goapm/apm"
)

// Infra is the infrastructure for the http api.
type Infra struct {
	*api.InfraBase
}

// NewInfra creates a new Infra for the http api.
func NewInfra(
	sc config.Configer[config.ServerConfig],
	mc config.Configer[config.MatchConfig],
) *Infra {
	return &Infra{
		InfraBase: api.NewInfraBase(sc, mc),
	}
}

// Start starts the http api infrastructure, if any errors occur, it will panic.
func (i *Infra) Start() {
	i.InfraBase.Start()

	mapi, shutdown := internalapi.Start(i.SC, i.MC)
	i.AppendCloser(shutdown)

	server := API{mapi}
	r := server.setupRouter()
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", i.SC.Get().HTTPPort),
		Handler:           r.Handler(),
		ReadHeaderTimeout: time.Minute,
	}

	apm.Logger.Info(context.Background(), "http server started", map[string]any{
		"port": i.SC.Get().HTTPPort,
		"name": i.InfraBase.Name,
	})

	// TODO: add graceful shutdown
	if err := srv.ListenAndServe(); err != nil {
		apm.Logger.Error(context.Background(), "error occurs in http server", err, nil)
	}
}

// Stop stops the http api infrastructure.
func (i *Infra) Stop() {
	i.InfraBase.Stop()
}
