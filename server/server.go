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
	baseServer *httptest.Server
	messages   []string
}

func (s *Server) Stop() {
	s.baseServer.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := strings.Join(s.messages, "<br>")

	fmt.Fprint(w, msg)
}

func (s *Server) Printf(format string, a ...any) {
	s.messages = append(s.messages, fmt.Sprintf(format, a...))
}

func (s *Server) Println(msg string) {
	s.messages = append(s.messages, msg)
}

func (s *Server) PrintFile(fileName string) {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	s.messages = append(s.messages, string(fileBytes))
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
