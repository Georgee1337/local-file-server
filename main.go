package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	listen = flag.String("listen", ":3000", "listen address")
)

func filesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(200000)
	if err != nil {
		log.Printf("Err: %v", err)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Err: %v", err)
		return
	}
	defer file.Close()

	resFile, err := os.Create("./data/" + handler.Filename)
	if err != nil {
		log.Printf("Err: %v", err)
	}
	defer resFile.Close()
	if err == nil {
		io.Copy(resFile, file)
		defer resFile.Close()
		log.Printf("Successfully uploaded file %q...", handler.Filename)
		http.Redirect(w, r, "/files/success", http.StatusSeeOther)
	}

}

func main() {
	flag.Parse()
	log.Printf("Starting server at port %q...", *listen)
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/files/upload", filesHandler)

	http.ListenAndServe(*listen, nil)

}
