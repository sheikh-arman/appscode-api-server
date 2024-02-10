package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	id = 0
)

type AppscodeEmployee struct {
	id, name, salary string
}

type Appscode struct {
	dbHost     string
	dbName     string
	dbPassword string
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "appscode"
	}
	dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
	if dbPass == "" {
		dbPass = "arman"
	}
	in := Appscode{
		dbHost:     dbHost,
		dbName:     dbName,
		dbPassword: dbPass,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", in.getFunc)
	r.Get("/info", in.getInfo)
	r.Post("/add", in.addInfo)
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println(server.ListenAndServe())
}

func (in Appscode) getFunc(w http.ResponseWriter, r *http.Request) {
	id += 1
	data := "Welcome\nCalling time:" + strconv.Itoa(id)
	writeJsonResponse(w, 200, data)
}

func (in Appscode) getInfo(w http.ResponseWriter, r *http.Request) {
	db := in.openConnection()
	rows, err := db.Query("select * from appscode")
	if err != nil {
		log.Fatalf("querying the books table %s\n", err.Error())
	}
	employee := []AppscodeEmployee{}
	for rows.Next() {
		var id, name, salary string
		err := rows.Scan(&id, &name, &salary)
		if err != nil {
			log.Fatalf("while scanning the row %s\n", err.Error())
		}
		log.Println(id, name, salary)
		employee = append(employee, AppscodeEmployee{id: id, name: name, salary: salary})
	}
	log.Println(employee)
	err = json.NewEncoder(w).Encode(employee)
	if err != nil {
		log.Fatalf("encoding employees: %s\n", err.Error())
	}
	//writeJsonResponse(w, 200, employee)
	in.closeConnection(db)
}

func (in Appscode) addInfo(w http.ResponseWriter, r *http.Request) {
	id += 1
	data := "Welcome\nCalling time:" + strconv.Itoa(id)

	writeJsonResponse(w, 200, data)
}

func writeJsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (in Appscode) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", in.dbPassword, in.dbHost, in.dbName))
	if err != nil {
		log.Fatalf("opening the connection to the database %s\n", err.Error())
	}
	return db
}

func (in Appscode) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
