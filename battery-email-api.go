package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ondrejsika/battery-email-api/battery_state"
	"github.com/ondrejsika/gosendmail/lib"
)

func notify(
	w http.ResponseWriter,
	smtpHost string,
	smtpPort string,
	from string,
	password string,
	to string,
	subject string,
	message string) {
	err := lib.GoSendMail(
		smtpHost,
		smtpPort,
		from,
		password,
		to,
		subject,
		message)
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
	dbDriver := flag.String("db-driver", "none", "")
	dbConnection := flag.String("db-connection", "", "")

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

	bs := battery_state.Init(*dbDriver, *dbConnection)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(
			"[battery-email-api] " +
				"battery-email-api v3 by Ondrej Sika (sika.io)"))
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
		if bs != nil {
			bs.Add(device[0], level[0])
		}
		if level[0] == "low" {
			notify(
				w,
				*smtpHost,
				*smtpPort,
				*from,
				*password,
				to[0],
				"["+device[0]+"] Low Battery",
				"Battery level of "+device[0]+" is under "+
					battery_level[0]+"%. Please, charge it. O. \n\n--\n"+r.Host)
			return
		}
		if level[0] == "high" {
			notify(
				w,
				*smtpHost,
				*smtpPort,
				*from,
				*password,
				to[0],
				"["+device[0]+"] High Battery",
				"Battery level of "+device[0]+" is over "+
					battery_level[0]+"%. Please, stop charging. O. \n\n--\n"+r.Host)
			return
		}
		w.WriteHeader(400)
		w.Write([]byte("[battery-email-api] No level"))
	})

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {
		tokenFromRequest := r.URL.Query()["token"]
		if *token != tokenFromRequest[0] {
			w.WriteHeader(403)
			w.Write([]byte("[battery-email-api] Wrong token"))
			return
		}
		if *dbDriver == "none" {
			w.WriteHeader(400)
			w.Write([]byte("[battery-email-api] Battery state is not enabled. " +
				"Driver is none."))
			return
		}
		device := r.URL.Query()["device"]
		b := bs.Get(device[0])
		var jsonOutput []byte
		if b.ID == 0 {
			jsonOutput, _ = json.Marshal(nil)
		} else {
			jsonOutput, _ = json.Marshal(b.Level)
		}
		w.WriteHeader(200)
		w.Write([]byte(jsonOutput))
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
