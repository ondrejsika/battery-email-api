package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ondrejsika/gosendmail/lib"
)

func notify(w http.ResponseWriter,
	smtpHost string,
	smtpPort string,
	from string,
	password string,
	to string,
	subject string,
	message string) {
	err := lib.GoSendMail(smtpHost, smtpPort, from, password, to, subject, message)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("[battery-email-api] ERR"))
		panic(err)
	}
	w.WriteHeader(200)
	w.Write([]byte("[battery-email-api] OK"))
}

func main() {
	token := flag.String("token", "", "")
	from := flag.String("from", "", "")
	smtpHost := flag.String("smtp-host", "", "")
	smtpPort := flag.String("smtp-port", "587", "optional")
	password := flag.String("password", "", "")

	flag.Parse()

	if *token == "" {
		fmt.Fprintf(os.Stderr, "-token is not defined\n")
		os.Exit(1)
	}
	if *from == "" {
		fmt.Fprintf(os.Stderr, "-from is not defined\n")
		os.Exit(1)
	}
	if *smtpHost == "" {
		fmt.Fprintf(os.Stderr, "-smtp-host is not defined\n")
		os.Exit(1)
	}
	if *smtpPort == "" {
		fmt.Fprintf(os.Stderr, "-smtp-port is not defined\n")
		os.Exit(1)
	}
	if *password == "" {
		fmt.Fprintf(os.Stderr, "-password is not defined\n")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("[battery-email-api] battery-email-api v2 by Ondrej Sika (sika.io)"))
	})

	http.HandleFunc("/api/notify", func(w http.ResponseWriter, r *http.Request) {
		tokenFromRequest := r.URL.Query()["token"]
		if *token != tokenFromRequest[0] {
			w.WriteHeader(403)
			w.Write([]byte("[battery-email-api] Wrong token"))
			return
		}
		device := r.URL.Query()["device"]
		battery_level := r.URL.Query()["battery_level"]
		to := r.URL.Query()["to"]
		level := r.URL.Query()["level"]
		if level[0] == "low" {
			notify(w, *smtpHost, *smtpPort, *from, *password, to[0], "["+device[0]+"] Low Battery", "Battery level of "+device[0]+" is under "+battery_level[0]+"%. Please, charge it. O.")
			return
		}
		if level[0] == "high" {
			notify(w, *smtpHost, *smtpPort, *from, *password, to[0], "["+device[0]+"] High Battery", "Battery level of "+device[0]+" is over "+battery_level[0]+"%. Please, stop charging. O.")
			return
		}
		w.WriteHeader(400)
		w.Write([]byte("[battery-email-api] No level"))
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
