package main

import (
	"fmt"
	"github.com/kubestaff/web-helper/server"
	"net/http"
	"time"
)

func main() {
	opts := server.Options{
		Port: 34567,
	}
	// we create the simplified web server
	s := server.NewServer(opts)

	// we close the server at the end
	defer s.Stop()

	s.Handle("/", HandleIndex)
	s.Handle("/status", HandleStatus)
	s.Handle("/months", HandleMonths)
	s.HandleJSON("/colors", HandleJsonOutput)
	s.HandleJSON("/add-color", HandleJsonInputFromParams)
	s.HandleJSON("/add-color-json", HandleJsonInputFromBody)

	s.Start()
}

type Color struct {
	Name string
	Code string
}

var colors = map[string]Color{
	"red": {
		Name: "red",
		Code: "#FF0000",
	},
	"blue": {
		Name: "blue",
		Code: "#0000FF",
	},
	"white": {
		Name: "white",
		Code: "#ffffff",
	},
	"black": {
		Name: "black",
		Code: "#000000",
	},
}

func HandleJsonInputFromBody(input server.Input) (o server.Output) {
	colorFromInput := Color{}

	err := input.Scan(&colorFromInput)

	if err != nil {
		return server.Output{
			Data: server.JsonError{
				Error: err.Error(),
				Code:  400,
			},
			Code: http.StatusBadRequest,
		}
	}

	if colorFromInput.Name == "" {
		return server.Output{
			Data: server.JsonError{
				Error: "Empty color name",
				Code:  400,
			},
			Code: http.StatusBadRequest,
		}
	}

	if colorFromInput.Code == "" {
		return server.Output{
			Data: server.JsonError{
				Error: "Empty color code",
				Code:  400,
			},
			Code: http.StatusBadRequest,
		}
	}

	colors[colorFromInput.Name] = colorFromInput

	return server.Output{
		Data: colorFromInput,
		Code: http.StatusOK,
	}
}

func HandleJsonInputFromParams(input server.Input) (o server.Output) {
	colorNameFromInput := input.Get("name")
	colorCodeFromInput := input.Get("code")

	if colorNameFromInput == "" {
		return server.Output{
			Data: server.JsonError{
				Error: "Empty color name",
				Code:  400,
			},
			Code: http.StatusBadRequest,
		}
	}

	if colorCodeFromInput == "" {
		return server.Output{
			Data: server.JsonError{
				Error: "Empty color code",
				Code:  400,
			},
			Code: http.StatusBadRequest,
		}
	}
	color := Color{
		Name: colorNameFromInput,
		Code: colorCodeFromInput,
	}
	colors[colorNameFromInput] = color

	return server.Output{
		Data: color,
		Code: http.StatusOK,
	}
}

func HandleJsonOutput(input server.Input) (o server.Output) {
	colorNameFromInput := input.Get("color")
	if colorNameFromInput == "" {
		return server.Output{
			Data: colors,
			Code: http.StatusOK,
		}
	}

	color, ok := colors[colorNameFromInput]
	if !ok {
		return server.Output{
			Data: server.JsonError{
				Error: fmt.Sprintf("%s color not found", colorNameFromInput),
				Code:  404,
			},
			Code: http.StatusNotFound,
		}
	}

	return server.Output{
		Data: color,
		Code: http.StatusOK,
	}
}

func HandleStatus(inputs server.Input) (filename string, placeholders map[string]string) {
	return "html/status.html", nil
}

func HandleIndex(inputs server.Input) (filename string, placeholders map[string]string) {
	variables := map[string]string{"%name%": "Max Mustermann"}
	return "html/index.html", variables
}

func HandleMonths(input server.Input) (filename string, placeholders map[string]string) {
	months := map[string]string{
		"1":  "Jan",
		"2":  "Feb",
		"3":  "Mar",
		"4":  "Apr",
		"5":  "Mai",
		"6":  "Jun",
		"7":  "Jul",
		"8":  "Aug",
		"9":  "Sep",
		"10": "Oct",
		"11": "Nov",
		"12": "Dec",
	}

	givenMonthNumber := input.Get("month")
	if givenMonthNumber == "" {
		givenMonthNumber = time.Now().Format("1")
	}

	output := map[string]string{
		"%month%": fmt.Sprintf("Number %s is month %s", givenMonthNumber, months[givenMonthNumber]),
	}

	return "html/month.html", output
}
