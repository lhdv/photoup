package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func main() {
	host := ""
	port := "3000"

	log.Println(host, port)

	http.HandleFunc("/", rootPage)
	http.HandleFunc("/photoup", photoUp)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	var data struct {}

	t := template.Must(template.ParseFiles("index.html"))
	err := t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func photoUp(w http.ResponseWriter, r *http.Request) {

	var totalWritten int64

	err := r.ParseMultipartForm(256000)
	if (err != nil) {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
	}

	origUser := r.PostForm["origUser"][0]

	log.Println(origUser)

	for _, files := range r.MultipartForm.File {
		log.Printf("[files] - %+v\n\n", files)
		for _, file := range files {

			// open uploaded  
			var infile multipart.File
	                if infile, err = file.Open(); nil != err {
				http.Error(w, "Error opening input file", http.StatusInternalServerError)
				return
			}

			// open destination  
			var outfile *os.File
			if outfile, err = os.Create("./uploaded/" + file.Filename); nil != err {
				http.Error(w, "Error opening output file", http.StatusInternalServerError)
				return
			}

			// 32K buffer copy  
			var written int64
			if written, err = io.Copy(outfile, infile); nil != err {
				http.Error(w, "Error writing file", http.StatusInternalServerError)
				return
			}

			if written > 1024 {
				fmt.Fprintf(w, "uploaded file: %s : length: %s kbytes\n",file.Filename, strconv.Itoa(int(written/1024)))
			} else {
				fmt.Fprintf(w, "uploaded file: %s : length: %s bytes\n",file.Filename, strconv.Itoa(int(written)))
			}

			totalWritten += written
		}
	}

	if totalWritten > 0 {

		fmt.Fprintf(w, "Total %s kbytes\n", strconv.Itoa(int(totalWritten/1024)))
	}
}

