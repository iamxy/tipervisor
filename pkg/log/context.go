package log

import (
	"context"
	"path"
)

type (
	loggerKey struct{}
	moduleKey struct{}
)

// WithLogger returns a new context with the provided logger. Use in
// combination with WithField(s) for great effect.
func WithLogger(ctx context.Context, entry *Entry) context.Context {
	return context.WithValue(ctx, loggerKey{}, entry)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *Entry {
	logger := ctx.Value(loggerKey{})
	if logger == nil {
		return (*Entry)(L)
	}
	return logger.(*Entry)
}

// WithModule adds the module to the context, appending it with a slash if a
// module already exists. A module is just an roughly correlated defined by the
// call tree for a given context.
func WithModule(ctx context.Context, module string) context.Context {
	parent := GetModulePath(ctx)
	if parent != "" {
		if path.Base(parent) == module {
			return ctx
		}
		module = path.Join(parent, module)
	}

	ctx = WithLogger(ctx, GetLogger(ctx).WithField("module", module))
	return context.WithValue(ctx, moduleKey{}, module)
}

// GetModulePath returns the module path for the provided context. If no module
// is set, an empty string is returned
func GetModulePath(ctx context.Context) string {
	module := ctx.Value(moduleKey{})
	if module == nil {
		return ""
	}
	return module.(string)
}
