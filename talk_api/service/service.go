package service

import "net/http"

type SocketService interface {
	AddClient(w http.ResponseWriter, req *http.Request)
	Run()
}
