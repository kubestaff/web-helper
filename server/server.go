package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const DefaultPort = 2345

type Server struct {
	body       string
	baseServer *http.Server
	opts       Options
	router     *http.ServeMux
}

func (s *Server) Stop() {
	if s.baseServer != nil {
		s.baseServer.Close()
	}
}

func (s *Server) Handle(url string, handler func(inputs map[string]string) (filename string, placeholders map[string]string)) {
	s.router.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
		inputs := make(map[string]string)

		queryParams := request.URL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				inputs[key] = value
			}
		}

		postParams := request.PostForm

		for key, values := range postParams {
			for _, value := range values {
				inputs[key] = value
			}
		}

		fileName, variables := handler(inputs)
		fileBytes, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		fileContent := string(fileBytes)
		for key, value := range variables {
			fileContent = strings.ReplaceAll(fileContent, key, value)
		}

		fmt.Fprintf(writer, fileContent)
	})
}

func (s *Server) PrintFile(url, fileName string, variables map[string]string) {
	s.router.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
		fileBytes, err := os.ReadFile(fileName)
		if err != nil {
			http.Error(writer, err.Error(), 500)
			return
		}

		fileContent := string(fileBytes)
		for key, value := range variables {
			fileContent = strings.ReplaceAll(fileContent, key, value)
		}

		queryParams := request.URL.Query()
		for key, values := range queryParams {
			for _, value := range values {
				fileContent = strings.ReplaceAll(fileContent, "%"+key+"%", value)
			}
		}

		postParams := request.PostForm

		for key, values := range postParams {
			for _, value := range values {
				fileContent = strings.ReplaceAll(fileContent, "%"+key+"%", value)
			}
		}

		fmt.Fprintf(writer, fileContent)
	})
}

func (s *Server) Start() {
	svr := s.initBaseServer()
	s.baseServer = svr

	fmt.Printf("started at http://127.0.0.1%s\n", s.baseServer.Addr)
	err := s.baseServer.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) GetUrl() string {
	if s.baseServer != nil {
		return s.baseServer.Addr
	}

	return ""
}

func (s *Server) initBaseServer() *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.opts.Port),
		Handler: s.router,
	}

	return server
}

func NewServer(opts Options) *Server {
	if opts.Port == 0 {
		opts.Port = DefaultPort
	}

	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mainServer := &Server{
		opts:   opts,
		router: router,
	}

	return mainServer
}
