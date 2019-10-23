package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	PageTitle   string
	Success     bool
	Message     string
	MessageType string
}

type Food struct {
	Name string
}

func main() {
	data := []Food{
		{
			Name: "burger",
		},
		{
			Name: "pizza",
		},
	}
	funcMap := template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	tmpl := template.Must(template.New("layout.html").Funcs(funcMap).ParseFiles("layout.html"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		// if req.URL.Path != "/" {
		// 	http.NotFound(w, req)
		// 	return
		// }
		data := PageData{
			PageTitle: "Welcome to my web server",
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			log.Print(err)
		}
	})
	mux.HandleFunc("/foods", func(w http.ResponseWriter, req *http.Request) {
		resp, err := json.Marshal(data)
		if err != nil {
			log.Print(err)
		}
		fmt.Fprint(w, string(resp))
	})
	mux.HandleFunc("/create-food", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.Method != http.MethodPost {
			http.NotFound(w, req)
			return
		}

		name := req.FormValue("foodName")

		data = append(data, Food{Name: name})

		err := tmpl.Execute(w, PageData{
			Success:     true,
			Message:     fmt.Sprintf("Food added: %s", name),
			MessageType: "success",
		})
		if err != nil {
			log.Print(err)
		}
	})

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	log.Println("Starting server at http://127.0.0.1")
	log.Println(srv.ListenAndServe())
}
