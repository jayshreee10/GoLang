package models

import (
	"database/sql"
	"errors"
	"time"
	"go-crud/auth"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Password is never sent to the client
	CreatedAt time.Time `json:"created_at"`
}

// UserSignup represents the data needed for registration
type UserSignup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLogin represents the data needed for login
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response for successful authentication
type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

// Custom errors
var (
	ErrUserExists   = errors.New("user with this email already exists")
	ErrInvalidLogin = errors.New("invalid email or password")
)

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

func GetUserByEmail(email string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, email, password, created_at FROM users WHERE email = ?", email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

// CreateUser creates a new user with plain text password
func CreateUser(email string) (int, error) {
	result, err := DB.Exec("INSERT INTO users (email) VALUES (?)", email)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

// RegisterUser creates a new user with hashed password
func RegisterUser(email, password string) (int, error) {
	// Check if user already exists
	var exists int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&exists)
	if err != nil {
		return 0, err
	}
	if exists > 0 {
		return 0, ErrUserExists
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return 0, err
	}

	// Create the user
	result, err := DB.Exec("INSERT INTO users (email, password) VALUES (?, ?)", email, hashedPassword)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	return int(id), err
}

// LoginUser validates credentials and returns user if valid
func LoginUser(email, password string) (User, error) {
	var user User
	err := DB.QueryRow("SELECT id, email, password, created_at FROM users WHERE email = ?", email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrInvalidLogin
		}
		return User{}, err
	}

	// Check password
	if err := auth.CheckPassword(user.Password, password); err != nil {
		return User{}, ErrInvalidLogin
	}

	return user, nil
}

func UpdateUser(id int, email string) error {
	_, err := DB.Exec("UPDATE users SET email = ? WHERE id = ?", email, id)
	return err
}

func UpdateUserPassword(id int, password string) error {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	
	_, err = DB.Exec("UPDATE users SET password = ? WHERE id = ?", hashedPassword, id)
	return err
}

func DeleteUser(id int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}