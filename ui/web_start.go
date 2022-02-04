package ui

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/gorilla/mux"
)

//Will Bind to this Intereface
//If Empty binds to all interfaces

//InterfaceIp : Bind to this Interface
var InterfaceIp binding.String = binding.NewString()

// Port ...
var Port binding.Int = binding.NewInt()

var r *mux.Router = mux.NewRouter()

// ServerAlive ...
var ServerAlive binding.Bool = binding.NewBool()

var srv *http.Server

// AddRouting ...
func AddRouting() {
	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "Hello %v\n Server is Alive\n", r.RemoteAddr)
	})

	r.HandleFunc("/page/{pname}", GenericPage)
	r.HandleFunc("/web/working", WorkingSubdomain)
	r.HandleFunc("/org/checklist", OrgCheckListRoute)
	r.HandleFunc("/web/checklist", WebCheckList)
	r.HandleFunc("/org/{orgitem}", OrgToolOutput)
	r.HandleFunc("/web/{sub}/{webitem}", WebToolOutput)
	r.HandleFunc("/commit", Commit)
}

// StartServer ...
func StartServer() {
	alive, _ := ServerAlive.Get()
	if !alive {
		AddRouting()
		ip, _ := InterfaceIp.Get()
		port, _ := Port.Get()
		Address := ip + ":" + strconv.Itoa(port)

		srv := &http.Server{
			Addr:         Address,
			WriteTimeout: time.Duration(15) * time.Second,
			ReadTimeout:  time.Duration(15) * time.Second,
			IdleTimeout:  time.Duration(60) * time.Second,
		}

		srv.Handler = r

		log.Printf("Listening at %v\n", Address)

		log.Fatal(srv.ListenAndServe())

		ServerAlive.Set(true)
	}
}

// StopServer ...
func StopServer() {
	alive, _ := ServerAlive.Get()
	if alive && srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
		log.Println("shutting down")
		ServerAlive.Set(false)
	}
}
