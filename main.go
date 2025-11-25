package main

import (
	"net/http"
	"power4/shared"
)

func main() {

	routes := []shared.Route{
		{
			Method:  "GET",
			Path:    "/health",
			Handler: func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) },
		},
	}

	shared.StartServer(routes, "127.0.0.1:80")
}
