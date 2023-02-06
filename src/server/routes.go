package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"marketplace-server/db"
	"marketplace-server/entity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupRoutes(r *gin.Engine) {
	p := r.Group("/api/posts")
	u := r.Group("/api/users")

	p.GET("/", GetPosts)
	p.GET("/:id", GetPostByID)
	p.POST("/", CreatePost)
	p.PUT("/:id", UpdatePostByID)
	p.POST("/approve", ApprovePosts)
	p.DELETE("/:id", DeletePostByID)

	u.GET("/:id", GetUserByID)
	u.POST("/", CreateUser)
	u.DELETE("/:id", DeleteUserByID)
}

func GetPosts(c *gin.Context) {
	ctx := context.Background()

	var (
		err        error
		inModQueue bool
	)

	inModQueue, err = strconv.ParseBool(c.Query("in_mod_queue"))
	if err != nil {
		inModQueue = false
	}

	posts, err := db.GetPosts(ctx, inModQueue)

	if err != nil {
		if err == db.ErrNotFound {
			c.String(http.StatusNotFound, "")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, posts)
	}
}

func GetPostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, "ID must be a valid UUID")
		return
	}

	ctx := context.Background()

	post, err := db.GetPost(ctx, id)

	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {

	var (
		p  entity.Post
		cr entity.CreatePostRequest
	)
	if err := c.BindJSON(&cr); err != nil {
		c.JSON(http.StatusBadRequest, "could not bind request body")
		return
	}

	p.CreatedBy = cr.CreatedBy

	if r := cr.RootParent; r != uuid.Nil {
		p.RootParent = &r
	}
	if d := cr.DirectParent; d != uuid.Nil {
		p.DirectParent = &d
	}

	p.Contents = cr.Contents
	p.Expiration = time.Now().Add(entity.PostLifetime)

	ctx := context.Background()

	id, err := db.CreatePost(ctx, &p)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func UpdatePostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, "ID must be a valid UUID")
		return
	}

	var (
		ur entity.UpdatePostRequest
	)

	if err := c.BindJSON(&ur); err != nil {
		c.JSON(http.StatusBadRequest, "could not bind request body")
		return
	}

	ctx := context.Background()

	_, err = db.UpdatePost(ctx, id, &ur)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func ApprovePosts(c *gin.Context) {
	var (
		err error
		ar  entity.ApprovePostsRequest
	)

	if err = c.BindJSON(&ar); err != nil {
		c.JSON(http.StatusBadRequest, "could not bind request body")
		return
	}

	ctx := context.Background()
	err = db.ApprovePosts(ctx, ar.Approved, ar.IDs)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

func DeletePostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, "ID must be a valid UUID")
		return
	}

	ctx := context.Background()

	err = db.DeletePost(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unexpected error: "+err.Error())
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func CreateUser(c *gin.Context) {
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

func GetUserByID(c *gin.Context) {
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

func DeleteUserByID(c *gin.Context) {
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
