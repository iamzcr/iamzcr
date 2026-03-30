package main

import (
	"iamzcr/config"
	"iamzcr/handlers/admin"
	"iamzcr/middleware"
	"iamzcr/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	models.InitDB(cfg)

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	api := r.Group("/api")
	{
		adminHandler := admin.NewAdminHandler()
		api.POST("/login", adminHandler.Login)

		adminGroup := api.Group("")
		adminGroup.Use(middleware.AuthMiddleware())
		{
			adminGroup.POST("/logout", adminHandler.Logout)
			adminGroup.GET("/admin/info", adminHandler.GetAdminInfo)

			adminGroup.GET("/articles", adminHandler.ListArticles)
			adminGroup.GET("/articles/:id", adminHandler.GetArticle)
			adminGroup.POST("/articles", adminHandler.CreateArticle)
			adminGroup.PUT("/articles/:id", adminHandler.UpdateArticle)
			adminGroup.DELETE("/articles/:id", adminHandler.DeleteArticle)

			categoryHandler := admin.NewCategoryHandler()
			adminGroup.GET("/categories", categoryHandler.List)
			adminGroup.GET("/categories/:id", categoryHandler.Get)
			adminGroup.POST("/categories", categoryHandler.Create)
			adminGroup.PUT("/categories/:id", categoryHandler.Update)
			adminGroup.DELETE("/categories/:id", categoryHandler.Delete)

			commentHandler := admin.NewCommentHandler()
			adminGroup.GET("/comments", commentHandler.List)
			adminGroup.GET("/comments/:id", commentHandler.Get)
			adminGroup.POST("/comments", commentHandler.Create)
			adminGroup.PUT("/comments/:id", commentHandler.Update)
			adminGroup.DELETE("/comments/:id", commentHandler.Delete)

			menuHandler := admin.NewMenuHandler()
			adminGroup.GET("/menus", menuHandler.List)
			adminGroup.GET("/menus/:id", menuHandler.Get)
			adminGroup.POST("/menus", menuHandler.Create)
			adminGroup.PUT("/menus/:id", menuHandler.Update)
			adminGroup.DELETE("/menus/:id", menuHandler.Delete)

			tagsHandler := admin.NewTagsHandler()
			adminGroup.GET("/tags", tagsHandler.List)
			adminGroup.GET("/tags/:id", tagsHandler.Get)
			adminGroup.POST("/tags", tagsHandler.Create)
			adminGroup.PUT("/tags/:id", tagsHandler.Update)
			adminGroup.DELETE("/tags/:id", tagsHandler.Delete)

			directoryHandler := admin.NewDirectoryHandler()
			adminGroup.GET("/directories", directoryHandler.List)
			adminGroup.GET("/directories/:id", directoryHandler.Get)
			adminGroup.POST("/directories", directoryHandler.Create)
			adminGroup.PUT("/directories/:id", directoryHandler.Update)
			adminGroup.DELETE("/directories/:id", directoryHandler.Delete)

			websiteHandler := admin.NewWebsiteHandler()
			adminGroup.GET("/website", websiteHandler.Get)
			adminGroup.PUT("/website", websiteHandler.Update)

			attachHandler := admin.NewAttachHandler()
			adminGroup.GET("/attaches", attachHandler.List)
			adminGroup.GET("/attaches/:id", attachHandler.Get)
			adminGroup.POST("/attaches", attachHandler.Create)
			adminGroup.PUT("/attaches/:id", attachHandler.Update)
			adminGroup.DELETE("/attaches/:id", attachHandler.Delete)

			langHandler := admin.NewLangHandler()
			adminGroup.GET("/langs", langHandler.List)
			adminGroup.GET("/langs/:id", langHandler.Get)
			adminGroup.POST("/langs", langHandler.Create)
			adminGroup.PUT("/langs/:id", langHandler.Update)
			adminGroup.DELETE("/langs/:id", langHandler.Delete)

			logHandler := admin.NewLogHandler()
			adminGroup.GET("/logs", logHandler.List)
			adminGroup.GET("/logs/:id", logHandler.Get)
			adminGroup.POST("/logs", logHandler.Create)
			adminGroup.DELETE("/logs/:id", logHandler.Delete)

			messageHandler := admin.NewMessageHandler()
			adminGroup.GET("/messages", messageHandler.List)
			adminGroup.GET("/messages/:id", messageHandler.Get)
			adminGroup.POST("/messages", messageHandler.Create)
			adminGroup.PUT("/messages/:id", messageHandler.Update)
			adminGroup.DELETE("/messages/:id", messageHandler.Delete)

			permitHandler := admin.NewPermitHandler()
			adminGroup.GET("/permits", permitHandler.List)
			adminGroup.GET("/permits/:id", permitHandler.Get)
			adminGroup.POST("/permits", permitHandler.Create)
			adminGroup.PUT("/permits/:id", permitHandler.Update)
			adminGroup.DELETE("/permits/:id", permitHandler.Delete)

			readHandler := admin.NewReadHandler()
			adminGroup.GET("/reads", readHandler.List)
			adminGroup.GET("/reads/:id", readHandler.Get)
			adminGroup.POST("/reads", readHandler.Create)
			adminGroup.DELETE("/reads/:id", readHandler.Delete)

			adminGroup.GET("/admins", adminHandler.ListAdmins)
			adminGroup.GET("/admins/:id", adminHandler.GetAdmin)
			adminGroup.POST("/admins", adminHandler.CreateAdmin)
			adminGroup.PUT("/admins/:id", adminHandler.UpdateAdmin)
			adminGroup.DELETE("/admins/:id", adminHandler.DeleteAdmin)
			adminGroup.POST("/admin/password", adminHandler.ChangePassword)

			adminGroupHandler := admin.NewAdminGroupHandler()
			adminGroup.GET("/admin_groups", adminGroupHandler.List)
			adminGroup.GET("/admin_groups/:id", adminGroupHandler.Get)
			adminGroup.POST("/admin_groups", adminGroupHandler.Create)
			adminGroup.PUT("/admin_groups/:id", adminGroupHandler.Update)
			adminGroup.DELETE("/admin_groups/:id", adminGroupHandler.Delete)
		}
	}

	log.Printf("Admin API server starting on port %s...", cfg.AdminPort)
	r.Run(":" + cfg.AdminPort)
}
