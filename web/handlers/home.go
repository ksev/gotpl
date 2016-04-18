package handlers

import "github.com/gorilla/mux"

// RegisterHome register all the handlers with the gorilla mux
func RegisterHome(n *mux.Router) {
	IHandleFunc(n, "/", RenderFunc("index"))
}
