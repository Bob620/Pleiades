package HttpServer

import (
	"net/http"
)

type Handler struct {
	DeleteRoutes []Route
	GetRoutes []Route
	PostRoutes []Route
	PutRoutes []Route
}

type Route struct {
	URL string
	Action func(w http.ResponseWriter, r http.Request)
}

func (handler Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		for i := range handler.DeleteRoutes {
			route := handler.DeleteRoutes[i]
			if (route.URL == req.URL.String()) {
			route.Action(res, *req)
			}
		}
	case http.MethodGet:
		for i := range handler.GetRoutes {
			route := handler.GetRoutes[i]
			if (route.URL == req.URL.String()) {
				route.Action(res, *req)
			}
		}
	case http.MethodPost:
		for i := range handler.PostRoutes {
			route := handler.PostRoutes[i]
			if (route.URL == req.URL.String()) {
				route.Action(res, *req)
			}
		}
	case http.MethodPut:
		for i := range handler.PutRoutes {
			route := handler.PutRoutes[i]
			if (route.URL == req.URL.String()) {
				route.Action(res, *req)
			}
		}
	}
}

func (handler *Handler) Delete(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.DeleteRoutes = append(handler.DeleteRoutes, Route{URL, f})
}

func (handler *Handler) Get(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.GetRoutes = append(handler.GetRoutes, Route{URL, f})
}

func (handler *Handler) Post(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.PostRoutes = append(handler.PostRoutes, Route{URL, f})
}

func (handler *Handler) Put(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.PutRoutes = append(handler.PutRoutes, Route{URL, f})
}