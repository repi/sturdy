package vcs

import (
	"fmt"

	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	git "github.com/libgit2/git2go/v33"
	"go.uber.org/zap"
)

func FetchTrackedToSturdytrunk(accessToken, ref string) func(vcs.RepoGitWriter) error {
	return func(repo vcs.RepoGitWriter) error {
		refspec := fmt.Sprintf("+%s:refs/heads/sturdytrunk", ref)
		if err := repo.FetchNamedRemoteWithCreds("origin", newCredentialsCallback(accessToken), []string{refspec}); err != nil {
			return fmt.Errorf("failed to perform remote fetch: %w", err)
		}

		// Make sure that sturdytrunk is the HEAD branch
		// This is the case for repositories that where empty the first time they where cloned to Sturdy
		if err := repo.SetDefaultBranch("sturdytrunk"); err != nil {
			return fmt.Errorf("could not set default branch: %w", err)
		}
		return nil
	}
}

func FetchBranchWithRefspec(accessToken, refspec string) func(vcs.RepoGitWriter) error {
	return func(repo vcs.RepoGitWriter) error {
		if err := repo.FetchNamedRemoteWithCreds("origin", newCredentialsCallback(accessToken), []string{refspec}); err != nil {
			return fmt.Errorf("failed to perform remote fetch: %w", err)
		}
		return nil
	}
}

func PushTrackedToGitHub(logger *zap.Logger, repo vcs.RepoGitWriter, accessToken, trackedBranchName string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/sturdytrunk:refs/heads/%s", trackedBranchName)
	userError, err = repo.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
	if err != nil {
		return userError, fmt.Errorf("failed to push %s: %w", refspec, err)
	}
	return "", nil
}

func PushBranchToGithubWithForce(logger *zap.Logger, executorProvider executor.Provider, codebaseID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("+refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().GitWrite(func(r vcs.RepoGitWriter) error {
		userError, err = r.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
		if err != nil {
			return fmt.Errorf("failed to push %s: %w", refspec, err)
		}
		return nil
	}).ExecTrunk(codebaseID, "PushBranchToGithubWithForce")
	if err != nil {
		return userError, err
	}
	return userError, nil
}

func PushBranchToGithubSafely(logger *zap.Logger, executorProvider executor.Provider, codebaseID, sturdyBranchName, remoteBranchName, accessToken string) (userError string, err error) {
	refspec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", sturdyBranchName, remoteBranchName)

	err = executorProvider.New().GitWrite(func(r vcs.RepoGitWriter) error {
		userError, err = r.PushNamedRemoteWithRefspec(logger, "origin", newCredentialsCallback(accessToken), []string{refspec})
		if err != nil {
			return fmt.Errorf("failed to push %s: %w", refspec, err)
		}
		return nil
	}).ExecTrunk(codebaseID, "PushBranchToGithubSafely")
	if err != nil {
		return userError, err
	}
	return userError, nil
}

func HaveTrackedBranch(executorProvider executor.Provider, codebaseID, remoteBranchName string) error {
	err := executorProvider.New().GitRead(func(r vcs.RepoGitReader) error {
		_, err := r.RemoteBranchCommit("origin", remoteBranchName)
		if err != nil {
			return fmt.Errorf("could not get remote branch: %w", err)
		}
		return nil
	}).ExecTrunk(codebaseID, "haveTrackedBranch")
	if err != nil {
		return err
	}
	return nil
}

func newCredentialsCallback(token string) git.CredentialsCallback {
	return func(url string, username string, allowedTypes git.CredentialType) (*git.Credential, error) {
		cred, _ := git.NewCredentialUserpassPlaintext("x-access-token", token)
		return cred, nil
	}
}
