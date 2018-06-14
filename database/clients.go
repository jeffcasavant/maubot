// jesaribot - A simple maubot plugin.
// Copyright (C) 2018 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	log "maunium.net/go/maulogger"
	"database/sql"
)

type MatrixClient struct {
	db  *Database
	sql *sql.DB

	UserID      string `json:"user_id"`
	Homeserver  string `json:"homeserver"`
	AccessToken string `json:"access_token"`
	NextBatch   string `json:"next_batch"`
	FilterID    string `json:"filter_id"`

	Sync          bool   `json:"sync"`
	AutoJoinRooms bool   `json:"auto_join_rooms"`
	DisplayName   string `json:"display_name"`
	AvatarURL     string `json:"avatar_url"`
}

type MatrixClientStatic struct {
	db  *Database
	sql *sql.DB
}

func (mcs *MatrixClientStatic) CreateTable() error {
	_, err := mcs.sql.Exec(`CREATE TABLE IF NOT EXISTS matrix_client (
		user_id      VARCHAR(255) PRIMARY KEY,
		homeserver   VARCHAR(255) NOT NULL,
		access_token VARCHAR(255) NOT NULL,
		next_batch   VARCHAR(255) NOT NULL,
		filter_id    VARCHAR(255) NOT NULL,

		sync         BOOLEAN      NOT NULL,
		autojoin     BOOLEAN      NOT NULL,
		display_name VARCHAR(255) NOT NULL,
		avatar_url   VARCHAR(255) NOT NULL
	)`)
	return err
}

func (mcs *MatrixClientStatic) Get(userID string) *MatrixClient {
	row := mcs.sql.QueryRow("SELECT user_id, homeserver, access_token, next_batch, filter_id, sync, autojoin, display_name, avatar_url FROM matrix_client WHERE user_id=?", userID)
	if row != nil {
		return mcs.New().Scan(row)
	}
	return nil
}

func (mcs *MatrixClientStatic) GetAll() (clients []*MatrixClient) {
	rows, err := mcs.sql.Query("SELECT user_id, homeserver, access_token, next_batch, filter_id, sync, autojoin, display_name, avatar_url FROM matrix_client")
	if err != nil || rows == nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		clients = append(clients, mcs.New().Scan(rows))
	}
	return
}

func (mcs *MatrixClientStatic) New() *MatrixClient {
	return &MatrixClient{
		db:  mcs.db,
		sql: mcs.sql,
	}
}

type Scannable interface {
	Scan(...interface{}) error
}

func (mxc *MatrixClient) Scan(row Scannable) *MatrixClient {
	err := row.Scan(&mxc.UserID, &mxc.Homeserver, &mxc.AccessToken, &mxc.NextBatch, &mxc.FilterID, &mxc.Sync, &mxc.AutoJoinRooms, &mxc.DisplayName, &mxc.AvatarURL)
	if err != nil {
		log.Fatalln("Database scan failed:", err)
	}
	return mxc
}

func (mxc *MatrixClient) Insert() error {
	_, err := mxc.sql.Exec("INSERT INTO matrix_client (user_id, homeserver, access_token, next_batch, filter_id, sync, autojoin, display_name, avatar_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		mxc.UserID, mxc.Homeserver, mxc.AccessToken, mxc.NextBatch, mxc.FilterID, mxc.Sync, mxc.AutoJoinRooms, mxc.DisplayName, mxc.AvatarURL)
	return err
}

func (mxc *MatrixClient) Update() error {
	_, err := mxc.sql.Exec("UPDATE matrix_client SET access_token=?, next_batch=?, filter_id=?, sync=?, autojoin=?, display_name=?, avatar_url=? WHERE user_id=?",
		mxc.AccessToken, mxc.NextBatch, mxc.FilterID, mxc.Sync, mxc.AutoJoinRooms, mxc.DisplayName, mxc.AvatarURL, mxc.UserID)
	return err
}