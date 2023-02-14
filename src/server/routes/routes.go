package routes

import (
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	
	setupPostRoutes("/api/posts", r)

	setupUserRoutes("/api/users", r)

	// Server-side HTML generation. Used on the client via elements with `hx-get=...`
	setupTemplateRoutes("/templates", r)
}

