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
		fmt.Printf("Get: Can't convert id '%s' to int! %v\n", id_str, err)
		return
	}

	path, err := s.svc.Get(id)

	if err != nil {
		http.Error(w, fmt.Sprintf("Can't get file with id '%d'!", id), http.StatusNotFound)
		fmt.Printf("Get: Got error while uploading file! %v\n", err)
		return
	}

	// TODO: Change filename by changing r.URL
	http.ServeFile(w, r, path)
}

func upload(w http.ResponseWriter, r *http.Request) (string, error) {
	// Generate filename
	filename := randStringRunes(20)
	url := r.URL.Query().Get("url")

	var err error
	if url != "" {
		err = loadByUrl(filename, url)
	} else {
		err = loadByMultipart(filename, r)
	}

	if err != nil {
		return "", fmt.Errorf("upload() error: %w", err)
	}

	return filename, nil
}

func loadByUrl(filename, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("load(): empty file: %w", err)
	}
	defer resp.Body.Close()

	// Copy data
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("load(): can't create temp file: %w", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("load(): can't save data to file: %w", err)
	}

	return nil
}

func loadByMultipart(filename string, r *http.Request) error {
	// Modified of https://stackoverflow.com/questions/45541656/golang-send-file-via-post-request
	file, _, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("loadByMultipart(): can't read multipart: %w", err)
	}
	defer file.Close()

	// copy example
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)

	return nil
}

func (s *ApiServer) Set(w http.ResponseWriter, r *http.Request) {
	path, err := upload(w, r)

	if err != nil {
		http.Error(w, "Got error while uploading file!", http.StatusBadRequest)
		fmt.Printf("Set: Got error while uploading file! %v\n", err)
		return
	}

	id, err := s.svc.Set(path)

	if err != nil {
		http.Error(w, "Got error while saving file!", http.StatusBadRequest)
		fmt.Printf("Set: Got error while saving file! %v\n", err)
		return
	}

	fmt.Fprint(w, id)
}

func (s *ApiServer) Delete(w http.ResponseWriter, r *http.Request) {
	id_str := r.URL.Query().Get("id")
	id, err := strconv.Atoi(id_str)

	if err != nil {
		http.Error(w, "Can't convert id to int!", http.StatusNotFound)
		fmt.Printf("Delete: Can't convert id to int! %v\n", err)
		return
	}

	err = s.svc.Delete(id)

	if err != nil {
		http.Error(w, "Can't delete file!", http.StatusNotFound)
		fmt.Printf("Delete: Can't delete file! %v\n", err)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprint(w, "1")
}
