package selfhosted

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/api"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"
	workers_license "getsturdy.com/api/pkg/installations/enterprise/selfhosted/worker"
	worker_installation_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/selfhosted/worker"

	"golang.org/x/sync/errgroup"
)

type API struct {
	ossAPI *api.API

	githubClonerQueue            *workers_github.ClonerQueue
	githubImporterQueue          workers_github.ImporterQueue
	licenseWorker                *workers_license.Worker
	installationStatisticsWorker *worker_installation_statistics.Worker
	githubWebhooksQueue          *workers_github.WebhooksQueue
}

func ProvideAPI(
	ossAPI *api.API,

	githubClonerQueue *workers_github.ClonerQueue,
	githubImporterQueue workers_github.ImporterQueue,
	licenseWorker *workers_license.Worker,
	installationStatisticsWorker *worker_installation_statistics.Worker,
	githubWebhooksQueue *workers_github.WebhooksQueue,
) *API {
	return &API{
		ossAPI:                       ossAPI,
		githubClonerQueue:            githubClonerQueue,
		githubImporterQueue:          githubImporterQueue,
		licenseWorker:                licenseWorker,
		installationStatisticsWorker: installationStatisticsWorker,
		githubWebhooksQueue:          githubWebhooksQueue,
	}
}

func (a *API) Start(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return a.ossAPI.Start(ctx)
	})

	wg.Go(func() error {
		if err := a.githubClonerQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github cloner queue: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.githubImporterQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github importer queue: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.licenseWorker.Start(ctx); err != nil {
			return fmt.Errorf("failed to start license worker: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.installationStatisticsWorker.Start(ctx); err != nil {
			return fmt.Errorf("failed to start installation statistics worker: %w", err)
		}
		return nil
	})

	wg.Go(func() error {
		if err := a.githubWebhooksQueue.Start(ctx); err != nil {
			return fmt.Errorf("failed to start github webhooks queue: %w", err)
		}
		return nil
	})

	return wg.Wait()
}
