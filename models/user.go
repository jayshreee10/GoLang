package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}


var DB *sql.DB

func GetUsers(limit int) ([]User, error) {
	rows, err := DB.Query("SELECT id, email, created_at FROM spree_users LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func CreateUser(email string) (int, error) {
	var id int
	err := DB.QueryRow(`INSERT INTO spree_users (email, created_at, updated_at) VALUES ($1, NOW(), NOW()) RETURNING id`, email).Scan(&id)
	return id, err
}

func UpdateUser(id int, email string) error {
	_, err := DB.Exec(`UPDATE spree_users SET email = $1 WHERE id = $2`, email, id)
	return err
}

func DeleteUser(id int) error {
	_, err := DB.Exec(`DELETE FROM spree_users WHERE id = $1`, id)
	return err
}
