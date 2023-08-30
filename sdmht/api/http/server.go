package http

import (
	"net/http"

	"sdmht/lib/kitx"
	mw "sdmht/sdmht/api/http/middleware"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/gorilla/mux"
)

func NewHTTPHandler(_ itfs.Service, _ *kitx.ServerOptions) http.Handler {
	// logger := opts.Logger()
	// options := []httptransport.ServerOption{
	// 	httptransport.ServerErrorEncoder(errorEncoder),
	// 	httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	// }
	r := mux.NewRouter()
	// r.Handle("/metrics", promhttp.Handler())
	r.Use(mw.AccountMW())

	r.Handle("/", XXX()).Methods("POST")

	return r
}
