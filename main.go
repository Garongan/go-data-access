package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Album struct {
	Id     int64   `db:"id"`
	Title  string  `db:"title"`
	Artist string  `db:"artist"`
	Price  float32 `db:"price"`
}

const (
	host   = "localhost"
	port   = 5432
	dbName = "recording"
)

var db *sql.DB

func main() {

	// information to connect psql
	psqlInfo := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", host, port, dbName, os.Getenv("user"), os.Getenv("pass"))

	// open connection
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	checkError(err)
	defer db.Close()

	// get realtime error from database
	pingErr := db.Ping()
	checkError(pingErr)

	fmt.Println("Connected!")

	// get data by artist
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// get data by id
	album, err := albumsById(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// query recording with specific artist name
func albumsByArtist(name string) ([]Album, error) {
	// slice albums
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %s %v", name, err)
	}
	defer rows.Close()

	// scan each row
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %s %v", name, err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %s %v", name, err)
	}
	return albums, nil
}

func albumsById(id int64) (Album, error) {
	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsById %d, no such album", id)
		}
		return album, fmt.Errorf("albumsById %d, %v", id, err)
	}
	return album, nil
}
