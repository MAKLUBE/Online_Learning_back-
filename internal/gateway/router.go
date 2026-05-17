package gateway

//
//import "net/http"
//
//func NewRouter(clients *Clients) http.Handler {
//	mux := http.NewServeMux()
//
//	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
//		writeJSON(w, http.StatusOK, map[string]string{"message": "api gateway is running"})
//	})
//
//	registerAuthRoutes(mux, clients)
//	registerCourseRoutes(mux, clients)
//	registerLearningRoutes(mux, clients)
//	registerNotificationRoutes(mux, clients)
//
//	return withMiddleware(clients, mux)
//}
