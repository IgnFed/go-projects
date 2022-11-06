package src

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type server interface {
	Start() bool
}

type s struct{}

var Server server = &s{}

type album struct {
	ID     int64
	Title  string
	Artist string
	Price  float64
}

var db *sql.DB

func getAlbums(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("SELECT * FROM ALBUM")
	var albums []album
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var album album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			log.Fatal(err)
		}
		albums = append(albums, album)
	}
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "data": albums})
}

func getAlbumByID(w http.ResponseWriter, r *http.Request) {
	var album album
	vars := mux.Vars(r)
	var id = vars["id"]
	row := db.QueryRow("SELECT * FROM ALBUM WHERE ID=?", id)
	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "data": album})
}

func postAlbum(r http.ResponseWriter, w *http.Request) {
	var new_album album
	err := json.NewDecoder(w.Body).Decode(&new_album)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec("INSERT INTO ALBUM (title, artist, price) VALUES (?,?,?)", new_album.Title, new_album.Artist, new_album.Price)

	if err != nil {
		log.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(r).Encode(map[string]any{"ok": true, "insert_id": id, "data": new_album})
}

func (s *s) Start() bool {
	Configs := mysql.Config{
		User:   os.Getenv("USERDB"),
		Passwd: os.Getenv("PASSDB"),
		Net:    "tcp",
		Addr:   "127.0.0.1",
		DBName: os.Getenv("DBNAME"),
	}

	r := mux.NewRouter()
	var err error
	db, err = sql.Open(os.Getenv("DBDRIVER"), Configs.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/albums", getAlbums).Methods("GET")
	api.HandleFunc("/albums/{id}", getAlbumByID).Methods("GET")
	api.HandleFunc("/albums", postAlbum).Methods("POST")
  http.ListenAndServe(":80", r)

	return true
}
