package models

import (
	
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUsers(limit int) ([]User, error) {
	rows, err := DB.Query("SELECT id, email, created_at FROM users LIMIT ?", limit)
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

func GetUserByID(id int) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, email, created_at FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Email, &user.CreatedAt)
	return user, err
}

func CreateUser(email string) (int, error) {
	result, err := DB.Exec("INSERT INTO users (email) VALUES (?)", email)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateUser(id int, email string) error {
	_, err := DB.Exec("UPDATE users SET email = ? WHERE id = ?", email, id)
	return err
}

func DeleteUser(id int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}