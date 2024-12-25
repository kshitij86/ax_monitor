package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
)

func CountUsers(db *sql.DB) int {
	rows, err := db.Query("select * from ax_monitor_users")
	if err != nil {
		panic(err)
	}

	var users []User

	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.user_id, &usr.username, &usr.password); err != nil {
			panic(err)
		}
		users = append(users, usr)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
	return len(users)
}

func LoginUser(db *sql.DB, username string) (string, bool, error) {

	// find the user
	var foundUser User
	if err := db.QueryRow("SELECT * FROM ax_monitor_users WHERE username = ?", username).Scan(&foundUser.user_id, &foundUser.username, &foundUser.password); err != nil {
		if err == sql.ErrNoRows {
			return "", false, fmt.Errorf("no such user found")
		}
		return "", false, fmt.Errorf("%s", err.Error())
	}
	fmt.Println("user found!!! -> ", foundUser.username)

	var password string
	fmt.Println("enter password: ")
	fmt.Scan(&password)

	hashMaker := sha256.New()
	hashMaker.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hashMaker.Sum(nil))

	var loginSuccess = (hashedPassword == foundUser.password)

	if loginSuccess {
		fmt.Println("user found...")
		fmt.Println("applying settings and importing projects for this user...")
		return foundUser.user_id, true, nil
	}
	return "", false, fmt.Errorf("user not found")
}
