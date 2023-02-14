package routes

import (
	"context"
	"net/http"

	"github.com/brettschalin/emporium-app/db"
	"github.com/brettschalin/emporium-app/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



func setupUserRoutes(group string, r *gin.Engine) {

	u := r.Group(group)

	u.GET("/:id", getUserByID)
	u.POST("/", createUser)
	u.DELETE("/:id", deleteUserByID)

}



func createUser(c *gin.Context) {
	var (
		cu entity.CreateUserRequest
	)
	if err := c.Bind(&cu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "could not bind request body",
			"error":   err,
		})
		return
	}
	ctx := context.Background()

	id, err := db.CreateUser(ctx, &entity.User{
		Name:    cu.Name,
		Email:   cu.Email,
		IsAdmin: cu.IsAdmin,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unexpected error: "+err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func getUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, "ID must be a valid UUID")
		return
	}

	ctx := context.Background()

	post, err := db.GetUser(ctx, id)

	if err != nil {
		c.JSON(http.StatusNotFound, "could not find user: "+err.Error())
	}

	c.JSON(http.StatusOK, post)
}

func deleteUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, "ID must be a valid UUID")
		return
	}

	ctx := context.Background()

	err = db.DeleteUser(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unexpected error: "+err.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
