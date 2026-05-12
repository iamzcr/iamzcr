package router

import (
	"iamzcr/config"
	adminH "iamzcr/handlers/admin"
	"iamzcr/middleware"
	adminS "iamzcr/services/admin"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(r *gin.Engine, cfg *config.Config) {
	articleSvc := adminS.NewArticleService()
	adminSvc := adminS.NewAdminService()
	categorySvc := adminS.NewCategoryService()
	directorySvc := adminS.NewDirectoryService()
	tagsSvc := adminS.NewTagsService()
	menuSvc := adminS.NewMenuService()
	commentSvc := adminS.NewCommentService()
	websiteSvc := adminS.NewWebsiteService()
	attachSvc := adminS.NewAttachService()
	langSvc := adminS.NewLangService()
	logSvc := adminS.NewLogService()
	messageSvc := adminS.NewMessageService()
	permitSvc := adminS.NewPermitService()
	readSvc := adminS.NewReadService()
	adminGroupSvc := adminS.NewAdminGroupService()
	articleMediaSvc := adminS.NewArticleMediaService()
	wechatSvc := adminS.NewWeChatService()
	platformSvc := adminS.NewPlatformService()
	attachMediaSvc := adminS.NewAttachMediaService()

	api := r.Group("/api")
	{
		adminHandler := adminH.NewAdminHandler(articleSvc, adminSvc, wechatSvc)
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

			articleMediaHandler := adminH.NewArticleMediaHandler(articleMediaSvc, articleSvc, wechatSvc)
			adminGroup.GET("/articles/:id/media", articleMediaHandler.ListMedia)
			adminGroup.POST("/articles/:id/media/publish", articleMediaHandler.PublishToMedia)

			categoryHandler := adminH.NewCategoryHandler(categorySvc)
			adminGroup.GET("/categories", categoryHandler.List)
			adminGroup.GET("/categories/:id", categoryHandler.Get)
			adminGroup.POST("/categories", categoryHandler.Create)
			adminGroup.PUT("/categories/:id", categoryHandler.Update)
			adminGroup.DELETE("/categories/:id", categoryHandler.Delete)

			commentHandler := adminH.NewCommentHandler(commentSvc)
			adminGroup.GET("/comments", commentHandler.List)
			adminGroup.GET("/comments/:id", commentHandler.Get)
			adminGroup.POST("/comments", commentHandler.Create)
			adminGroup.PUT("/comments/:id", commentHandler.Update)
			adminGroup.DELETE("/comments/:id", commentHandler.Delete)

			menuHandler := adminH.NewMenuHandler(menuSvc)
			adminGroup.GET("/menus", menuHandler.List)
			adminGroup.GET("/menus/:id", menuHandler.Get)
			adminGroup.POST("/menus", menuHandler.Create)
			adminGroup.PUT("/menus/:id", menuHandler.Update)
			adminGroup.DELETE("/menus/:id", menuHandler.Delete)

			tagsHandler := adminH.NewTagsHandler(tagsSvc)
			adminGroup.GET("/tags", tagsHandler.List)
			adminGroup.GET("/tags/:id", tagsHandler.Get)
			adminGroup.POST("/tags", tagsHandler.Create)
			adminGroup.PUT("/tags/:id", tagsHandler.Update)
			adminGroup.DELETE("/tags/:id", tagsHandler.Delete)

			directoryHandler := adminH.NewDirectoryHandler(directorySvc)
			adminGroup.GET("/directories", directoryHandler.List)
			adminGroup.GET("/directories/:id", directoryHandler.Get)
			adminGroup.POST("/directories", directoryHandler.Create)
			adminGroup.PUT("/directories/:id", directoryHandler.Update)
			adminGroup.DELETE("/directories/:id", directoryHandler.Delete)

			websiteHandler := adminH.NewWebsiteHandler(websiteSvc)
			adminGroup.GET("/website", websiteHandler.Get)
			adminGroup.GET("/website/list", websiteHandler.List)
			adminGroup.PUT("/website", websiteHandler.Update)
			adminGroup.DELETE("/website/:id", websiteHandler.Delete)

			attachHandler := adminH.NewAttachHandler(attachSvc)
			adminGroup.GET("/attaches", attachHandler.List)
			adminGroup.GET("/attaches/:id", attachHandler.Get)
			adminGroup.POST("/attaches", attachHandler.Create)
			adminGroup.PUT("/attaches/:id", attachHandler.Update)
			adminGroup.DELETE("/attaches/:id", attachHandler.Delete)
			adminGroup.POST("/upload", attachHandler.Upload)

			langHandler := adminH.NewLangHandler(langSvc)
			adminGroup.GET("/langs", langHandler.List)
			adminGroup.GET("/langs/:id", langHandler.Get)
			adminGroup.POST("/langs", langHandler.Create)
			adminGroup.PUT("/langs/:id", langHandler.Update)
			adminGroup.DELETE("/langs/:id", langHandler.Delete)

			logHandler := adminH.NewLogHandler(logSvc)
			adminGroup.GET("/logs", logHandler.List)
			adminGroup.GET("/logs/:id", logHandler.Get)
			adminGroup.POST("/logs", logHandler.Create)
			adminGroup.DELETE("/logs/:id", logHandler.Delete)

			messageHandler := adminH.NewMessageHandler(messageSvc)
			adminGroup.GET("/messages", messageHandler.List)
			adminGroup.GET("/messages/:id", messageHandler.Get)
			adminGroup.POST("/messages", messageHandler.Create)
			adminGroup.PUT("/messages/:id", messageHandler.Update)
			adminGroup.DELETE("/messages/:id", messageHandler.Delete)

			permitHandler := adminH.NewPermitHandler(permitSvc)
			adminGroup.GET("/permits", permitHandler.List)
			adminGroup.GET("/permits/:id", permitHandler.Get)
			adminGroup.POST("/permits", permitHandler.Create)
			adminGroup.PUT("/permits/:id", permitHandler.Update)
			adminGroup.DELETE("/permits/:id", permitHandler.Delete)

			readHandler := adminH.NewReadHandler(readSvc)
			adminGroup.GET("/reads", readHandler.List)
			adminGroup.GET("/reads/:id", readHandler.Get)
			adminGroup.POST("/reads", readHandler.Create)
			adminGroup.DELETE("/reads/:id", readHandler.Delete)

			adminGroup.GET("/admins", adminHandler.ListAdmins)
			adminGroup.GET("/admins/:id", adminHandler.GetAdmin)
			adminGroup.POST("/admins", adminHandler.CreateAdmin)
			adminGroup.PUT("/admins/:id", adminHandler.UpdateAdmin)
			adminGroup.DELETE("/admins/:id", adminHandler.DeleteAdmin)
			adminGroup.POST("/admins/:id/password", adminHandler.ChangeAdminPassword)
			adminGroup.POST("/admin/password", adminHandler.ChangePassword)

			adminGroupHandler := adminH.NewAdminGroupHandler(adminGroupSvc)
			adminGroup.GET("/admin_groups", adminGroupHandler.List)
			adminGroup.GET("/admin_groups/:id", adminGroupHandler.Get)
			adminGroup.POST("/admin_groups", adminGroupHandler.Create)
			adminGroup.PUT("/admin_groups/:id", adminGroupHandler.Update)
			adminGroup.DELETE("/admin_groups/:id", adminGroupHandler.Delete)

			platformHandler := adminH.NewPlatformHandler(platformSvc)
			adminGroup.GET("/platforms", platformHandler.List)
			adminGroup.GET("/platforms/:id", platformHandler.Get)
			adminGroup.POST("/platforms", platformHandler.Create)
			adminGroup.PUT("/platforms/:id", platformHandler.Update)
			adminGroup.DELETE("/platforms/:id", platformHandler.Delete)

			attachMediaHandler := adminH.NewAttachMediaHandler(attachMediaSvc, attachSvc, wechatSvc)
			adminGroup.GET("/attach_media", attachMediaHandler.List)
			adminGroup.GET("/attach_media/:id", attachMediaHandler.Get)
			adminGroup.POST("/attaches/:id/sync_wechat", attachMediaHandler.SyncToWechat)
			adminGroup.DELETE("/attach_media/:id", attachMediaHandler.Delete)
		}
	}

	r.Static("/asset", cfg.AssetDir())
}
