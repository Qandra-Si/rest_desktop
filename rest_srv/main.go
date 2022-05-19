package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"rest_srv/internal/desktopstore"
	"time"
)

type desktopServer struct {
	store *desktopstore.DesktopStore
}

func NewDesktopServer() *desktopServer {
	store := desktopstore.New()
	return &desktopServer{store: store}
}

type RequestDesktop struct {
	ComputerName string    `json:"cname"`
	Ip           string    `json:"cip"`
	UserName     string    `json:"user"`
	At           time.Time `json:"at"`
}

func (ts *desktopServer) registerDesktopHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/register/" || req.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("expect method POST at /register/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	log.Printf("handling desktop create at %s\n", req.URL.Path)

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rd RequestDesktop
	if err := dec.Decode(&rd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateDesktop(rd.ComputerName, rd.UserName, rd.Ip, rd.At)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *desktopServer) unregisterDesktopHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/unregister/" || req.Method != http.MethodDelete {
		http.Error(w, fmt.Sprintf("expect method DELETE at /unregister/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	log.Printf("handling delete desktop at %s\n", req.URL.Path)

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rd RequestDesktop
	if err1 := dec.Decode(&rd); err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	err2 := ts.store.DeleteDesktop(rd.ComputerName, rd.UserName, rd.Ip, rd.At)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusNotFound)
	}
}

func (ts *desktopServer) updateDesktopHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/update/" || req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET at /update/, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	log.Printf("handling update desktop at %s\n", req.URL.Path)

	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rd RequestDesktop
	if err1 := dec.Decode(&rd); err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	err2 := ts.store.UpdateDesktop(rd.ComputerName, rd.UserName, rd.Ip, rd.At)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	server := NewDesktopServer()
	mux.HandleFunc("/register/", server.registerDesktopHandler)
	mux.HandleFunc("/unregister/", server.unregisterDesktopHandler)
	mux.HandleFunc("/update/", server.updateDesktopHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}
