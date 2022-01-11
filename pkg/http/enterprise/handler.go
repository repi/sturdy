package enterprise

import (
	authz "mash/pkg/auth"
	db_change "mash/pkg/change/db"
	workers_ci "mash/pkg/ci/workers"
	db_codebase "mash/pkg/codebase/db"
	service_comments "mash/pkg/comments/service"
	ghappclient "mash/pkg/github/client"
	"mash/pkg/github/config"
	db_github "mash/pkg/github/db"
	routes_v3_ghapp "mash/pkg/github/enterprise/routes"
	workers_github "mash/pkg/github/enterprise/workers"
	service_github "mash/pkg/github/service"
	service_jwt "mash/pkg/jwt/service"
	db_review "mash/pkg/review/db"
	service_statuses "mash/pkg/statuses/service"
	service_sync "mash/pkg/sync/service"
	db_user "mash/pkg/user/db"
	"mash/pkg/view/events"
	activity_sender "mash/pkg/workspace/activity/sender"
	db_workspace "mash/pkg/workspace/db"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs/executor"
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

type DevelopmentAllowExtraCorsOrigin string

type Engine = gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	postHogClient posthog.Client,
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceReader db_workspace.WorkspaceReader,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	gitHubInstallationsRepo db_github.GitHubInstallationRepo,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPRRepo db_github.GitHubPRRepo,
	gitHubAppConfig config.GitHubAppConfig,
	gitHubClientProvider ghappclient.ClientProvider,
	gitHubClonerPublisher *workers_github.ClonerQueue,
	workspaceWriter db_workspace.WorkspaceWriter,
	executorProvider executor.Provider,
	reviewRepo db_review.ReviewRepository,
	activitySender activity_sender.ActivitySender,
	eventSender events.EventSender,
	workspaceService service_workspace.Service,
	statusesService *service_statuses.Service,
	syncService *service_sync.Service,
	jwtService *service_jwt.Service,
	commentsService *service_comments.Service,
	gitHubService *service_github.Service,
	ciBuildQueue *workers_ci.BuildQueue,
	ossEngine *gin.Engine,
) http.Handler {
	publ := ossEngine.Group("")
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubAppConfig, postHogClient, gitHubInstallationsRepo, gitHubRepositoryRepo, codebaseRepo, executorProvider, gitHubClientProvider, gitHubUserRepo, codebaseUserRepo, gitHubClonerPublisher, gitHubPRRepo, workspaceReader, workspaceWriter, workspaceService, syncService, changeRepo, changeCommitRepo, reviewRepo, eventSender, activitySender, statusesService, commentsService, gitHubService, ciBuildQueue))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))
	return ossEngine
}