package main

import (
	"flag"
	"os"
	"net/http"
	"os/signal"
	"syscall"
	"log"
	"strings"
	"github.com/ssOleg/go_service/go_server/storage"
	"github.com/ssOleg/go_service/go_server/web"
	"fmt"
)

var port = flag.String("port", "", "port to run the server (Required)")

// main function to boot up everything
func main() {
	flag.Parse()
	if *port == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	dbStructure := storage.DataBase{ConnectionPoint: "localhost"}
	session, err := dbStructure.Connect()
	if err != nil {
		fmt.Println("Hello it is an error occured:", err)
		os.Exit(1)
	}

	defer session.Close()
	store := storage.NewStorage(session.DB("testDB"))

	//Remove old data from DataBase
	info, err := dbStructure.RemoveAll()
	log.Println("Dtatabase clear info: ", info)
	if err != nil {
		fmt.Println("Hello it is an error occured:", err)
		os.Exit(1)
	}
	//Insert data to database
	store.InsertInitialData()

	webRouter := web.Router{Storage: &dbStructure}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		<-sigs
		store.SaveData()
		os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(strings.Join([]string{"", *port}, ":"), web.GetRouter(webRouter)))

}
