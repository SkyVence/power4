package main

import (
	"net/http"
	"power4/base/handlers"
	bonusHandlers "power4/bonus/handlers"
	"power4/shared"
)

var routes = []shared.Route{
	{
		Method:  "GET",
		Path:    "/health",
		Handler: func(w http.ResponseWriter, r *http.Request) { http.Error(w, "OK", http.StatusOK) },
	},
	{
		Method:  "GET",
		Path:    "/",
		Handler: handlers.HomeHandler,
	},
	{
		Method:  "POST",
		Path:    "/move",
		Handler: handlers.MoveHandler,
	},
	{
		Method:  "POST",
		Path:    "/new-game",
		Handler: handlers.NewGameHandler,
	},
	{
		Method:  "POST",
		Path:    "/reset-scores",
		Handler: handlers.ResetScoresHandler,
	},
	{
		Method:  "GET",
		Path:    "/bonus/setup",
		Handler: bonusHandlers.SetupHandler,
	},
	{
		Method:  "POST",
		Path:    "/bonus/start-game",
		Handler: bonusHandlers.StartGameHandler,
	},
	{
		Method:  "GET",
		Path:    "/bonus/game",
		Handler: bonusHandlers.GameHandler,
	},
	{
		Method:  "POST",
		Path:    "/bonus/move",
		Handler: bonusHandlers.MakeMove,
	},
	{
		Method:  "POST",
		Path:    "/bonus/new-game",
		Handler: bonusHandlers.NewGameHandler,
	},
	{
		Method:  "POST",
		Path:    "/bonus/reset-scores",
		Handler: bonusHandlers.ResetScoresHandler,
	},
	// Redirect root to setup
	{
		Method: "GET",
		Path:   "/bonus",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/bonus/setup", http.StatusSeeOther)
		},
	},
}

func main() {

	shared.StartServer(routes, "0.0.0.0:8080")
}
