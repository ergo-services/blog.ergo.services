package blockinfoapp

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func createCommonSup() gen.SupervisorBehavior {
	return &CommonSup{}
}

type CommonSup struct {
	gen.Supervisor
}

func (sup *CommonSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: "commonsup",
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "dispatcher",
				Child: createDispatcher(),
			},
			gen.SupervisorChildSpec{
				Name:  "crawlersup",
				Child: createCrawlerSup(),
			},
		},
		Strategy: gen.SupervisorStrategy{
			Type:      gen.SupervisorStrategyOneForAll,
			Intensity: 2,
			Period:    5,
			Restart:   gen.SupervisorStrategyRestartTemporary,
		},
	}
	return spec, nil
}
