package ci

import (
	"context"
	"time"

	"github.com/shopware/shopware-cli/logging"
)

// DefaultCi implements CiHelper for non-CI environments.
type DefaultCi struct{}

type DefaultCiSection struct {
	name  string
	ctx   context.Context
	start time.Time
}

// Section starts a new log section.
func (d *DefaultCi) Section(ctx context.Context, name string) Section {
	logging.FromContext(ctx).Infof("Starting %s", name)
	return DefaultCiSection{
		name:  name,
		ctx:   ctx,
		start: time.Now(),
	}
}

// SectionEnd ends the current log section.
func (d DefaultCiSection) End() {
	logging.FromContext(d.ctx).Infof("%s ended after %s", d.name, time.Since(d.start))
}
