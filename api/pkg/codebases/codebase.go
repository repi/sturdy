package codebases

import (
	"time"

	"getsturdy.com/api/pkg/author"
	"getsturdy.com/api/pkg/users"

	"github.com/gosimple/slug"
)

type ShortCodebaseID string

func (s ShortCodebaseID) String() string {
	return string(s)
}

type ID string

func (id ID) String() string {
	return string(id)
}

type Codebase struct {
	ID              ID              `db:"id" json:"id"`
	ShortCodebaseID ShortCodebaseID `db:"short_id" json:"short_id"` // Used in Web slugs
	Name            string          `db:"name" json:"name"`
	Description     string          `db:"description" json:"description"`
	Emoji           string          `db:"emoji" json:"emoji"`
	InviteCode      *string         `db:"invite_code" json:"invite_code"`
	CreatedAt       *time.Time      `db:"created_at" json:"created_at"`
	ArchivedAt      *time.Time      `db:"archived_at" json:"archived_at"`
	OrganizationID  *string         `db:"organization_id"`

	IsReady  bool `json:"is_ready" db:"is_ready"`
	IsPublic bool `json:"is_public" db:"is_public"`

	// Use through ChangeService.HeadChange()
	CalculatedHeadChangeID bool    `json:"-" db:"calculated_head_change_id"`
	CachedHeadChangeID     *string `json:"-" db:"cached_head_change_id"`
}

type CodebaseUser struct {
	ID         string     `db:"id"`
	UserID     users.ID   `db:"user_id"`
	CodebaseID ID         `db:"codebase_id"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`
}

type CodebaseWithMetadata struct {
	Codebase
	Members []author.Author `json:"members"`
}

func (c Codebase) Slug() string {
	return slug.Make(c.Name)
}

func (c Codebase) GenerateSlug() string {
	// TODO: Remove the "Slug" method on the frontend, and generate all slugs on the backend
	return slug.Make(c.Name) + "-" + string(c.ShortCodebaseID)
}
