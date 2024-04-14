package routes

import (
	"parameter-store-be/controllers"

	"github.com/gin-gonic/gin"
)

func setupGroupAgent(r *gin.RouterGroup) {
	agentGroup := r.Group("/agents")
	{
		agentGroup.POST("/:agent_id/rerun-workflow", controllers.RerunWorkFlowByAgent)
	}
}
