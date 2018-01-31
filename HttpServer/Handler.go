package HttpServer

import (
	"net/http"
	"strings"
)

type Handler struct {
	Routes []Route
	RouteNotFound Route
	ServerError Route
}

type Route struct {
	URL string
	Method string
	Routes []Route
	Action func(w http.ResponseWriter, r http.Request)
}

func (handler Handler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "/" {
		req.RequestURI = ""
	}
	req.ParseForm()
	requestArray := strings.Split(strings.Split(req.RequestURI, "?")[0], "/")

	k := 0
	var currentRoute Route
	currentRouteList := handler.Routes

	for k < len(requestArray) {
		requestRoute := requestArray[k]
		var broke = false
		for i := range currentRouteList {
			route := currentRouteList[i]
			if route.URL == requestRoute {
				currentRoute = route
				currentRouteList = route.Routes
				broke = true
				break
			}
		}
		if !broke {
			currentRoute.Action(res, *req)
			return
		}
		k++
	}

	currentRoute.Action(res, *req)
}

func (handler *Handler) AddError(errorType string, f func(res http.ResponseWriter, req http.Request)) {
	switch errorType {
	case "404":
		handler.RouteNotFound = Route{"404", "GET", []Route{}, f}
	case "500":
		handler.RouteNotFound = Route{"500", "GET", []Route{}, f}
	}
}
/*
func (handler *Handler) Delete(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.DeleteRoutes = append(handler.DeleteRoutes, Route{URL, f})
}
*/
func (handler *Handler) Get(URL string, f func(res http.ResponseWriter, req http.Request)) {
	if URL == "/" {
		URL = ""
	}
	requestArray := strings.Split(URL, "/")

	k := 0
	var currentRoute Route
	currentRouteList := handler.Routes

	for k < len(requestArray) {
		requestRoute := requestArray[k]
		var broke = false
		for i := range currentRouteList {
			route := currentRouteList[i]
			if route.URL == requestRoute {
				currentRoute = route
				currentRouteList = route.Routes
				broke = true
				break
			}
		}
		if !broke {
			var route Route;
			if k+1 == len(requestArray) {
				route = Route{currentRoute.URL, "GET", []Route{}, f}
			} else {
				route = Route{requestRoute, "GET", []Route{}, func(res http.ResponseWriter, req http.Request) {
					handler.RouteNotFound.Action(res, req)
				}}
			}
			if currentRoute.Action == nil {
				handler.Routes = append(handler.Routes, route)
			} else {
				currentRoute.Routes = append(currentRoute.Routes, route)
			}
			currentRoute = route
			currentRouteList = route.Routes
		}
		k++
	}
}
/*
func (handler *Handler) Post(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.PostRoutes = append(handler.PostRoutes, Route{URL, f})
}

func (handler *Handler) Put(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.PutRoutes = append(handler.PutRoutes, Route{URL, f})
}
*/