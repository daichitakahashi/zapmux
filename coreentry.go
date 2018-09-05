package zapmux

import "go.uber.org/zap/zapcore"

// CoreEntry is
type CoreEntry struct {
	Core  zapcore.Core
	Entry zapcore.Entry
}

// With is wrapper
func (ce *CoreEntry) With(fields ...zapcore.Field) {
	if ce == nil || len(fields) == 0 {
		return
	}
	ce.Core = ce.Core.With(fields)
	return
}
