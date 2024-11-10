package main

import (
	"benchmark/config"
	"benchmark/repository"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, repos, err := prepare()
	if err != nil {
		logrus.Error("Failed to prepare values ", err)
		return
	}

	wgEnd := &sync.WaitGroup{}
	nRequests := new(int64)
	chStart := make(chan struct{})

	for i := 0; i < cfg.ParallelProc; i++ {
		wgEnd.Add(1)
		i := i
		go repoRun(ctx, repos[i], chStart, cfg, nRequests, wgEnd)
	}

	tStart := time.Now()
	close(chStart)
	wgEnd.Wait()

	tEnd := time.Now().Sub(tStart).Milliseconds()
	outputResults(tEnd, nRequests)
}

func prepare() (*config.Config, []*repository.Repo, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, nil, errors.WithMessage(err, "config init failed")
	}

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
		return nil, nil, errors.WithMessage(err, "parse level failed")
	}
	logrus.SetLevel(lvl)

	maxProc := cfg.ParallelProc

	repos := make([]*repository.Repo, maxProc)
	for i := 0; i < maxProc; i++ {
		repos[i], err = repository.NewRepo(cfg.DSN)
		if err != nil {
			return nil, nil, errors.WithMessage(err, "repo init failed")
		}
	}

	return cfg, repos, nil
}

func repoRun(ctx context.Context, repo *repository.Repo, chStart chan struct{}, cfg *config.Config, nRequests *int64, wgEnd *sync.WaitGroup) {
	logrus.Debug("repo run")
	defer wgEnd.Done()

	timer := time.NewTimer(cfg.TestTime)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			logrus.Debug("ctx done ", ctx.Err())
			return
		case <-timer.C:
			logrus.Debug("timer done")
			return
		case <-chStart:
			err := repo.DoRequest(ctx, cfg.TextRequest)
			if err != nil {
				logrus.Error(err)
			}
			atomic.AddInt64(nRequests, 1)
		}
	}
}

func outputResults(tEnd int64, nRequests *int64) {
	logrus.Info("RPS measurement time in milliseconds = ", tEnd)
	logrus.Info("Number of completed requests = ", *nRequests)
	RPS := *nRequests * 1000 / tEnd
	logrus.Info(fmt.Sprintf("RPS = %s op/s", strconv.Itoa(int(RPS))))
}
