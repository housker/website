package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"./database"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

type webservice struct {
	client string
	host   string
	port   string
}

func setWebservice() *webservice {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	client := fmt.Sprintf("%s/client", dir)
	return &webservice{client: client, host: os.Getenv("CONN_HOST"), port: os.Getenv("CONN_PORT")}
}

func (ws *webservice) setRoutes() {
	mux := http.NewServeMux()
	ip := database.SetIntentProvider()

	mux.Handle("/", http.FileServer(http.Dir(ws.client)))
	mux.HandleFunc("/tags", ip.TagHandler)

	fmt.Println("Serving on", ws.host+":"+ws.port)
	err := http.ListenAndServe(ws.host+":"+ws.port,
		handlers.CompressHandler(mux))
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}

func importPython() {
	cmd := exec.Command("python", "-c", "import ai/prediction; print prediction.cat_strings('foo', 'bar')")
	fmt.Println(cmd.Args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	importPython()
	ws := setWebservice()
	ws.setRoutes()
}
