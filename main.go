package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"signer/config"
	"signer/repo"
	serv "signer/server"
	"signer/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := mux.NewRouter()

	connStr := fmt.Sprintf(cfg.Database.ConnectionString, cfg.Database.Username, cfg.Database.Password, cfg.Database.DatabaseIP)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Print(err)
		return
	}

	tr := repo.NewTestRepo(db)
	ur := repo.NewUserRepo(db)
	qar := repo.NewQARepo(db)

	ts := service.NewTestService(cfg, &tr, &qar)
	us := service.NewUserService(cfg, &ur)

	handler := serv.NewHandler(us, ts)

	router.HandleFunc("/login", handler.GetToken).Methods("POST")
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/sign", handler.Sign).Methods("POST")
	router.HandleFunc("/check-sig", handler.CheckSignature).Methods("POST")

	fmt.Println("Starting the server")
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port), router)
	if err != nil {
		fmt.Println("Could not start the server", err)
	}
	fmt.Println("Server started. Listenning on port 4000")
}
