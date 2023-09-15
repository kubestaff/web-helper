package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

type Server struct {
	URL        string
	body       string
	baseServer *httptest.Server
}

func (s *Server) Stop() {
	s.baseServer.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, s.body)
}

func (s *Server) PrintFile(fileName string, variables map[string]string) {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fileContent := string(fileBytes)
	for key, value := range variables {
		fileContent = strings.ReplaceAll(fileContent, key, value)
	}

	s.body = fileContent
}

func (s *Server) Start() {
	s.baseServer.Start()
	fmt.Println("started at " + s.baseServer.URL)
	select {}
}

func NewServer() *Server {
	mainServer := &Server{}

	router := http.NewServeMux()
	router.Handle("/", mainServer)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	baseServer := httptest.NewUnstartedServer(router)
	mainServer.baseServer = baseServer
	mainServer.URL = baseServer.URL

	return mainServer
}
