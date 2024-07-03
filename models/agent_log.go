package models

import (
	"gorm.io/gorm"
)

type AgentLog struct {
	gorm.Model
	AgentID        uint   `gorm:"foreignKey:AgentID" json:"agent_id"`
	Action         string `gorm:"type:varchar(100);not null" json:"action"`
	ProjectID      uint   `gorm:"foreignKey:ProjectID" json:"project_id"`
	Path           string `gorm:"type:varchar(100);not null" json:"path"`
	ResponseStatus int    `gorm:"not null" json:"response_status"`
	Message        string `gorm:"type:text" json:"message"`
	Latency        int    `gorm:"not null" json:"latency"`

	ExecutedInWorkflowLogID uint        `json:"executed_in_workflow_log_id"`
	ExecutedInWorkflowLog   WorkflowLog `gorm:"foreignKey:ExecutedInWorkflowLogID; references:ID" json:"executed_in_workflow_log"`

	AgentPullParameterLog []AgentPullParameterLog `gorm:"foreignKey:AgentLogID; references:ID" json:"agent_pull_parameter_log"`

	Project Project `json:"project"`
	Agent   Agent   `json:"agent"`
}
