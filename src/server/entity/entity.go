package entity

import (
	"time"

	"github.com/google/uuid"
)

// Database tables

type base struct {
	ID      uuid.UUID `json:"id"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// For posts and replies
type Post struct {
	base `json:",inline"`

	CreatedBy  uuid.UUID `json:"created_by"` // must be ID for a User
	Expiration time.Time `json:"expiration"`

	// nullable Post IDs
	RootParent   *uuid.UUID `json:"root_parent,omitempty"`
	DirectParent *uuid.UUID `json:"direct_parent,omitempty"`

	Contents   string `json:"contents"`
	Approved   bool   `json:"approved"`
	InModQueue bool   `json:"in_mod_queue"`
}

func (p *Post) ScanFields() []any {
	return []any{
		&p.ID, &p.Created, &p.Updated, &p.CreatedBy, &p.Expiration, &p.Contents, &p.Approved, &p.InModQueue,
	}
}

type User struct {
	base `json:",inline"`

	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func (u *User) ScanFields() []any {
	return []any{
		&u.ID, &u.Created, &u.Updated, &u.Name, &u.Email, &u.IsAdmin,
	}
}

// Various constants
var (
	PostLifetime = 7 * 24 * time.Hour
)

// Request objects

type CreatePostRequest struct {
	CreatedBy    uuid.UUID `json:"created_by" binding:"required" form:"created_by"`
	RootParent   uuid.UUID `json:"root_parent" form:"root_parent"`
	DirectParent uuid.UUID `json:"direct_parent" form:"direct_parent"`
	Contents     string    `json:"contents" binding:"required,lt=10000" form:"contents"`

	// Currently hardcoded for seven days. TODO: allow more options
	//Expiration time.Time
}

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required,lt=100" form:"name"`
	Email   string `json:"email" binding:"required,email" form:"email"`
	IsAdmin bool   `json:"is_admin" form:"is_admin"`
}

type UpdatePostRequest struct {
	Contents *string `json:"contents" binding:"lt=10000" form:"contents"`
}

type ApprovePostsRequest struct {
	Approved bool        `json:"approved" form:"approved"`
	IDs      []uuid.UUID `json:"ids" binding:"required" form:"ids"`
}
