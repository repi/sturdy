package notification

import (
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"
)

type Notification struct {
	ID               string           `db:"id"`
	CodebaseID       codebases.ID     `db:"codebase_id"`
	UserID           users.ID         `db:"user_id"`
	NotificationType NotificationType `db:"type"`
	ReferenceID      string           `db:"reference_id"`
	CreatedAt        time.Time        `db:"created_at"`
	ArchivedAt       *time.Time       `db:"archived_at"`
}

type NotificationType string

const (
	NotificationTypeUndefined       NotificationType = ""
	CommentNotificationType         NotificationType = "comment"
	ReviewNotificationType          NotificationType = "review"
	RequestedReviewNotificationType NotificationType = "requested_review"
	NewSuggestionNotificationType   NotificationType = "new_suggesion"
	GitHubRepositoryImported        NotificationType = "github_repository_imported"
)
