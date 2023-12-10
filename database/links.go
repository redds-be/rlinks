//    rlinks, a simple link shortener written in Go.
//    Copyright (C) 2023 redd
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

func CreateLinksTable(db *sql.DB, maxShort int) {
	sqlCreateTable := fmt.Sprintf("CREATE TABLE IF NOT EXISTS links (id UUID PRIMARY KEY, created_at TIMESTAMP NOT NULL, expire_at TIMESTAMP NOT NULL, url varchar NOT NULL, short varchar(%d) UNIQUE NOT NULL, password varchar(97));", maxShort)
	_, err := db.Exec(sqlCreateTable)
	if err != nil {
		log.Fatal("Unable to create the 'links' table:", err)
	}
}

func CreateLink(db *sql.DB, id uuid.UUID, createdAt time.Time, expireAt time.Time, url, short, password string) (Link, error) {
	sqlCreateLink := `INSERT INTO links (id, created_at, expire_at, url, short, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING expire_at, url, short;`
	var returnValues Link
	err := db.QueryRow(sqlCreateLink, id, createdAt, expireAt, url, short, password).Scan(
		&returnValues.ExpireAt,
		&returnValues.Url,
		&returnValues.Short,
	)

	return returnValues, err
}

func GetUrlByShort(db *sql.DB, short string) (string, error) {
	sqlGetUrlByShort := `SELECT url FROM links WHERE short = $1;`
	var url string
	err := db.QueryRow(sqlGetUrlByShort, short).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func GetHashByShort(db *sql.DB, short string) (string, error) {
	sqlGetPasswordByShort := `SELECT password FROM links WHERE short = $1;`
	var password string
	err := db.QueryRow(sqlGetPasswordByShort, short).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

func GetLinks(db *sql.DB) ([]Link, error) {
	sqlGetLinks := `SELECT expire_at, url, short FROM links;`
	var links []Link
	rows, err := db.Query(sqlGetLinks)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var i Link
		err := rows.Scan(
			&i.ExpireAt,
			&i.Url,
			&i.Short,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, i)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return links, nil
}

func RemoveLink(db *sql.DB, short string) error {
	sqlRemoveLink := `DELETE FROM links WHERE short = $1;`
	_, err := db.Exec(sqlRemoveLink, short)
	if err != nil {
		return err
	}

	return nil
}
