package blockinfoapp

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	log "github.com/sirupsen/logrus"
)

const GLASSNODEAPIKEY = "xxxxxxxxxxxxxxxxxxxxxxxxxxx" // replace by your API-KEY

func CreateBlockInfoApp() gen.ApplicationBehavior {
	return &BlockInfoApp{}
}

type BlockInfoApp struct {
	gen.Application
}

func (app *BlockInfoApp) Load(args ...etf.Term) (gen.ApplicationSpec, error) {
	return gen.ApplicationSpec{
		Name:        "blockinfoapp",
		Description: "description of this application",
		Version:     "v.1.0",
		Children: []gen.ApplicationChildSpec{
			gen.ApplicationChildSpec{
				Name:  "commonsup",
				Child: createCommonSup(),
			},
		},
	}, nil
}

func (app *BlockInfoApp) Start(process gen.Process, args ...etf.Term) {
	log.Infof("Application BlockInfoApp started with Pid %s\n", process.Self())
}
