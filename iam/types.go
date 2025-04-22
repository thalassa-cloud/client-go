package iam

import (
	"time"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

type Team struct {
	Identity    string            `json:"identity"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   *time.Time        `json:"updatedAt"`
}

type CreateTeam struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type UpdateTeam struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

type TeamMember struct {
	Identity  string       `json:"identity"`
	Role      string       `json:"role"`
	User      base.AppUser `json:"user"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt *time.Time   `json:"updatedAt"`
}
