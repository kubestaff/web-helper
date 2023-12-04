package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func (s *Server) inputsFromRequest(request *http.Request) url.Values {
	inputs := request.URL.Query()

	postParams := request.PostForm
	for key, values := range postParams {
		inputs[key] = values
	}

	return inputs
}

type Output struct {
	Data any
	Code int
}

type JsonError struct {
	Code  int
	Error string
}

func (s *Server) jsonErr(writer http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = 200
	}

	jsonErr := JsonError{
		Code:  code,
		Error: err.Error(),
	}

	writer.WriteHeader(jsonErr.Code)

	jsonData, err := json.Marshal(jsonErr)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Input struct {
	url.Values
	body io.Reader
}

func (i Input) Scan(target any) error {
	err := json.NewDecoder(i.body).Decode(target)
	if err == io.EOF {
		return errors.New("empty request data provided")
	}

	return err
}

func (s *Server) HandleJSON(url string, handler func(input Input) (o Output)) {
	s.router.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusOK)
			return
		}

		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)
		if err != nil {
			s.jsonErr(writer, fmt.Errorf("can't read request body: %v", err), 0)
			return
		}

		buffer := bytes.NewBuffer(body)

		requestValues := s.inputsFromRequest(request)
		inpt := Input{
			Values: requestValues,
			body:   buffer,
		}
		output := handler(inpt)

		if output.Code == 0 {
			output.Code = 200
		}

		jsonData, err := json.Marshal(output.Data)
		if err != nil {
			s.jsonErr(writer, err, 0)
			return
		}

		writer.WriteHeader(output.Code)

		_, err = writer.Write(jsonData)
		if err != nil {
			s.jsonErr(writer, err, 0)
			return
		}
	})
}

func (s *Server) Handle(url string, handler func(input Input) (filename string, placeholders map[string]string)) {
	s.router.HandleFunc(url, func(writer http.ResponseWriter, request *http.Request) {
		values := s.inputsFromRequest(request)

		inpt := Input{
			Values: values,
			body:   request.Body,
		}

		fileName, variables := handler(inpt)
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
