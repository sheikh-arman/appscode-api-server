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
	ID = 0
)

type AppscodeEmployee struct {
	id     string `json:"id"`
	name   string `json:"name"`
	salary string `json:"salary"`
}

type Appscode struct {
	dbHost     string `json:"dbhost"`
	dbName     string `json:"dbname"`
	dbPassword string `json:"dbpassword"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
	Id    int    `json:"id"`
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
	//r.Get("/", in.getFunc)
	//r.Get("/info", in.getInfo)
	//r.Post("/add", in.addInfo)
	r.Group(func(r chi.Router) {
		// jwtauth-> will learn later
		r.Route("/", func(r chi.Router) {
			r.Get("/", in.getFunc)
			r.Get("/info", in.getInfo)
			r.Post("/add", in.addInfo)
		})
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	fmt.Println("server started on port 8080")
	fmt.Println(server.ListenAndServe())

}

func (in Appscode) getFunc(w http.ResponseWriter, r *http.Request) {
	ID += 1
	data := "Welcome\nCalling time:" + strconv.Itoa(ID)
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
		var id string
		var name, salary string
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
	//var feeds []Item
	//var feed Item
	//feed = Item{
	//	Id:    ID,
	//	Title: "Nothing",
	//	Post:  "Lorem Ipsum Doller Site",
	//}
	////feeds2[ID] = feed
	//ID++
	//feeds = append(feeds, feed)
	//
	//feed = Item{
	//	Id:    ID,
	//	Title: "Nothing2",
	//	Post:  "Lorem Ipsum Doller Site2",
	//}
	////feeds2[ID] = feed
	//ID++
	//feeds = append(feeds, feed)
	//
	//feed = Item{
	//	Id:    ID,
	//	Title: "Nothing3",
	//	Post:  "Lorem Ipsum Doller Site3",
	//}
	////feeds2[ID] = feed
	//ID++
	//feeds = append(feeds, feed)
	//writeJsonResponse(w, 200, feeds)
}

func (in Appscode) addInfo(w http.ResponseWriter, r *http.Request) {
	data := AppscodeEmployee{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatalf("issue: ", err.Error())
	}
	log.Println(data)
	db := in.openConnection()
	insertQuery, err := db.Prepare("insert into appscode values (?, ?, ?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(data.id, data.name, data.salary)
	if err != nil {
		log.Fatalf("execing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commint the transaction %s\n", err.Error())
	}
	in.closeConnection(db)
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

//{
//"id":"3",
//"name":"asdsa",
//"salary":"ada"
//}
