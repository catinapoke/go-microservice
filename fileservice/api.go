package fileservice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type ApiServer struct {
	svc FileService
}

func CreateAPIServer(svc FileService) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	http.HandleFunc("/get", s.Get)
	http.HandleFunc("/set", s.Set)
	http.HandleFunc("/delete", s.Delete)
	return http.ListenAndServe(listenAddr, nil)
}

func DoNothing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "wat")
}

func (s *ApiServer) Get(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_str)

	if err != nil {
		http.Error(w, fmt.Sprintf("Can't convert id '%s' to int!", id_str), http.StatusNotFound)
		return
	}

	path, err := s.svc.Get(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Can't get file with id '%d'!", id), http.StatusNotFound)
		return
	}

	// TODO: Change filename by changing r.URL
	http.ServeFile(w, r, path)
}

// Modified of https://stackoverflow.com/questions/45541656/golang-send-file-via-post-request
func upload(w http.ResponseWriter, r *http.Request) (string, error) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		//panic(err) //dont do this
		return "", err
	}
	defer file.Close()

	// Generate filename
	filename := randStringRunes(10) + "_" + handler.Filename

	// copy example
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	io.Copy(f, file)

	return filename, nil
}

func (s *ApiServer) Set(w http.ResponseWriter, r *http.Request) {
	path, err := upload(w, r)

	if err != nil {
		http.Error(w, "Got error while uploading file!", http.StatusBadRequest)
		return
	}

	id, err := s.svc.Set(path)

	if err != nil {
		http.Error(w, "Got error while saving file!", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, id)
}

func (s *ApiServer) Delete(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_str)

	if err != nil {
		http.Error(w, "Can't convert id to int!", http.StatusNotFound)
		return
	}

	err = s.svc.Delete(id)

	if err != nil {
		http.Error(w, "Can't delete file!", http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(w, "1")
}
