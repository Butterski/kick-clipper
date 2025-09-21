package models

import (
	"sync"
	"time"
)

// Proxy represents a proxy server configuration
type Proxy struct {
	IP       string
	Port     string
	Username string
	Password string
}

// BotStats tracks individual bot performance
type BotStats struct {
	ID         int
	Successful int
	Failed     int
	Status     string
	LastError  string
}

// GlobalStats tracks overall application metrics
type GlobalStats struct {
	Mu                sync.RWMutex
	TotalSuccessful   int
	TotalFailed       int
	InitialViews      int
	CurrentViews      int
	StartTime         time.Time
	Bots              map[int]*BotStats
	ActiveBots        int
	FinishedBots      int
	ProxiesCount      int
	TargetViews       int
	ViewsGained       int
	ViewsPerSecond    float64
	ViewsPerMinute    float64
	SuccessRate       float64
	EstimatedTimeLeft time.Duration
}

// NewGlobalStats creates a new GlobalStats instance
func NewGlobalStats() *GlobalStats {
	return &GlobalStats{
		Bots: make(map[int]*BotStats),
	}
}

// UpdateBotStats safely updates bot statistics
func (gs *GlobalStats) UpdateBotStats(botID int, successful, failed int, status, lastError string) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Bots[botID] == nil {
		gs.Bots[botID] = &BotStats{ID: botID}
	}

	bot := gs.Bots[botID]

	// Update totals
	gs.TotalSuccessful += successful - bot.Successful
	gs.TotalFailed += failed - bot.Failed

	// Update bot stats
	bot.Successful = successful
	bot.Failed = failed
	bot.Status = status
	bot.LastError = lastError
}

// IncrementSuccess safely increments successful views
func (gs *GlobalStats) IncrementSuccess(botID int) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Bots[botID] == nil {
		gs.Bots[botID] = &BotStats{ID: botID}
	}

	gs.Bots[botID].Successful++
	gs.TotalSuccessful++
}

// IncrementFailed safely increments failed views
func (gs *GlobalStats) IncrementFailed(botID int, error string) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Bots[botID] == nil {
		gs.Bots[botID] = &BotStats{ID: botID}
	}

	gs.Bots[botID].Failed++
	gs.Bots[botID].LastError = error
	gs.TotalFailed++
}

// UpdateBotStatus safely updates bot status
func (gs *GlobalStats) UpdateBotStatus(botID int, status string) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Bots[botID] == nil {
		gs.Bots[botID] = &BotStats{ID: botID}
	}

	gs.Bots[botID].Status = status
}

// GetStats returns a snapshot of current statistics
func (gs *GlobalStats) GetStats() (int, int, int, int, float64) {
	gs.Mu.RLock()
	defer gs.Mu.RUnlock()

	return gs.TotalSuccessful, gs.TotalFailed, gs.ActiveBots, gs.FinishedBots, gs.SuccessRate
}

// SetActiveBot sets a bot as active
func (gs *GlobalStats) SetActiveBot(botID int) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	if gs.Bots[botID] == nil {
		gs.Bots[botID] = &BotStats{ID: botID}
	}

	gs.ActiveBots++
}

// SetFinishedBot marks a bot as finished
func (gs *GlobalStats) SetFinishedBot(botID int) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	gs.ActiveBots--
	gs.FinishedBots++

	if gs.Bots[botID] != nil {
		gs.Bots[botID].Status = "Finished"
	}
}

// UpdateGlobalMetrics calculates and updates global metrics
func (gs *GlobalStats) UpdateGlobalMetrics(currentViews int) {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()

	gs.CurrentViews = currentViews

	// Calculate runtime
	runtime := time.Since(gs.StartTime)

	// Calculate success rate
	total := gs.TotalSuccessful + gs.TotalFailed
	if total > 0 {
		gs.SuccessRate = float64(gs.TotalSuccessful) / float64(total) * 100
	}

	// Calculate views gained and rates
	if gs.CurrentViews > 0 && gs.InitialViews > 0 {
		gs.ViewsGained = gs.CurrentViews - gs.InitialViews
		if runtime.Seconds() > 0 {
			gs.ViewsPerSecond = float64(gs.ViewsGained) / runtime.Seconds()
			gs.ViewsPerMinute = gs.ViewsPerSecond * 60
		}
	}

	// Estimate time left
	if gs.ViewsPerSecond > 0 && gs.TargetViews > 0 {
		remaining := gs.TargetViews - gs.TotalSuccessful
		if remaining > 0 {
			gs.EstimatedTimeLeft = time.Duration(float64(remaining)/gs.ViewsPerSecond) * time.Second
		}
	}
}

// ClipResponse represents the API response structure
type ClipResponse struct {
	Clip struct {
		ViewCount int `json:"view_count"`
	} `json:"clip"`
}
