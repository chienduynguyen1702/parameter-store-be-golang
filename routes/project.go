package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupProject(r *gin.RouterGroup) {
	projectGroup := r.Group("/projects/:project_id", middleware.RequiredAuth, middleware.RequiredBelongToProject)
	{
		projectGroup.GET("/", controllers.GetProjectAllInfo)
		overviewGroup := projectGroup.Group("/overview")
		{
			overviewGroup.GET("/", controllers.GetProjectOverView)
			overviewGroup.PUT("/", controllers.UpdateProjectInformation, middleware.RequiredIsAdmin)
			overviewGroup.POST("/add-user", controllers.AddUserToProject, middleware.RequiredIsAdmin)
		}
		agentGroup := projectGroup.Group("/agents")
		{
			agentGroup.GET("/", controllers.GetAgents)
			agentGroup.POST("/", controllers.CreateNewAgent, middleware.RequiredIsAdmin)
			agentGroup.GET("/:agent_id", controllers.GetAgentDetail)
			agentGroup.PATCH("/:agent_id/archive", controllers.ArchiveAgent, middleware.RequiredIsAdmin)
			agentGroup.PATCH("/:agent_id/unarchive", controllers.RestoreAgent, middleware.RequiredIsAdmin)
			agentGroup.GET("/archived", controllers.GetArchivedAgents)
			// agentGroup.PUT("/:agent_id", controllers.UpdateProjectInformation)
			// agentGroup.DELETE("/:agent_id", controllers.DeleteAgent)
		}
		versionGroup := projectGroup.Group("/versions")
		{
			versionGroup.GET("/", controllers.GetProjectVersions)
			versionGroup.POST("/", controllers.CreateNewVersion, middleware.RequiredIsAdmin)

		}
		parameterGroup := projectGroup.Group("/parameters")
		{
			parameterGroup.GET("/", controllers.GetProjectParameters)
			parameterGroup.POST("/", controllers.CreateParameter, middleware.RequiredIsAdmin)
			parameterGroup.PUT("/:parameter_id", controllers.UpdateParameter, middleware.RequiredIsAdmin)

			parameterGroup.GET("/:parameter_id", controllers.GetParameterByID)

			parameterGroup.GET("/archived", controllers.GetArchivedParameters)
			parameterGroup.PATCH("/:parameter_id/archive", controllers.ArchiveParameter, middleware.RequiredIsAdmin)
			parameterGroup.PATCH("/:parameter_id/unarchive", controllers.UnarchiveParameter, middleware.RequiredIsAdmin)
		}
		trackingGroup := projectGroup.Group("/tracking")
		{
			trackingGroup.GET("/", controllers.GetProjectTracking)
			// trackingGroup.POST("/", controllers.CreateNewTracking)
			// trackingGroup.PUT("/:tracking_id", controllers.UpdateTracking)
			// trackingGroup.DELETE("/:tracking_id", controllers.DeleteTracking)
		}
	}
}
