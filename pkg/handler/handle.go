package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sheikh-arman/api-server-containerized/pkg/structure"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Appscode struct {
	DBHost     string `json:"dbhost"`
	DBName     string `json:"dbname"`
	DBPassword string `json:"dbpassword"`
}

var (
	ID = 0
)

func Handle() {
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
		DBHost:     dbHost,
		DBName:     dbName,
		DBPassword: dbPass,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		// jwtauth-> will learn later
		r.Route("/", func(r chi.Router) {
			r.Get("/", in.getFunc)
			r.Get("/employee", in.getInfo)
			r.Post("/employee", in.addInfo)
			r.Put("/employee", in.editInfo)
			r.Delete("/employee", in.deleteInfo)
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
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "Welcome")
	fmt.Fprintln(w, "Calling time:"+strconv.Itoa(ID))
	w.WriteHeader(403)
}

func (in Appscode) getInfo(w http.ResponseWriter, r *http.Request) {
	db := in.openConnection()
	rows, err := db.Query("select * from appscode.employee")
	if err != nil {
		log.Fatalf("querying the books table %s\n", err.Error())
	}
	employee := []structure.AppscodeEmployee{}
	for rows.Next() {
		var id int
		var name, salary string
		err := rows.Scan(&id, &name, &salary)
		if err != nil {
			log.Fatalf("while scanning the row %s\n", err.Error())
		}
		log.Println(id, name, salary)
		obj := structure.AppscodeEmployee{
			Id: id, Name: name, Salary: salary,
		}
		employee = append(employee, obj)
	}
	in.closeConnection(db)
	writeJsonResponse(w, 200, employee)
}

func (in Appscode) addInfo(w http.ResponseWriter, r *http.Request) {
	data := structure.AppscodeEmployee{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatalf("issue: %s\n", err.Error())
	}
	log.Println(data)
	db := in.openConnection()
	insertQuery, err := db.Prepare("insert into appscode.employee ( name, salary) values ( ?, ?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(data.Name, data.Salary)
	if err != nil {
		log.Fatalf("execing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commint the transaction %s\n", err.Error())
	}
	in.closeConnection(db)
}
func (in Appscode) editInfo(w http.ResponseWriter, r *http.Request) {
	data := structure.AppscodeEmployee{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatalf("issue: %s\n", err.Error())
	}
	log.Println(data)
	db := in.openConnection()
	insertQuery, err := db.Prepare("UPDATE appscode.employee SET name = ?, salary = ?  WHERE id = ?")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(data.Name, data.Salary, data.Id)
	if err != nil {
		log.Fatalf("execing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commint the transaction %s\n", err.Error())
	}
	in.closeConnection(db)
}
func (in Appscode) deleteInfo(w http.ResponseWriter, r *http.Request) {
	data := 0
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Fatalf("issue: %s\n", err.Error())
	}
	log.Println(data)
	db := in.openConnection()
	insertQuery, err := db.Prepare("DELETE FROM appscode.employee WHERE id = ?")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(data)
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
	err := json.NewEncoder(w).Encode(data)
	log.Fatalf("issue writing json response %s\n", err)
}

func (in Appscode) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", in.DBPassword, in.DBHost, in.DBName))
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
