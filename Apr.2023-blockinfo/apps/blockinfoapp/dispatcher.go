package blockinfoapp

import (
	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func createDispatcher() gen.ServerBehavior {
	return &Dispatcher{}
}

type Dispatcher struct {
	gen.Server
	log *logrus.Entry
}

const CrawlEvent gen.Event = "crawl"

type CrawlEventMessage struct {
	Date int64
}

// Init invoked on a start this process.
func (s *Dispatcher) Init(process *gen.ServerProcess, args ...etf.Term) error {
	s.log = log.WithFields(log.Fields{
		process.Name(): process.Self(),
	})

	if err := process.RegisterEvent(CrawlEvent, CrawlEventMessage{}); err != nil {
		return err
	}

	s.log.Info("Dispatcher process is started")
	return nil
}

// HandleInfo invoked if this process received message sent with Process.Send(...).
func (s *Dispatcher) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	if date, ok := message.(int64); ok {
		s.log.Infof("generate crawl event with date %d", date)
		process.SendEventMessage(CrawlEvent, CrawlEventMessage{Date: date})
	}
	return gen.ServerStatusOK
}
