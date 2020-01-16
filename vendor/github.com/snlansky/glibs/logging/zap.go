/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package logging

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
)

// NewZapLogger creates a zap logger around a new zap.Core. The core will use
// the provided encoder and sinks and a level enabler that is associated with
// the provided logger name. The logger that is returned will be named the same
// as the logger.
func NewZapLogger(core zapcore.Core, options ...zap.Option) *zap.Logger {
	return zap.New(
		core,
		append([]zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		}, options...)...,
	)
}

// NewGRPCLogger creates a grpc.Logger that delegates to a zap.Logger.
func NewGRPCLogger(l *zap.Logger) *zapgrpc.Logger {
	l = l.WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(3),
	)
	return zapgrpc.NewLogger(l, zapgrpc.WithDebug())
}

// NewGlibsLogger creates a logger that delegates to the zap.SugaredLogger.
func NewGlibsLogger(l *zap.Logger, options ...zap.Option) *GlibsLogger {
	return &GlibsLogger{
		s: l.WithOptions(append(options, zap.AddCallerSkip(1))...).Sugar(),
	}
}

// A GlibsLogger is an adapter around a zap.SugaredLogger that provides
// structured logging capabilities while preserving much of the legacy logging
// behavior.
//
// The most significant difference between the GlibsLogger and the
// zap.SugaredLogger is that methods without a formatting suffix (f or w) build
// the log entry message with fmt.Sprintln instead of fmt.Sprint. Without this
// change, arguments are not separated by spaces.
type GlibsLogger struct{ s *zap.SugaredLogger }

func (f *GlibsLogger) DPanic(args ...interface{})                    { f.s.DPanicf(formatArgs(args)) }
func (f *GlibsLogger) DPanicf(template string, args ...interface{})  { f.s.DPanicf(template, args...) }
func (f *GlibsLogger) DPanicw(msg string, kvPairs ...interface{})    { f.s.DPanicw(msg, kvPairs...) }
func (f *GlibsLogger) Debug(args ...interface{})                     { f.s.Debugf(formatArgs(args)) }
func (f *GlibsLogger) Debugf(template string, args ...interface{})   { f.s.Debugf(template, args...) }
func (f *GlibsLogger) Debugw(msg string, kvPairs ...interface{})     { f.s.Debugw(msg, kvPairs...) }
func (f *GlibsLogger) Error(args ...interface{})                     { f.s.Errorf(formatArgs(args)) }
func (f *GlibsLogger) Errorf(template string, args ...interface{})   { f.s.Errorf(template, args...) }
func (f *GlibsLogger) Errorw(msg string, kvPairs ...interface{})     { f.s.Errorw(msg, kvPairs...) }
func (f *GlibsLogger) Fatal(args ...interface{})                     { f.s.Fatalf(formatArgs(args)) }
func (f *GlibsLogger) Fatalf(template string, args ...interface{})   { f.s.Fatalf(template, args...) }
func (f *GlibsLogger) Fatalw(msg string, kvPairs ...interface{})     { f.s.Fatalw(msg, kvPairs...) }
func (f *GlibsLogger) Info(args ...interface{})                      { f.s.Infof(formatArgs(args)) }
func (f *GlibsLogger) Infof(template string, args ...interface{})    { f.s.Infof(template, args...) }
func (f *GlibsLogger) Infow(msg string, kvPairs ...interface{})      { f.s.Infow(msg, kvPairs...) }
func (f *GlibsLogger) Panic(args ...interface{})                     { f.s.Panicf(formatArgs(args)) }
func (f *GlibsLogger) Panicf(template string, args ...interface{})   { f.s.Panicf(template, args...) }
func (f *GlibsLogger) Panicw(msg string, kvPairs ...interface{})     { f.s.Panicw(msg, kvPairs...) }
func (f *GlibsLogger) Warn(args ...interface{})                      { f.s.Warnf(formatArgs(args)) }
func (f *GlibsLogger) Warnf(template string, args ...interface{})    { f.s.Warnf(template, args...) }
func (f *GlibsLogger) Warnw(msg string, kvPairs ...interface{})      { f.s.Warnw(msg, kvPairs...) }
func (f *GlibsLogger) Warning(args ...interface{})                   { f.s.Warnf(formatArgs(args)) }
func (f *GlibsLogger) Warningf(template string, args ...interface{}) { f.s.Warnf(template, args...) }

// for backwards compatibility
func (f *GlibsLogger) Critical(args ...interface{})                   { f.s.Errorf(formatArgs(args)) }
func (f *GlibsLogger) Criticalf(template string, args ...interface{}) { f.s.Errorf(template, args...) }
func (f *GlibsLogger) Notice(args ...interface{})                     { f.s.Infof(formatArgs(args)) }
func (f *GlibsLogger) Noticef(template string, args ...interface{})   { f.s.Infof(template, args...) }

func (f *GlibsLogger) Named(name string) *GlibsLogger { return &GlibsLogger{s: f.s.Named(name)} }
func (f *GlibsLogger) Sync() error                     { return f.s.Sync() }
func (f *GlibsLogger) Zap() *zap.Logger                { return f.s.Desugar() }

func (f *GlibsLogger) IsEnabledFor(level zapcore.Level) bool {
	return f.s.Desugar().Core().Enabled(level)
}

func (f *GlibsLogger) With(args ...interface{}) *GlibsLogger {
	return &GlibsLogger{s: f.s.With(args...)}
}

func (f *GlibsLogger) WithOptions(opts ...zap.Option) *GlibsLogger {
	l := f.s.Desugar().WithOptions(opts...)
	return &GlibsLogger{s: l.Sugar()}
}

func formatArgs(args []interface{}) string { return strings.TrimSuffix(fmt.Sprintln(args...), "\n") }
