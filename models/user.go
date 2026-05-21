package models

import (
	"time"
)

type UserTable struct {
	UserId     uint       `json:"user_id"`
	Username   string     `json:"username"`
	Telephone  string     `json:"telephone"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Level      string     `json:"level"`
	RoleAction string     `json:"role_action"`
	OpId       *uint      `json:"op_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type AuthUser struct {
	UserId     uint   `json:"user_id"`
	Telephone  string `json:"telephone"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	RoleAction string `json:"role_action"`
	Password   string `json:"-"`
	OpId       *uint  `json:"op_id"`
	Level      string `json:"level"`
	Status     string `json:"-"`
}

type User struct {
	UserId     uint    `json:"user_id"`
	Telephone  string  `json:"telephone"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	OpId       *uint   `json:"op_id"`
	Level      string  `json:"level"`
	RoleAction string  `json:"role_action"`
	Status     string  `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
	DeletedAt  *string `json:"deleted_at"`
}

type SaveUser struct {
	Telephone  string  `json:"telephone"`
	Email      string  `json:"email"`
	Username   string  `json:"username"`
	Level      string  `json:"level"`
	RoleAction string  `json:"role_action"`
	Status     *string `json:"status,omitempty"`
	Password   *string `json:"password,omitempty"`
}

type SaveUserPassword struct {
	Password string `json:"password"`
}

type JwtUserClaim struct {
	UserId     uint   `json:"user_id"`
	Telephone  string `json:"telephone"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	OpId       *uint  `json:"op_id"`
	Level      string `json:"level"`
	RoleAction string `json:"role_action"`
}

func (User) TableName() string {
	return "users"
}

func (AuthUser) TableName() string {
	return "users"
}

func (UserTable) TableName() string {
	return "users"
}

func (SaveUser) TableName() string {
	return "users"
}

func (SaveUserPassword) TableName() string {
	return "users"
}
