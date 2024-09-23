package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type PageData struct {
	Title  string
	Result string
}

func main() {
	http.HandleFunc("/", home)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		filepath.Join("ui", "html", "base.tmpl"),
		filepath.Join("ui", "html", "partials", "nav.tmpl"),
		filepath.Join("ui", "html", "pages", "home.tmpl"),
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	if r.Method == http.MethodPost {
		valueStr := r.FormValue("value")
		fromUnit := r.FormValue("from_unit")
		toUnit := r.FormValue("to_unit")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		result, err := convert(value, fromUnit, toUnit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := PageData{
			Title:  "Unit Converter",
			Result: fmt.Sprintf("%.2f %s = %.2f %s", value, fromUnit, result, toUnit),
		}

		ts.Execute(w, data)
		return
	}

	data := PageData{
		Title: "Unit Converter",
	}

	ts.Execute(w, data)
}

func convert(value float64, fromUnit string, toUnit string) (float64, error) {
	switch fromUnit {
	case "meters":
		switch toUnit {
		case "kilometers":
			return value / 1000, nil
		case "meters":
			return value, nil
		default:
			return 0, fmt.Errorf("unsupported conversion: %s to %s", fromUnit, toUnit)
		}

	case "kilometers":
		switch toUnit {
		case "meters":
			return value * 1000, nil
		case "kilometers":
			return value, nil
		default:
			return 0, fmt.Errorf("unsupported conversion: %s to %s", fromUnit, toUnit)
		}

	case "grams":
		switch toUnit {
		case "kilograms":
			return value / 1000, nil
		case "grams":
			return value, nil
		default:
			return 0, fmt.Errorf("unsupported conversion: %s to %s", fromUnit, toUnit)
		}

	case "kilograms":
		switch toUnit {
		case "grams":
			return value * 1000, nil
		case "kilograms":
			return value, nil
		default:
			return 0, fmt.Errorf("unsupported conversion: %s to %s", fromUnit, toUnit)
		}

	default:
		return 0, fmt.Errorf("unsupported conversion: %s", fromUnit)
	}
}
