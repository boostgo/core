// Package appx provide the whole app lifetime control.
// Features:
// - Global context. To see if app lifetime is end
// - Global context cancel func. You can provide it to other packages & it can be called to graceful shutdown the app
// - Teardown functions. They call at the end of app lifetime
// - Signal catch and calling global context shutdown
//
//nolint:govet
package appx

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"sync"
	"time"

	"github.com/boostgo/core/errorx"
)

var (
	_appOnce sync.Once
	_app     *app
)

// app controls the application lifetime.
//
// Contain global app context, context's cancel function and teardown functions
type app struct {
	ctx         context.Context
	cancel      context.CancelFunc
	tears       []func() error
	gracefulLog func()
}

// Cancel global app context
func (app *app) Cancel() {
	app.cancel()
}

func getApp() *app {
	_appOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		_app = &app{
			ctx:    ctx,
			cancel: cancel,
			tears:  make([]func() error, 0),
		}
	})

	return _app
}

// Context returns global context
func Context() context.Context {
	return getApp().ctx
}

// Cancel call global context cancel function
func Cancel() {
	getApp().cancel()
}

// GracefulLog set function which calls when Cancel called
func GracefulLog(gracefulLog func()) {
	getApp().gracefulLog = gracefulLog
}

// Tear add teardown function which calls after
func Tear(tear func() error) {
	l := getApp()
	l.tears = append(l.tears, tear)
}

// Wait hold current goroutine till global context cancel.
//
// If provide wait time it will wait provided time after calling global context cancel
func Wait(waitTime ...time.Duration) {
	l := getApp()

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		Cancel()
	}()

	<-l.ctx.Done()

	if l.gracefulLog != nil {
		errorx.TryMust(func() error {
			l.gracefulLog()
			return nil
		})
	}

	tears := make([]func() error, len(l.tears))
	copy(tears, l.tears)
	slices.Reverse(tears)

	for _, tear := range tears {
		errorx.TryMust(tear)
	}

	if len(waitTime) > 0 && waitTime[0] > 0 {
		time.Sleep(waitTime[0])
	}
}
