package zapmux

import (
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

// DuxCore is
type DuxCore struct {
	maincore zapcore.Core
	subcore  zapcore.Core
	// interceptor func(mainCore, subCore *zapcore.Core, mainEnt, subEnt *zapcore.Entry)
	interceptor func(main, sub *CoreEntry)
	subEnt      zapcore.Entry
}

// NewDuxCore is
func NewDuxCore(main, sub zapcore.Core) *DuxCore {
	if main == nil {
		main = zapcore.NewNopCore()
	}
	if sub == nil {
		sub = zapcore.NewNopCore()
	}

	return &DuxCore{
		maincore: main,
		subcore:  sub,
	}
}

// WithInterceptor ???
func (dc *DuxCore) WithInterceptor(interceptor func(main, sub *CoreEntry)) *DuxCore {
	newCore := dc.clone()

	if dc.interceptor != nil {
		newCore.interceptor = func(main, sub *CoreEntry) {
			newCore.interceptor(main, sub)
			interceptor(main, sub)
		}
	} else {
		newCore.interceptor = interceptor
	}
	return newCore
}

// Enabled is
func (dc *DuxCore) Enabled(lvl zapcore.Level) bool {
	return dc.maincore.Enabled(lvl) || dc.subcore.Enabled(lvl)
}

// With is
func (dc *DuxCore) With(fields []zapcore.Field) zapcore.Core {
	if len(fields) == 0 {
		return dc
	}
	newCore := dc.clone()
	newCore.maincore = newCore.maincore.With(fields)
	newCore.subcore = newCore.subcore.With(fields)
	return newCore
}

// Check is
func (dc *DuxCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	// if dc.maincore.Enabled(e.Level) && dc.subcore.Enabled(e.Level) && dc.interceptor != nil {
	if dc.interceptor != nil {

		newCore := dc.clone()

		main := CoreEntry{
			Core:  dc.maincore,
			Entry: e,
		}
		sub := CoreEntry{
			Core:  dc.subcore,
			Entry: e,
		}

		newCore.interceptor(&main, &sub)

		newCore.maincore = main.Core
		newCore.subcore = sub.Core
		newCore.subEnt = sub.Entry

		if newCore.maincore != nil && newCore.subcore == nil {
			return ce.AddCore(main.Entry, newCore.maincore)
		} else if newCore.maincore == nil && newCore.subcore != nil {
			return ce.AddCore(newCore.subEnt, newCore.subcore)
		}
		return ce.AddCore(main.Entry, newCore)
	}

	ce = dc.maincore.Check(e, ce)
	return dc.subcore.Check(e, ce)
}

// Write is
func (dc *DuxCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	defer dc.subcore.Check(dc.subEnt, nil).Write(fields...)
	return dc.maincore.Write(e, fields)
}

// Sync is
func (dc *DuxCore) Sync() error {
	var err error
	err = multierr.Append(err, dc.maincore.Sync())
	err = multierr.Append(err, dc.subcore.Sync())
	return err
}

func (dc *DuxCore) clone() *DuxCore {
	copy := *dc
	return &copy
}
