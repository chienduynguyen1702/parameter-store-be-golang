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
		organizationGroup.GET("/dashboard/logs", controllers.GetOrganizationInformation)
		organizationGroup.GET("/dashboard/totals", controllers.GetOrganizationDashboardTotals)
		organizationGroup.PUT("/:organization_id", middleware.RequiredIsOrgAdmin, controllers.UpdateOrganizationInformation)
	}
}
