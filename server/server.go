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
	body           string
	baseTestServer *httptest.Server
	baseServer     *http.Server
	opts           Options
}

func (s *Server) Stop() {
	if s.baseServer != nil {
		s.baseServer.Close()
		return
	}

	s.baseTestServer.Close()
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
	if s.baseServer != nil {
		fmt.Printf("started at http://127.0.0.1%s\n", s.baseServer.Addr)
		err := s.baseServer.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}

	s.baseTestServer.Start()
	fmt.Printf("started at %s\n", s.baseTestServer.URL)
	select {}
}

func (s *Server) GetUrl() string {
	if s.baseServer != nil {
		return s.baseServer.Addr
	}

	return s.baseTestServer.URL
}

func NewServer(opts Options) *Server {
	mainServer := &Server{
		opts: opts,
	}

	router := http.NewServeMux()
	router.Handle("/", mainServer)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	if opts.Port > 0 {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}
		mainServer.baseServer = server
	} else {
		mainServer.baseTestServer = httptest.NewUnstartedServer(router)
	}

	return mainServer
}
