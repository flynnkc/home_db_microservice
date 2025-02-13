package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/flynnkc/home_db_microservice/pkg/handlers/models"
)

const (
	insertTemp string = "INSERT INTO temperature(room, temperature, humidity) VALUES(?, ?, ?)"
)

var (
	log            *slog.Logger = slog.Default()
	db             models.Database
	tempInsertStmt *sql.Stmt
)

func LoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				log.Error("error in parse form",
					"error", err)
			}
			log.Info("incoming request",
				"request", fmt.Sprintf("%+v", r))
			h.ServeHTTP(w, r)
		})
}

func TestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Test")
			h.ServeHTTP(w, r)
		})
}

func RecoveryMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					log.Error("recovering",
						"error", r)
				}
			}()

			h.ServeHTTP(w, r)
		})
}

// Room varchar255 temp float humidity float
func MysqlTempHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure database and prepared statements are initialized
	if db == nil {
		d, err := models.NewMysqlDB(os.Getenv("DB_CONN"))
		if err != nil {
			panic(err)
		}

		db = d
		tempInsertStmt, err = db.PrepareStmt(insertTemp)
		if err != nil {
			panic(err)
		}
	}

	switch r.Method {
	case "GET":
		w.Write([]byte("success\n"))
	case "POST":
		room := r.PostFormValue("room")
		temp := r.PostFormValue("temp")
		hum := r.PostFormValue("humidity")
		if room == "" || temp == "" || hum == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request -- missing value"))
			return
		}

		result, err := tempInsertStmt.Exec(room, temp, hum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Error(err.Error())
		}

		rows, err := result.RowsAffected()
		if err != nil {
			log.Error(err.Error())
		}

		w.Write([]byte(fmt.Sprintf("id %v - rows affected %v\n", id, rows)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}
}
