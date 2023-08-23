package http

import (
	"net/http"

	"sdmht/lib/kitx"
	itfs "sdmht/sdmht/svc/interfaces"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewHTTPHandler(_ itfs.Service, _ *kitx.ServerOptions) http.Handler {
	// logger := opts.Logger()
	// options := []httptransport.ServerOption{
	// 	httptransport.ServerErrorEncoder(errorEncoder),
	// 	httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	// }
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())

	return r
}
