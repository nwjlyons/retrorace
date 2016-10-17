package handlers

import (
	"fmt"
	"github.com/nwjlyons/retrorace/retroengine"
	"gopkg.in/macaron.v1"
	"runtime"
	"syscall"
)

func Stats(ctx *macaron.Context) (string, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return ctx.JSONString(struct {
		GOVERSION    string
		GOARCH       string
		GOOS         string
		NumCPU       int
		NumGoroutine int
		PID          int
		Memory       string
		NumGames     int
	}{
		GOVERSION:    runtime.Version(),
		GOARCH:       runtime.GOARCH,
		GOOS:         runtime.GOOS,
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		Memory:       fmt.Sprintf("%.3fMB", float64(mem.Alloc)/1000000),
		PID:          syscall.Getpid(),
		NumGames:     len(retroengine.GameStore.Games),
	})
}
