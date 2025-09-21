package bot

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"kick-clipper/internal/client"
	"kick-clipper/internal/config"
	"kick-clipper/internal/models"
	"kick-clipper/internal/proxy"
)

// Manager handles bot operations and coordination
type Manager struct {
	config   *config.Config
	stats    *models.GlobalStats
	client   *client.KickClient
	proxyMgr *proxy.Manager
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewManager creates a new bot manager
func NewManager(cfg *config.Config, stats *models.GlobalStats, proxyMgr *proxy.Manager) *Manager {
	timeout := time.Duration(cfg.Timeout) * time.Second
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		config:   cfg,
		stats:    stats,
		client:   client.NewKickClient(timeout),
		proxyMgr: proxyMgr,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// StartBots launches all bots concurrently
func (bm *Manager) StartBots(clipURL, channelName, clipID string) error {
	var wg sync.WaitGroup

	// Set start time
	bm.stats.Mu.Lock()
	bm.stats.StartTime = time.Now()
	bm.stats.Mu.Unlock()

	// Launch bots
	for i := 1; i <= bm.config.BotCount; i++ {
		wg.Add(1)
		go bm.runBot(i, clipURL, channelName, clipID, &wg)

		// Small delay between bot starts
		time.Sleep(50 * time.Millisecond)
	}

	// Wait for all bots to complete
	wg.Wait()
	return nil
}

// Stop cancels all running bots
func (bm *Manager) Stop() {
	if bm.cancel != nil {
		bm.cancel()
	}
}

// GetContext returns the context for cancellation
func (bm *Manager) GetContext() context.Context {
	return bm.ctx
}

// runBot executes the bot logic for a single bot
func (bm *Manager) runBot(botID int, clipURL, channelName, clipID string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initialize bot
	bm.stats.SetActiveBot(botID)
	defer bm.stats.SetFinishedBot(botID)

	successCount := 0
	failCount := 0

	for viewNum := 1; viewNum <= bm.config.ViewsPerBot; viewNum++ {
		select {
		case <-bm.ctx.Done():
			// Context canceled, stop this bot
			bm.stats.UpdateBotStatus(botID, "Stopped")
			return
		default:
		}

		// Update bot status
		status := fmt.Sprintf("View %d/%d", viewNum, bm.config.ViewsPerBot)
		bm.stats.UpdateBotStatus(botID, status)

		// Perform view simulation
		useProxy := !bm.config.NoproxyMode
		err := bm.client.SimulateView(bm.ctx, clipURL, bm.proxyMgr, useProxy)

		if err != nil {
			failCount++
			bm.stats.IncrementFailed(botID, err.Error())
		} else {
			successCount++
			bm.stats.IncrementSuccess(botID)
		}

		// Random delay between requests
		minDelay, maxDelay := bm.config.GetDelayRange()
		delay := time.Duration(rand.Intn(maxDelay-minDelay+1)+minDelay) * time.Second

		select {
		case <-bm.ctx.Done():
			bm.stats.UpdateBotStatus(botID, "Stopped")
			return
		case <-time.After(delay):
		}
	}

	// Bot completed all views
	bm.stats.UpdateBotStatus(botID, "Finished")
}
