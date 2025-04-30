package models

import (
	"errors"
	"time"
)

type Address struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	StreetLine1 string    `json:"street_line1"`
	StreetLine2 string    `json:"street_line2"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Get all addresses with optional limit
func GetAddresses(limit int) ([]Address, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, street_line1, street_line2, city, state, postal_code, country, 
		       is_default, created_at, updated_at
		FROM addresses
		ORDER BY id
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var addresses []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.StreetLine1, &a.StreetLine2, &a.City, 
				&a.State, &a.PostalCode, &a.Country, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

// Get address by ID
func GetAddressByID(id int) (Address, error) {
	var address Address
	err := DB.QueryRow(`
		SELECT id, user_id, street_line1, street_line2, city, state, postal_code, country, 
		       is_default, created_at, updated_at
		FROM addresses 
		WHERE id = ?`, id).Scan(
		&address.ID, &address.UserID, &address.StreetLine1, &address.StreetLine2, &address.City,
		&address.State, &address.PostalCode, &address.Country, &address.IsDefault, 
		&address.CreatedAt, &address.UpdatedAt)
	
	if err != nil {
		return address, err
	}
	
	return address, nil
}

// Get addresses by user ID
func GetAddressesByUserID(userID int) ([]Address, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, street_line1, street_line2, city, state, postal_code, country, 
		       is_default, created_at, updated_at
		FROM addresses
		WHERE user_id = ?
		ORDER BY is_default DESC, id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var addresses []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.StreetLine1, &a.StreetLine2, &a.City, 
				&a.State, &a.PostalCode, &a.Country, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

// Create a new address
func CreateAddress(userID int, streetLine1, streetLine2, city, state, postalCode, country string, isDefault bool) (int, error) {
	// If this is a default address, unset default flag on other addresses for this user
	if isDefault {
		_, err := DB.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ?", userID)
		if err != nil {
			return 0, err
		}
	}
	
	result, err := DB.Exec(`
		INSERT INTO addresses (user_id, street_line1, street_line2, city, state, postal_code, country, is_default, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		userID, streetLine1, streetLine2, city, state, postalCode, country, isDefault)
	if err != nil {
		return 0, err
	}
	
	id, err := result.LastInsertId()
	return int(id), err
}

// Update an existing address
func UpdateAddress(id, userID int, streetLine1, streetLine2, city, state, postalCode, country string, isDefault bool) error {
	// If this is a default address, unset default flag on other addresses for this user
	if isDefault {
		_, err := DB.Exec("UPDATE addresses SET is_default = 0 WHERE user_id = ? AND id != ?", userID, id)
		if err != nil {
			return err
		}
	}
	
	_, err := DB.Exec(`
		UPDATE addresses 
		SET street_line1 = ?, street_line2 = ?, city = ?, state = ?, postal_code = ?, country = ?, 
		    is_default = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND user_id = ?`,
		streetLine1, streetLine2, city, state, postalCode, country, isDefault, id, userID)
	return err
}

// Delete an address
func DeleteAddress(id, userID int) error {
	// Check if the address is referenced by any orders
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM orders WHERE address_id = ?", id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return ErrAddressInUse
	}
	
	_, err = DB.Exec("DELETE FROM addresses WHERE id = ? AND user_id = ?", id, userID)
	return err
}

// Custom error
var ErrAddressInUse = errors.New("address is in use by one or more orders")