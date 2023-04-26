package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ergo-services/ergo/etf"
	"github.com/ergo-services/ergo/gen"
)

func createWebHandler() gen.WebHandlerBehavior {
	return &WebHandler{}
}

type WebHandler struct {
	gen.WebHandler
}

//
// Mandatory callbacks
//

// HandleRequest invokes on a HTTP-request
func (r *WebHandler) HandleRequest(process *gen.WebHandlerProcess, request gen.WebMessageRequest) gen.WebHandlerStatus {
	var data []BlockInfo
	if args, ok := request.Request.URL.Query()["crawl"]; ok {
		date, err := time.Parse(time.DateOnly, args[0])
		if err != nil {
			request.Response.WriteHeader(http.StatusBadRequest)
			request.Response.Write([]byte(err.Error()))
			return gen.WebHandlerStatusDone
		}
		// Got request for crawling data. Send this request to the Dispatcher
		process.Send("dispatcher", date.Unix())
		request.Response.Write([]byte("get data on " + date.String()))
		return gen.WebHandlerStatusDone
	}

	// No args, return data from the DB
	blockinfoDB.Find(&data)
	resp, err := json.Marshal(data)
	if err != nil {
		request.Response.WriteHeader(http.StatusInternalServerError)
		request.Response.Write([]byte(err.Error()))
		return gen.WebHandlerStatusDone
	}

	request.Response.Header().Set("Content-Type", "application/json")
	request.Response.Write(resp)

	return gen.WebHandlerStatusDone
}

//
// Optional gen.Server's callbacks
//

// HandleWebHandlerCall this callback is invoked on ServerProcess.Call(...).
func (r *WebHandler) HandleWebHandlerCall(process *gen.WebHandlerProcess, from gen.ServerFrom, message etf.Term) (etf.Term, gen.ServerStatus) {
	return nil, gen.ServerStatusOK
}

// HandleWebHandlerCast this callback is invoked on ServerProcess.Cast(...).
func (r *WebHandler) HandleWebHandlerCast(process *gen.WebHandlerProcess, message etf.Term) gen.ServerStatus {
	return gen.ServerStatusOK
}

// HandleWebHandlerInfo this callback is invoked on Process.Send(...).
func (r *WebHandler) HandleWebHandlerInfo(process *gen.WebHandlerProcess, message etf.Term) gen.ServerStatus {
	return gen.ServerStatusOK
}

// HandleWebHandlerTerminate this callback is invoked on the process termiation, providing the reason of termination
// along with the counter of handled requests
func (r *WebHandler) HandleWebHandlerTerminate(process *gen.WebHandlerProcess, reason string, count int64) {

}
