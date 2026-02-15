package api

import "net/http"

// corsMiddleware adds CORS headers to handle cross-origin requests
func (a *API) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *API) routes() {
	a.mux = http.NewServeMux()
	a.mux.HandleFunc("POST /v1/config", a.handlerApplyConfig)
	a.mux.HandleFunc("GET /v1/producers/{service_id}/{topic_id}", a.handlerGetProducerDetail)
	a.mux.HandleFunc("GET /v1/producers", a.handlerListProducers)
	a.mux.HandleFunc("GET /v1/events", a.handlerListEvents)
	a.mux.HandleFunc("GET /v1/consumers", a.handlerListConsumers)
	a.mux.HandleFunc("GET /v1/overview", a.handlerOverview)
	a.mux.HandleFunc("GET /v1/topics/{name}", a.handlerGetTopic)
	a.mux.HandleFunc("GET /v1/topics", a.handlerListTopics)
	a.mux.HandleFunc("GET /v1/services/{name}", a.handlerGetService)
}
