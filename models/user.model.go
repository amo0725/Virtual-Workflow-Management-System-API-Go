package models

import (
	"virtual_workflow_management_system_gin/common"
)

type UserRole string

const (
	Admin    UserRole = "Admin"
	Employer UserRole = "Employer"
)

type UserAction string

const (
	Edit     UserAction = "Edit"
	Delete   UserAction = "Delete"
	Transfer UserAction = "Transfer"
)

type User struct {
	common.BaseModel `bson:",inline"`
	Username         string   `json:"username" bson:"username"`
	Password         string   `json:"password" bson:"password"`
	Role             UserRole `json:"role" bson:"role"`
}

type JWTUser struct {
	Username string   `json:"username"`
	Role     UserRole `json:"role"`
}
