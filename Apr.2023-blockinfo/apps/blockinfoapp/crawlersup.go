package blockinfoapp

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func createCrawlerSup() gen.SupervisorBehavior {
	return &CrawlerSup{}
}

type CrawlerSup struct {
	gen.Supervisor
}

func (sup *CrawlerSup) Init(args ...etf.Term) (gen.SupervisorSpec, error) {
	spec := gen.SupervisorSpec{
		Name: "crawlersup",
		Children: []gen.SupervisorChildSpec{
			gen.SupervisorChildSpec{
				Name:  "crawlerblock",
				Child: createCrawlerBlock(),
			},
			gen.SupervisorChildSpec{
				Name:  "crawleraddress",
				Child: createCrawlerAddress(),
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
