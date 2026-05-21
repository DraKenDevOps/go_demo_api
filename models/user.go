package models

import (
	"time"
)

type UserTable struct {
	ID            uint       `json:"id"`
	Telephone     string     `json:"telephone"`
	Name          string     `json:"name"`
	Password      string     `json:"password"`
	ReportGroupID *uint      `json:"reportGroupId"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     *time.Time `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

type AuthUser struct {
	ID            uint   `json:"id"`
	Telephone     string `json:"telephone"`
	Name          string `json:"name"`
	Password      string `json:"-"`
	ReportGroupID *uint  `json:"reportGroupId"`
	Role          string `json:"role"`
	Status        string `json:"-"`
}

type User struct {
	ID            uint    `json:"id"`
	Telephone     string  `json:"telephone"`
	Name          string  `json:"name"`
	ReportGroupID *uint   `json:"reportGroupId"`
	Role          string  `json:"role"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     *string `json:"updatedAt"`
	DeletedAt     *string `json:"deletedAt"`
}

type SaveUser struct {
	Telephone string `json:"telephone"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

type SaveUserPassword struct {
	Password string `json:"password"`
}

type JwtUserClaim struct {
	ID            uint   `json:"id"`
	Telephone     string `json:"telephone"`
	Name          string `json:"name"`
	ReportGroupID *uint  `json:"reportGroupId"`
	Role          string `json:"role"`
}
