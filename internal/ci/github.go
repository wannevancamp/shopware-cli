package ci

import (
	"context"
	"fmt"
	"time"

	"github.com/shopware/shopware-cli/logging"
)

// GithubActions implements CiHelper for GitHub Actions.
type GithubActions struct{}

type GithubActionsSection struct {
	name  string
	ctx   context.Context
	start time.Time
}

// SectionStart starts a new log section.
func (g *GithubActions) Section(ctx context.Context, name string) Section {
	fmt.Printf("::group::%s\n", name)
	return GithubActionsSection{
		ctx:   ctx,
		name:  name,
		start: time.Now(),
	}
}

func (s GithubActionsSection) End() {
	duration := time.Since(s.start)
	fmt.Printf("::endgroup::\n")
	logging.FromContext(s.ctx).Infof("Section '%s' ended after %s", s.name, duration)
}
