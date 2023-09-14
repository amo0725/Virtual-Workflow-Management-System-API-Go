package models

import (
	"virtual_workflow_management_system_gin/common"
)

type TaskStatus string

const (
	Pending    TaskStatus = "Pending"
	InProgress TaskStatus = "In Progress"
	Completed  TaskStatus = "Completed"
)

type Task struct {
	common.BaseModel `bson:",inline"`
	Name             string     `json:"name" bson:"name"`
	Description      string     `json:"description" bson:"description"`
	Status           TaskStatus `json:"status" bson:"status"`
	Order            int        `json:"order" bson:"order"`
}

type Workflow struct {
	common.BaseModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Tasks            []Task `json:"tasks" bson:"tasks"`
	Owner            string `json:"owner" bson:"owner"`
}

func (workflow *Workflow) CheckWorkflowAccess(user JWTUser, action UserAction) bool {
	if user.Username == workflow.Owner {
		return true
	}

	switch action {
	case "transfer":
		return user.Username == workflow.Owner
	case "edit":
	case "delete":
		return user.Role == "admin"
	default:
		return false
	}

	return false
}
