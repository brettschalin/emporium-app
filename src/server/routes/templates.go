package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pongo "github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"
)

func setupTemplateRoutes(group string, r *gin.Engine) {

	t := r.Group(group)
	t.GET("/*path", getHTMLTemplate)

}

func getHTMLTemplate(c *gin.Context) {

	path := strings.TrimPrefix(c.Param("path"), "/")
	switch path {
	case "postlist":
		showModButtons, _ := strconv.ParseBool(c.Query("showModButtons"))
		c.HTML(http.StatusOK, "postlist.html", pongo.Context{"showModButtons": showModButtons})
	default:
		c.AbortWithError(http.StatusNotFound,
			fmt.Errorf(`No template at path %q`, path))
	}

}
