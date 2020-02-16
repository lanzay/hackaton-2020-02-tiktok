package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	API_ADDR     = ":8080"
	DIR_DATA_IN  = "data_in"
	DIR_DATA_OUT = "data_out"
)

func WebRun() {

	os.MkdirAll(DIR_DATA_IN, 0777)
	os.MkdirAll(DIR_DATA_OUT, 0777)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./tiktok/build"))
	mux.Handle("/", fs)

	mux.HandleFunc("/ping", wwwPing)
	mux.HandleFunc("/f1/", wwwf1)
	mux.HandleFunc("/v/", wwwV)

	log.Println("[I] Server run on", API_ADDR)
	log.Println(http.ListenAndServe(API_ADDR, mux))

}

func wwwf1(w http.ResponseWriter, r *http.Request) {

	log.Println("[D] Reqest on [f1] fronm", r.RemoteAddr, r.Method, r.URL.Path)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, headers := range r.MultipartForm.File {
		for fileN, header := range headers {
			log.Printf("[D] Receive file #%d from user %s - %dKb\n", fileN, header.Filename, int(header.Size)/1024)
			fd, err := header.Open()
			if err != nil {
				log.Println("[E]", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fBody, err := ioutil.ReadAll(fd)
			fNew, err := os.Create(DIR_DATA_IN + "/f1_" + time.Now().Format("02-01-2006_15-04-05") + "_" + strconv.Itoa(time.Now().Nanosecond()) + ".jpg")
			if err != nil {
				log.Println("[E]", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			_, err = fNew.Write(fBody)
			if err != nil {
				log.Println("[E]", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fd.Close()
			fNew.Close()
			//log.Println(fd, err)
		}
	}

	//v := "/v/Dance_tutorial_on_Tik_Tok.wmv"
	v := "/v/2.mp4"
	w.Write([]byte(v))
}

func wwwV(w http.ResponseWriter, r *http.Request) {

	tmp := strings.Split(r.URL.Path, "/v/")
	f, err := os.Open(DIR_DATA_OUT + "/" + tmp[len(tmp)-1])
	if err != nil {
		log.Println("[E]", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	body, _ := ioutil.ReadAll(f)
	w.Write(body)
	w.Header().Add("Access-Control-Allow-Origin", "*;http://localhost:3000")
}

func wwwPing(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*;http://localhost:3000")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Write([]byte("Pong " + r.RemoteAddr))
}

func Cmd(arg ...string) {

	// docker run --gpus all -it -v /host/data:/data -v /host/config:/config ufoym/deepo bash

	cmd := exec.Command(arg[0], arg...)
	cmd.Env = append(os.Environ(),
		"FOO=duplicate_value", // ignored
		"FOO=actual_value",    // this value is used
	)
	var out bytes.Buffer
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("in all caps: %q\n", out.String())

}
