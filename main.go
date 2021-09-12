package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func scanFiles() {

	files, err := ioutil.ReadDir("./tmp/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}
}

func main() {
	clearConsole()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("TEMP DIR:", path+"/tmp/")
	http.ListenAndServe(":9000", http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		clearConsole()
		scanFiles()
		if req.Method == "POST" {
			src, hdr, err := req.FormFile("my-file")
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			defer src.Close()

			dst, err := os.Create("./tmp/" + hdr.Filename)
			if err != nil {
				http.Error(res, err.Error(), 500)
				return
			}
			clearConsole()
			scanFiles()

			defer dst.Close()

			io.Copy(dst, src)
		}

		res.Header().Set("Content-Type", "text/html")
		io.WriteString(res, `
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.1/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-F3w7mX95PdgyTmZZMECAngseQB83DfGTowi0iMjiWaeVhAn4FJkqJByhZMI3AhiU" crossorigin="anonymous">
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.1/dist/js/bootstrap.bundle.min.js" integrity="sha384-/bQdsTh/da6pkI1MST/rWKFNjaCP5gBSY4sEBT38Q/9RBh9AH40zEOg7Hlq2THRZ" crossorigin="anonymous"></script>
		<div class="container">
		<div class="container">
		<div class="row">
		  <div class="col">
			1 of 3
		  </div>
		  <div class="col-6">
			2 of 3 (wider)
		  </div>
		  <div class="col">
			3 of 3
		  </div>
		</div>
		<div class="row">
		  <div class="col">
			1 of 3
		  </div>
		  <div class="col-5">
			2 of 3 (wider)
		  </div>
		  <div class="col">
			3 of 3
		  </div>
		</div>
	  </div>
		<form method="POST" enctype="multipart/form-data">
        <input type="file" name="my-file">
        <input type="submit">
        </form>
        `)
	}))
}
