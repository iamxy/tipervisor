package log

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, GetLogger(ctx), (*Entry)(L))

	ctx = WithLogger(ctx, GetLogger(ctx).WithField("test", "one"))
	assert.Equal(t, (*logrus.Entry)(GetLogger(ctx)).Data["test"], "one")
}

func TestModuleContext(t *testing.T) {
	ctx := context.Background()
	assert.Equal(t, GetModulePath(ctx), "")

	ctx = WithModule(ctx, "a")
	assert.Equal(t, GetModulePath(ctx), "a")
	logger := (*logrus.Entry)(GetLogger(ctx))
	assert.Equal(t, logger.Data["module"], "a")

	parent, ctx := ctx, WithModule(ctx, "a")
	assert.Equal(t, ctx, parent)
	assert.Equal(t, GetModulePath(ctx), "a")
	assert.Equal(t, (*logrus.Entry)(GetLogger(ctx)).Data["module"], "a")

	ctx = WithModule(ctx, "b")
	assert.Equal(t, GetModulePath(ctx), "a/b")
	assert.Equal(t, (*logrus.Entry)(GetLogger(ctx)).Data["module"], "a/b")

	ctx = WithModule(ctx, "c")
	assert.Equal(t, GetModulePath(ctx), "a/b/c")
	assert.Equal(t, (*logrus.Entry)(GetLogger(ctx)).Data["module"], "a/b/c")
}
