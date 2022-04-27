package api

import "net/http"

type StaticHandler struct {
	http.Handler
}

func NewStaticHander(dir string) *StaticHandler {
	return &StaticHandler{http.FileServer(http.Dir(dir))}
}
