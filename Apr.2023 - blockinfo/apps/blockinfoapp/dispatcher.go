package blockinfoapp

import (
	"fmt"

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

//
// Methods below are optional, so you can remove those that aren't be used
//

// HandleInfo invoked if this process received message sent with Process.Send(...).
func (s *Dispatcher) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	if date, ok := message.(int64); ok {
		s.log.Infof("generate crawl event with date %d", date)
		process.SendEventMessage(CrawlEvent, CrawlEventMessage{Date: date})
	}
	return gen.ServerStatusOK
}

// HandleCast invoked if this process received message sent with ServerProcess.Cast(...).
// Return ServerStatusStop to stop server with "normal" reason. Use ServerStatus(error)
// for the custom reason
func (s *Dispatcher) HandleCast(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	fmt.Printf("HandleCast: %#v \n", message)
	return gen.ServerStatusOK
}

// HandleCall invoked if this process got sync request using ServerProcess.Call(...)
func (s *Dispatcher) HandleCall(process *gen.ServerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	fmt.Printf("HandleCall: %#v \n", message)
	return nil, gen.ServerStatusOK
}

// HandleDirect invoked on a direct request made with Process.Direct(...)
func (s *Dispatcher) HandleDirect(process *gen.ServerProcess, ref etf.Ref, message interface{}) (interface{}, gen.DirectStatus) {
	fmt.Printf("HandleDirect: %#v \n", message)
	return nil, nil
}

// Terminate invoked on a termination process. ServerProcess.State is not locked during this callback.
func (s *Dispatcher) Terminate(process *gen.ServerProcess, reason string) {
	fmt.Printf("Terminated: %s with reason %s", process.Self(), reason)
}
