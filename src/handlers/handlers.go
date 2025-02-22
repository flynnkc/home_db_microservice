package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/flynnkc/home_db_microservice/src/pkg/handlers/models"
)

const (
	insertTemp string = "INSERT INTO temperature(room, temperature, humidity) VALUES(?, ?, ?)"
)

var (
	log            *slog.Logger = slog.Default()
	homeDb         models.Database
	tempInsertStmt *sql.Stmt
)

// Room varchar255 temp float humidity float
func MysqlTempHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure database and prepared statements are initialized
	if homeDb == nil {
		db, err := models.NewMysqlDB(os.Getenv("DB_CONN"))
		if err != nil {
			panic(err)
		}

		homeDb = db
		tempInsertStmt, err = homeDb.PrepareStmt(insertTemp)
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
