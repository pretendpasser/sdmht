package http

// func NmsAlertData(svc *erranalysis.Service, options []httptransport.ServerOption, mws ...endpoint.Middleware) http.Handler {
// 	ep := api.MakeNmsAlertDataEndpoint(svc)
// 	for i := len(mws) - 1; i >= 0; i-- { // reverse
// 		ep = mws[i](ep)
// 	}
// 	return httptransport.NewServer(
// 		ep,
// 		decodeNmsAlertDataReq,
// 		encodeResponse,
// 		options...,
// 	)
// }
