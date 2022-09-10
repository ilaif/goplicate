package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Publisher struct {
	baseBranch string

	repo     *git.Repository
	worktree *git.Worktree
	status   git.Status
}

func NewPublisher(baseBranch string) *Publisher {
	return &Publisher{baseBranch: baseBranch}
}

func (p *Publisher) Init(ctx context.Context, gitPath string) error {
	var err error

	p.repo, err = git.PlainOpen(gitPath)
	if err != nil {
		return errors.Wrap(err, "Failed to open repository")
	}

	p.worktree, err = p.repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "Failed to open worktree")
	}

	p.status, err = p.worktree.Status()
	if err != nil {
		return errors.Wrap(err, "Failed to get worktree status")
	}

	return nil
}

func (p *Publisher) IsClean() bool {
	return p.status.IsClean()
}

func (p *Publisher) Publish(ctx context.Context, filePaths []string) error {
	fmt.Println("Publishing changes...")

	origHeadRef, err := p.repo.Head()
	if err != nil {
		return errors.Wrap(err, "Failed to get HEAD reference")
	}

	if p.baseBranch != "" {
		if outBytes, err := exec. // nolint:gosec
						CommandContext(ctx, "git", "checkout", p.baseBranch).
						CombinedOutput(); err != nil {
			return errors.Wrapf(err, "Failed to checkout base branch '%s': %s", p.baseBranch, string(outBytes))
		}
	}
	defer func() {
		if err := p.worktree.Checkout(&git.CheckoutOptions{Branch: origHeadRef.Name()}); err != nil {
			fmt.Printf("ERROR: Failed to checkout back to original branch '%s'\n", p.baseBranch)
		}
	}()

	if outBytes, err := exec.CommandContext(ctx, "git", "pull").CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Failed to pull branch: %s", string(outBytes))
	}

	headRef, err := p.repo.Head()
	if err != nil {
		return errors.Wrap(err, "Failed to get HEAD reference")
	}
	branchName := "chore/update-goplicate-snippets"
	branchRefName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(branchRefName, headRef.Hash())
	if err := p.repo.Storer.SetReference(ref); err != nil {
		return errors.Wrap(err, "Failed to create new branch")
	}

	if err := p.worktree.Checkout(&git.CheckoutOptions{Branch: ref.Name()}); err != nil {
		return errors.Wrap(err, "Failed to checkout branch")
	}

	filePaths = lo.Uniq(append(filePaths, lo.Keys(p.status)...))
	for _, path := range filePaths {
		if _, err := p.worktree.Add(path); err != nil {
			return errors.Wrap(err, "Failed to add all files to the worktree")
		}
	}

	commitMsg := "chore: update goplicate snippets"
	if _, err := p.worktree.Commit(commitMsg, &git.CommitOptions{}); err != nil {
		return errors.Wrap(err, "Failed to commit changes")
	}

	if outBytes, err := exec.CommandContext(ctx, "git", "push", "-u", "origin", branchName).CombinedOutput(); err != nil {
		return errors.Wrapf(err, "Failed to push changes: %s", string(outBytes))
	}

	changedPathsStr := strings.Join(lo.Map(filePaths, func(path string, _ int) string { return "* " + path }), "\n")
	prBody := fmt.Sprintf("# Update goplicate snippets\n\nUpdated files:\n\n%s", changedPathsStr)
	outBytes, err := exec.
		CommandContext(ctx, "gh", "pr", "create", "--title", commitMsg, "--body", prBody).
		CombinedOutput()
	resp := string(outBytes)
	alreadyExists := strings.Contains(resp, "already exists:")
	if err != nil && !alreadyExists {
		return errors.Wrapf(err, "Failed to create a PR: %s", resp)
	}

	if alreadyExists {
		fmt.Printf("WARNING: PR already exists: %s", resp)
	} else {
		fmt.Printf("Created a PR: %s", resp)
	}

	return nil
}
