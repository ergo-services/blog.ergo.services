package blockinfoapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/sirupsen/logrus"

	log "github.com/sirupsen/logrus"
)

func createCrawlerAddress() gen.ServerBehavior {
	return &CrawlerAddress{}
}

const getAddressURI = "https://api.glassnode.com/v1/metrics/addresses/count?api_key=%s&a=%s&s=%d&u=%d"

type MessageAddressData struct {
	T int64
	V int64
}

type CrawlerAddress struct {
	gen.Server
	log *logrus.Entry
}

// Init invoked on a start this process.
func (s *CrawlerAddress) Init(process *gen.ServerProcess, args ...etf.Term) error {
	s.log = log.WithFields(log.Fields{
		process.Name(): process.Self(),
	})

	if err := process.MonitorEvent(CrawlEvent); err != nil {
		return err
	}

	s.log.Info("Crawler process is started")
	return nil
}

// HandleInfo invoked if this process received message sent with Process.Send(...).
func (s *CrawlerAddress) HandleInfo(process *gen.ServerProcess, message etf.Term) gen.ServerStatus {
	switch m := message.(type) {
	case CrawlEventMessage:
		var data []MessageAddressData
		s.log.Infof("got crawl request on %d", m.Date)

		// form the request to the API
		uri := fmt.Sprintf(getAddressURI, GLASSNODEAPIKEY, "btc", m.Date, m.Date+24*60*60)
		resp, err := http.Get(uri)
		if err != nil {
			s.log.Errorf("Can't receive address data: %s", err)
			return gen.ServerStatusOK
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &data); err != nil {
			s.log.Errorf("Can't unmarshal data: %s", err)
			return gen.ServerStatusOK
		}
		if len(data) == 0 {
			s.log.Warningf("no data on %d", m.Date)
			return gen.ServerStatusOK
		}

		// send it to the storage
		process.Send("storage", data[0])

	case gen.MessageEventDown:
		s.log.Warning("producer has terminated. stopping...")
		return gen.ServerStatusStop
	default:
		s.log.Errorf("unknown message", message)
	}
	return gen.ServerStatusOK
}
