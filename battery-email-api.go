
package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"log"

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
		w.Write([]byte("ERR"))
		panic(err)
	}
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	from := flag.String("from", "", "")
	smtpHost := flag.String("smtp-host", "", "")
	smtpPort := flag.String("smtp-port", "587", "optional")
	password := flag.String("password", "", "")

	flag.Parse()

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
		w.Write([]byte("battery-email-api by Ondrej Sika (sika.io)"))
	})

	http.HandleFunc("/notify-low", func(w http.ResponseWriter, r *http.Request) {
		device, _ := r.URL.Query()["device"]
		battery_level, _ := r.URL.Query()["battery_level"]
		to, _ := r.URL.Query()["to"]
		notify(w, *smtpHost, *smtpPort, *from, *password, to[0], "["+device[0]+"] Low Battery", "Battery level of "+device[0]+" is under "+battery_level[0]+"%. Please, charge it. O.")
	})

	http.HandleFunc("/notify-high", func(w http.ResponseWriter, r *http.Request) {
		device, _ := r.URL.Query()["device"]
		battery_level, _ := r.URL.Query()["battery_level"]
		to, _ := r.URL.Query()["to"]
		notify(w, *smtpHost, *smtpPort, *from, *password, to[0], "["+device[0]+"] High Battery", "Battery level of "+device[0]+" is over "+battery_level[0]+"%. Please, stop charging. O.")
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}