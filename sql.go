package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Song struct {
	id      int
	name    string
	cat     string // Is either English1, Kannada, Nyishi Songs, Spanish, Sunday School, VV Hindi 2021, VV Malayalam 2021, VV Tamil 2021
	font    string
	font2   string
	lyrics  string
	lyrics2 string
}

func (s Song) print() {
	fmt.Printf("ID: %d, NAME: %s, CAT: %s, FONT: %s, FONT2: %s, LYRICS: %s, LYRICS2: %s\n",
		s.id, s.name, s.cat,
		s.font, s.font2,
		s.lyrics, s.lyrics2)
}

func lookupSong(songs []Song, title string) Song{
	for _, song := range songs {
		if song.name == title {
			return song
		}
	}
	return Song{}
}

func db_query(args ...string) []Song {
	var songs []Song

	db, err := sql.Open("sqlite3", "songs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var rows *sql.Rows

	switch len(args) {
	case 0:
		rows, err = db.Query("SELECT id, name, cat, font, font2, lyrics, lyrics2 from sm")
	case 1:
		pattern := "%"+args[0]+"%"
		rows, err = db.Query(
			`SELECT id, name, cat, font, font2, lyrics, lyrics2 from sm WHERE (lyrics like ? OR lyrics2 like ? OR name like ?)`, pattern, pattern, pattern)

	case 2:
		pattern := "%"+args[0]+"%"
		rows, err = db.Query(
			`SELECT id, name, cat, font, font2, lyrics, lyrics2 from sm WHERE (lyrics like ? OR lyrics2 like ? OR name like ?) AND cat = ?`, pattern, pattern, pattern, args[1])
	}
	
	for rows.Next() {
		var s Song
		
		err := rows.Scan(&s.id, &s.name, &s.cat, &s.font, &s.font2, &s.lyrics, &s.lyrics2)

		if err != nil {
			log.Fatal(err)
		}
		songs = append(songs, s)
	}

	return songs
}


func test_sql() {
	s := db_query("g", "English1")
	for _, song := range s {
		song.print()
	}
}
