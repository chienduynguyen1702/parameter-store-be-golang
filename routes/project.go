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
		projectGroup.POST("/apply-parameters", middleware.RequiredIsAdmin, controllers.ApplyParametersInProject)
		overviewGroup := projectGroup.Group("/overview")
		{
			overviewGroup.GET("/", controllers.GetProjectOverView)
			overviewGroup.PUT("/", middleware.RequiredIsAdmin, controllers.UpdateProjectInformation)
			overviewGroup.POST("/add-user", middleware.RequiredIsAdmin, controllers.AddUserToProject)
			overviewGroup.POST("/remove-user", middleware.RequiredIsAdmin, controllers.RemoveUserFromProject)
			overviewGroup.GET("/users/:user_id", controllers.GetUserInProject)
			overviewGroup.PUT("/users/:user_id", controllers.UpdateUserInProject)
			overviewGroup.DELETE("/users/:user_id", middleware.RequiredIsAdmin, controllers.RemoveUserFromProject)
		}
		workflowGroup := projectGroup.Group("/workflows")
		{
			workflowGroup.GET("/", controllers.GetProjectWorkflows)
			workflowGroup.GET("/:workflow_id/run", controllers.GetWorkflowProcess)
		}
		dashboardGrop := projectGroup.Group("/dashboard")
		{
			dashboardGrop.GET("/totals", controllers.GetProjectDashboardTotals)
			dashboardGrop.GET("/logs", controllers.GetProjectDashboardLogs)
		}
		agentGroup := projectGroup.Group("/agents")
		{
			agentGroup.GET("/", controllers.GetAgents)
			agentGroup.POST("/", middleware.RequiredIsAdmin, controllers.CreateNewAgent)
			agentGroup.GET("/:agent_id", controllers.GetAgentDetail)
			agentGroup.PATCH("/:agent_id/archive", middleware.RequiredIsAdmin, controllers.ArchiveAgent)
			agentGroup.PATCH("/:agent_id/unarchive", middleware.RequiredIsAdmin, controllers.RestoreAgent)
			agentGroup.GET("/archived", controllers.GetArchivedAgents)
			agentGroup.PUT("/:agent_id", middleware.RequiredIsAdmin, controllers.UpdateAgent)
			// agentGroup.DELETE("/:agent_id", controllers.DeleteAgent)
		}
		versionGroup := projectGroup.Group("/versions")
		{
			versionGroup.GET("/", controllers.GetProjectVersions)
			versionGroup.POST("/", middleware.RequiredIsAdmin, controllers.CreateNewVersion)

		}
		parameterGroup := projectGroup.Group("/parameters")
		{
			parameterGroup.GET("/", controllers.GetProjectParameters)
			parameterGroup.GET("/download", controllers.DownloadLatestParameters)
			parameterGroup.POST("/", middleware.RequiredIsAdmin, controllers.CreateParameter)
			parameterGroup.PUT("/:parameter_id", middleware.RequiredIsAdmin, controllers.UpdateParameter)

			parameterGroup.GET("/:parameter_id", controllers.GetParameterByID)

			parameterGroup.GET("/archived", controllers.GetArchivedParameters)
			parameterGroup.PATCH("/:parameter_id/archive", middleware.RequiredIsAdmin, controllers.ArchiveParameter)
			parameterGroup.PATCH("/:parameter_id/unarchive", middleware.RequiredIsAdmin, controllers.UnarchiveParameter)

			parameterGroup.GET("/download-template", controllers.DownloadExecelTemplateParameters)
			parameterGroup.POST("/upload", middleware.RequiredIsAdmin, controllers.UploadParameters)
			parameterGroup.GET("/:parameter_id/search-in-repo", controllers.SearchParameterInRepo)
			parameterGroup.GET("/:parameter_id/get-file-content", controllers.TestGetFileContent)

			parameterGroup.POST("/check-using", controllers.CheckParameterUsing)
		}
		trackingGroup := projectGroup.Group("/tracking")
		{
			trackingGroup.GET("/logs", controllers.GetProjectTracking)
			// trackingGroup.POST("/", controllers.CreateNewTracking)
			// trackingGroup.PUT("/:tracking_id", controllers.UpdateTracking)
			// trackingGroup.DELETE("/:tracking_id", controllers.DeleteTracking)
		}
		stageGroup := projectGroup.Group("/stages")
		{
			stageGroup.GET("/", controllers.GetListStageInProject)
			stageGroup.POST("/", middleware.RequiredIsAdmin, controllers.CreateStageInProject)
			stageGroup.PUT("/:stage_id", middleware.RequiredIsAdmin, controllers.UpdateStageInProject)

			stageGroup.GET("/:stage_id", controllers.GetStageInProject)

			stageGroup.GET("/archived", controllers.GetListArchivedStageInProject)
			stageGroup.PATCH("/:stage_id/archive", middleware.RequiredIsAdmin, controllers.ArchiveStageInProject)
			stageGroup.PATCH("/:stage_id/unarchive", middleware.RequiredIsAdmin, controllers.UnarchiveStageInProject)
		}
		environmentGroup := projectGroup.Group("/environments")
		{
			environmentGroup.GET("/", controllers.GetListEnvironmentInProject)
			environmentGroup.POST("/", middleware.RequiredIsAdmin, controllers.CreateEnvironmentInProject)
			environmentGroup.PUT("/:environment_id", middleware.RequiredIsAdmin, controllers.UpdateEnvironmentInProject)

			environmentGroup.GET("/:environment_id", controllers.GetEnvironmentInProject)

			environmentGroup.GET("/archived", controllers.GetListArchivedEnvironmentInProject)
			environmentGroup.PATCH("/:environment_id/archive", middleware.RequiredIsAdmin, controllers.ArchiveEnvironmentInProject)
			environmentGroup.PATCH("/:environment_id/unarchive", middleware.RequiredIsAdmin, controllers.UnarchiveEnvironmentInProject)
		}
		// githubGroup := projectGroup.Group("/github")
		// {
		// githubGroup.POST("/repositories", controllers.ListRepositoryByGithubUser)
		// githubGroup.GET("/repos", controllers.GetProjectGithubRepos)
		// githubGroup.GET("/repos/:repo_id", controllers.GetProjectGithubRepoDetail)
		// githubGroup.POST("/repos/:repo_id/apply", middleware.RequiredIsAdmin, controllers.ApplyParametersFromGithubRepo)
		// }
	}
	r.POST("/projects/github/repositories", controllers.ListRepositoryByGithubUser)
}
