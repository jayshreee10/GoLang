package models

import (
	"time"
)

type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func GetRoles(limit int) ([]Role, error) {
	rows, err := DB.Query("SELECT id, name, description, created_at FROM roles LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name, &r.Description, &r.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func GetRoleByID(id int) (Role, error) {
	var role Role
	err := DB.QueryRow("SELECT id, name, description, created_at FROM roles WHERE id = ?", id).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	return role, err
}

func CreateRole(name string, description string) (int, error) {
	result, err := DB.Exec(
		"INSERT INTO roles (name, description) VALUES (?, ?)", 
		name, description)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

func UpdateRole(id int, name string, description string) error {
	_, err := DB.Exec(
		"UPDATE roles SET name = ?, description = ? WHERE id = ?", 
		name, description, id)
	return err
}

func DeleteRole(id int) error {
	_, err := DB.Exec("DELETE FROM roles WHERE id = ?", id)
	return err
}