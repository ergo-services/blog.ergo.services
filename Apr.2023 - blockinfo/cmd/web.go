package main

import (
	"net/http"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"crypto/tls"

	"github.com/ergo-services/ergo/lib"
)

func createWeb() gen.WebBehavior {
	return &Web{}
}

type Web struct {
	gen.Web
	log *logrus.Entry
}

//
// Mandatory callbacks
//

// InitWeb invoked on starting Web server
func (w *Web) InitWeb(process *gen.WebProcess, args ...etf.Term) (gen.WebOptions, error) {
	var options gen.WebOptions
	options.Port = 9090
	options.Host = "localhost"
	// enable TLS with self-signed certificate
	cert, _ := lib.GenerateSelfSignedCert("Web Service")
	certUpdater := lib.CreateCertUpdater(cert)
	tlsConfig := &tls.Config{
		GetCertificate: certUpdater.GetCertificateFunc(),
	}
	options.TLS = tlsConfig

	mux := http.NewServeMux()
	handlerOptions := gen.WebHandlerOptions{
		NumHandlers:    3,
		IdleTimeout:    10,
		RequestTimeout: 20,
	}
	webRoot := process.StartWebHandler(createWebHandler(), handlerOptions)
	mux.Handle("/", webRoot)
	options.Handler = mux
	w.log = log.WithFields(log.Fields{
		process.Name(): process.Self(),
	})
	w.log.Infof("Starting Web server on https://%s:%d/\n", options.Host, options.Port)
	return options, nil
}

//
// Optional gen.Server's callbacks
//

// HandleWebCall this callback is invoked on ServerProcess.Call(...).
func (w *Web) HandleWebCall(process *gen.WebProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	return nil, gen.ServerStatusOK
}

// HandleWebCast this callback is invoked on ServerProcess.Cast(...).
func (w *Web) HandleWebCast(process *gen.WebProcess, message etf.Term) gen.ServerStatus {
	return gen.ServerStatusOK
}

// HandleWebInfo this callback is invoked on Process.Send(...).
func (w *Web) HandleWebInfo(process *gen.WebProcess, message etf.Term) gen.ServerStatus {
	return gen.ServerStatusOK
}
