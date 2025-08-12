package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Not included in JSON responses
	CreatedAt time.Time `json:"created_at"`
}

// UserModel handles user database operations
type UserModel struct {
	DB *sql.DB
}

// NewUserModel creates a new UserModel
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// CreateTable creates the users table if it doesn't exist
func (m *UserModel) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := m.DB.Exec(query)
	return err
}

// Create creates a new user with hashed password
func (m *UserModel) Create(username, password string) (*User, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Insert user
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	result, err := m.DB.Exec(query, username, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	// Get the created user
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return m.GetByID(int(id))
}

// GetByUsername retrieves a user by username
func (m *UserModel) GetByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, created_at FROM users WHERE username = ?`
	err := m.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID retrieves a user by ID
func (m *UserModel) GetByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, created_at FROM users WHERE id = ?`
	err := m.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ValidatePassword checks if the provided password matches the user's password
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
