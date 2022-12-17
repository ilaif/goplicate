package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/caarlos0/log"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/ilaif/goplicate/pkg/utils"
)

// Publisher publishes changes to git, including opening PRs
type Publisher struct {
	baseBranch string
	dir        string
	cmdRunner  *utils.CommandRunner

	repo   *git.Repository
	status git.Status
}

func NewPublisher(baseBranch string, dir string) *Publisher {
	cmdRunner := utils.NewCommandRunner(dir)

	return &Publisher{baseBranch: baseBranch, dir: dir, cmdRunner: cmdRunner}
}

func (p *Publisher) Init(ctx context.Context) error {
	var err error

	log.Debugf("Opening repository '%s'", p.dir)
	p.repo, err = git.PlainOpen(p.dir)
	if err != nil {
		return errors.Wrap(err, "Failed to open repository")
	}

	log.Debug("Opening worktree")
	worktree, err := p.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed to open worktree")
	}

	log.Debug("Getting worktree status")
	p.status, err = worktree.Status()
	if err != nil {
		return errors.Wrap(err, "Failed to get worktree status")
	}

	return nil
}

func (p *Publisher) StashChanges(ctx context.Context) (func() error, error) {
	log.Debug("Stashing working directory changes")
	if output, err := p.cmdRunner.Run(ctx, "git", "stash"); err != nil {
		return nil, errors.Wrapf(err, "Failed to stash local changes: %s", output)
	}

	return func() error {
		log.Debug("Cleanup: Un-stashing working directory changes")
		if output, err := p.cmdRunner.Run(ctx, "git", "stash", "pop"); err != nil {
			return errors.Wrapf(err, "Cleanup: Failed to restore local changes: %s", output)
		}

		return nil
	}, nil
}

func (p *Publisher) IsClean() bool {
	return p.status.IsClean()
}

func (p *Publisher) Publish(ctx context.Context, filePaths []string, confirm bool) error {
	log.Info("Publishing changes...")

	log.Debug("Fetching current branch name")
	origBranchName, err := p.cmdRunner.Run(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return errors.Wrapf(err, "Failed to fetch current branch name: %s", origBranchName)
	}
	origBranchName = strings.Trim(origBranchName, "\n")

	if p.baseBranch != "" {
		log.Debugf("Checking out base branch '%s'", p.baseBranch)
		if output, err := p.cmdRunner.Run(ctx, "git", "checkout", p.baseBranch); err != nil {
			return errors.Wrapf(err, "Failed to checkout base branch '%s': %s", p.baseBranch, output)
		}
	}
	defer func() {
		log.Debugf("Cleanup: Checking out original branch '%s'", origBranchName)
		if output, err := p.cmdRunner.Run(ctx, "git", "checkout", string(origBranchName)); err != nil {
			log.WithError(err).Errorf("Cleanup: Failed to checkout back to original branch '%s': %s", p.baseBranch, output)
		}
	}()

	log.Debugf("Pulling from remote")
	if output, err := p.cmdRunner.Run(ctx, "git", "pull"); err != nil {
		return errors.Wrapf(err, "Failed to pull branch: %s", output)
	}

	log.Debug("Fetching HEAD reference")
	branchName := "chore/update-goplicate-snippets"

	log.Debugf("Deleting existing branch '%s' if exists", branchName)
	if output, err := p.cmdRunner.Run(ctx, "git", "branch", "-D", branchName); err != nil {
		log.WithError(err).Debugf("Failed to delete existing branch '%s': %s", branchName, output)
	}

	remoteOriginURL, err := p.cmdRunner.Run(ctx, "git", "config", "--get", "remote.origin.url")
	remoteOriginURL = strings.Trim(remoteOriginURL, "\n")
	if err != nil {
		return errors.Wrapf(err, "Failed to get remote origin url: %s", remoteOriginURL)
	}

	output, err := p.cmdRunner.Run(ctx, "git", "ls-remote", "--heads", remoteOriginURL, branchName)
	if err != nil {
		return errors.Wrapf(err, "Failed to list remote branches: %s", output)
	}
	if strings.Contains(output, fmt.Sprintf("refs/heads/%s", branchName)) {
		// Remote branch exists
		if !confirm && !utils.AskUserYesNoQuestion(
			fmt.Sprintf("Found branch '%s' in origin. Do you want to delete it?", branchName),
		) {
			return errors.New("User aborted")
		}

		output, err := p.cmdRunner.Run(ctx, "git", "push", "-d", "origin", branchName)
		if err != nil {
			return errors.Wrapf(err, "Failed to delete existing remote branch '%s': %s", branchName, output)
		}
	}

	log.Debugf("Checking out new branch '%s'", branchName)
	if output, err := p.cmdRunner.Run(ctx, "git", "checkout", "-b", branchName); err != nil {
		return errors.Wrapf(err, "Failed to checkout new branch '%s': %s", branchName, output)
	}

	filePaths = lo.Uniq(append(filePaths, lo.Keys(p.status)...))
	for _, path := range filePaths {
		log.Debugf("Adding file '%s' to the worktree", path)
		if output, err := p.cmdRunner.Run(ctx, "git", "add", path); err != nil {
			return errors.Wrapf(err, "Failed to add files to the worktree: %s", output)
		}
	}

	log.Debug("Committing changes")
	commitMsg := "chore: update goplicate snippets"
	if output, err := p.cmdRunner.Run(ctx, "git", "commit", "-m", commitMsg); err != nil {
		return errors.Wrapf(err, "Failed to commit changes: %s", output)
	}

	log.Debug("Pushing changes")
	if output, err := p.cmdRunner.Run(ctx, "git", "push", "-u", "origin", branchName); err != nil {
		return errors.Wrapf(err, "Failed to push changes: %s", output)
	}

	changedPathsStr := strings.Join(lo.Map(filePaths, func(path string, _ int) string { return "* " + path }), "\n")
	prBody := fmt.Sprintf("# Update goplicate snippets\n\nUpdated files:\n\n%s", changedPathsStr)

	log.Debug("Creating pull request")
	resp, err := p.cmdRunner.Run(ctx, "gh", "pr", "create", "--title", commitMsg, "--body", prBody, "--head", branchName)
	resp = strings.TrimSuffix(resp, "\n")
	alreadyExists := strings.Contains(resp, "already exists:")
	if err != nil && !alreadyExists {
		return errors.Wrapf(err, "Failed to create a PR: %s", resp)
	}

	if alreadyExists {
		log.Warnf("PR already exists: %s", resp)
	} else {
		log.Infof("Created PR: %s", resp)
	}

	return nil
}
