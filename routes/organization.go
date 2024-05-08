package routes

import (
	"parameter-store-be/controllers"
	"parameter-store-be/middleware"

	"github.com/gin-gonic/gin"
)

func setupGroupOrganization(r *gin.RouterGroup) {
	organizationGroup := r.Group("/organizations", middleware.RequiredAuth)
	{
		organizationGroup.GET("/", controllers.GetOrganizationInformation)
		organizationGroup.GET("/dashboard/logs", middleware.RequiredIsOrgAdmin, controllers.GetOrganizationDashboardLogs)
		organizationGroup.GET("/dashboard/totals", middleware.RequiredIsOrgAdmin, controllers.GetOrganizationDashboardTotals)
		organizationGroup.PUT("/:organization_id", middleware.RequiredIsOrgAdmin, controllers.UpdateOrganizationInformation)
	}
}
