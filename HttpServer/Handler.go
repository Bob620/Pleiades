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
	route := handler.searchURL(req.RequestURI)

	route.Action(res, *req)
}

func (handler *Handler) AddError(errorType string, f func(res http.ResponseWriter, req http.Request)) {
	switch errorType {
	case "404":
		handler.RouteNotFound = Route{"404", http.MethodGet, []Route{}, f}
	case "500":
		handler.RouteNotFound = Route{"500", http.MethodGet, []Route{}, f}
	}
}

func (handler *Handler) Delete(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.addURL(URL, http.MethodDelete, f)
}

func (handler *Handler) Get(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.addURL(URL, http.MethodGet, f)
}

func (handler *Handler) Post(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.addURL(URL, http.MethodPost, f)
}

func (handler *Handler) Put(URL string, f func(res http.ResponseWriter, req http.Request)) {
	handler.addURL(URL, http.MethodPut, f)
}

func (handler Handler) searchURL(URL string) Route {
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
			return currentRoute
		}
		k++
	}

	return currentRoute
}

func (handler *Handler) addURL(URL string, method string, f func(res http.ResponseWriter, req http.Request)) {
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
			var route Route
			if k+1 == len(requestArray) {
				route = Route{currentRoute.URL, method, []Route{}, f}
				route = Route{requestRoute, method, []Route{}, func(res http.ResponseWriter, req http.Request) {
					handler.RouteNotFound.Action(res, req)
				}}
			} else {
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