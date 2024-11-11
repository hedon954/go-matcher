package api

import (
	"time"

	"github.com/hedon954/go-matcher/internal/config"
	"github.com/hedon954/goapm"
	"github.com/hedon954/goapm/apm"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// InfraBase is the base infrastructure for the api.
type InfraBase struct {
	*goapm.Infra

	SC config.Configer[config.ServerConfig]
	MC config.Configer[config.MatchConfig]
}

// NewInfraBase creates a new InfraBase.
func NewInfraBase(
	sc config.Configer[config.ServerConfig],
	mc config.Configer[config.MatchConfig],
) *InfraBase {
	infra := &InfraBase{
		SC: sc,
		MC: mc,
	}
	infra.Infra = goapm.NewInfra("go-matcher",
		goapm.WithAPM(sc.Get().OtelExporterEndpoint),
		goapm.WithMetrics(),
		goapm.WithRotateLog("logs/go-matcher.log",
			rotatelogs.WithRotationTime(time.Hour*24),
			rotatelogs.WithRotationCount(30),
		),
		goapm.WithAutoPProf(&apm.AutoPProfOpt{
			EnableCPU:       true,
			EnableMem:       true,
			EnableGoroutine: true,
		}),
	)
	return infra
}

// Start starts the infrastructure, if any errors occur, it will panic.
func (i *InfraBase) Start() {}

// Stop stops the infrastructure.
func (i *InfraBase) Stop() {
	i.Infra.Stop()
}
