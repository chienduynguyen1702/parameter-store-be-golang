package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupProject(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects/:project_id", middleware.RequiredAuth)
	{
		overviewGroup := projectGroup.Group("/overview")
		{
			overviewGroup.GET("/", controllers.GetProjectOverView)
			overviewGroup.PUT("/", controllers.UpdateProjectInformation)
		}
		agentGroup := projectGroup.Group("/agents")
		{
			agentGroup.GET("/", controllers.GetProjectAgents)
			agentGroup.POST("/", controllers.CreateNewAgent)
			// agentGroup.PUT("/:agent_id", controllers.UpdateProjectInformation)
			// agentGroup.DELETE("/:agent_id", controllers.DeleteAgent)
		}
		stageGroup := projectGroup.Group("/stages")
		{
			stageGroup.GET("/", controllers.GetStages)
			stageGroup.POST("/", controllers.CreateStage)
			stageGroup.PUT("/:stage_id", controllers.UpdateStage)
			stageGroup.DELETE("/:stage_id", controllers.DeleteStage)
		}
		environmentGroup := projectGroup.Group("/environments")
		{
			environmentGroup.GET("/", controllers.GetEnvironments)
			environmentGroup.POST("/", controllers.CreateEnvironment)
			environmentGroup.PUT("/:environment_id", controllers.UpdateEnvironment)
			environmentGroup.DELETE("/:environment_id", controllers.DeleteEnvironment)
		}
		// parameterGroup := projectGroup.Group("/parameters")
		// {
		// parameterGroup.GET("/", controllers.GetProjectParameters)
		// parameterGroup.POST("/", controllers.CreateNewParameter)
		// parameterGroup.PUT("/:parameter_id", controllers.UpdateParameter)
		// parameterGroup.DELETE("/:parameter_id", controllers.DeleteParameter)
		// }
		// trackingGroup := projectGroup.Group("/tracking")
		// {
		// trackingGroup.GET("/", controllers.GetProjectTracking)
		// trackingGroup.POST("/", controllers.CreateNewTracking)
		// trackingGroup.PUT("/:tracking_id", controllers.UpdateTracking)
		// trackingGroup.DELETE("/:tracking_id", controllers.DeleteTracking)
		// }
	}
}
