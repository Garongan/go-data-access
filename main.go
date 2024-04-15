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
	logError(err)
	fmt.Printf("Albums found: %v\n", albums)

	// get data by id
	album, err := getAlbumById(1)
	logError(err)
	fmt.Printf("Album found: %v\n", album)

	// add data album
	addInfo, err := addAlbum(Album{
		Title:  "Dumes",
		Artist: "Denny Chacknan",
		Price:  55.65,
	})
	logError(err)
	fmt.Println(addInfo)

	// update album by id
	updatedAlbum, err := updateAlbumById(5, Album{
		Title: "Kisinan",
		Artist: "Massdho",
		Price: 66.54,
	})
	logError(err)
	fmt.Printf("Updated Album: %v\n", updatedAlbum)

	// delete album by id
	deleteInfo, err := deleteAlbumById(5)
	logError(err)
	fmt.Println(deleteInfo)

	// get all album
	allAlbums, err := getAlbums()
	logError(err)
	fmt.Printf("All Albums: %v\n", allAlbums)

}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// get albums by artist name
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

// get albums by id
func getAlbumById(id int64) (Album, error) {
	var album Album
	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("getAlbumById %d, no such album", id)
		}
		return album, fmt.Errorf("getAlbumById %d, %v", id, err)
	}
	return album, nil
}

// add new album
func addAlbum(album Album) (string, error) {
	_, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)", album.Title, album.Artist, album.Price)
	if err != nil {
		return "", fmt.Errorf("addAlbum %v", err)
	}
	return "success created new album", nil
}

// get all albums
func getAlbums() ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("getAlbums: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("getAlbums: %v", err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAlbums: %v", err)
	}
	return albums, nil
}

// delete albums by id
func deleteAlbumById(id int64) (string, error) {
	_, err := db.Exec("DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return "", fmt.Errorf("deleteAlbumById: %v", err)
	}
	return fmt.Sprintf("success delete album of id: %d", id), nil
}

// update album by id
func updateAlbumById(id int64, album Album) (Album, error) {

	_, err := db.Exec("UPDATE album SET title = $1, artist = $2, price = $3 WHERE id = $4",
		album.Title, album.Artist, album.Price, id,
	)
	if err != nil {
		return album, fmt.Errorf("updateAlbumById: %v", err)
	}

	return getAlbumById(id)
}
