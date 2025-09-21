package dashboard

import (
	"fmt"
	"time"

	"kick-clipper/internal/models"
	"kick-clipper/pkg/color"
	"kick-clipper/pkg/utils"
)

// Dashboard handles the display of real-time statistics
type Dashboard struct {
	stats *models.GlobalStats
}

// NewDashboard creates a new dashboard instance
func NewDashboard(stats *models.GlobalStats) *Dashboard {
	return &Dashboard{
		stats: stats,
	}
}

// Display renders the enhanced dashboard
func (d *Dashboard) Display(clipID string, config interface{}) {
	d.stats.Mu.RLock()
	defer d.stats.Mu.RUnlock()

	// Header
	headerColor := color.New(color.FgHiCyan, color.Bold)
	headerColor.Println("🚀 KICK.COM CLIP VIEW BOT - GOLANG EDITION 🚀")
	color.Cyan("═══════════════════════════════════════════════════════════════════════")

	// Runtime information
	d.displayRuntimeInfo(clipID)

	// View statistics
	d.displayViewStats()

	// Bot statistics
	d.displayBotStats()

	// Progress bar
	d.displayProgress()

	// Configuration info
	d.displayConfigInfo(config)

	// Performance metrics
	d.displayPerformanceMetrics()

	// Recent bot activity
	d.displayBotActivity()

	color.Cyan("\n⏹️  Press Ctrl+C to stop")
	color.Cyan("═══════════════════════════════════════════════════════════════════════")
}

// displayRuntimeInfo shows runtime and ETA information
func (d *Dashboard) displayRuntimeInfo(clipID string) {
	runtime := time.Since(d.stats.StartTime)
	hours := int(runtime.Hours())
	minutes := int(runtime.Minutes()) % 60
	seconds := int(runtime.Seconds()) % 60

	color.White("📎 Clip ID: %s", clipID)
	color.White("⏱️  Runtime: %02d:%02d:%02d", hours, minutes, seconds)

	if d.stats.EstimatedTimeLeft > 0 {
		color.White("⏳ ETA: %v", d.stats.EstimatedTimeLeft.Round(time.Second))
	}

	color.Cyan("═══════════════════════════════════════════════════════════════════════")
}

// displayViewStats shows view-related statistics
func (d *Dashboard) displayViewStats() {
	color.Yellow("\n📊 VIEW STATISTICS:")

	if d.stats.CurrentViews > 0 {
		color.Green("   Current Views: %s", utils.FormatNumber(d.stats.CurrentViews))
		color.Green("   Initial Views: %s", utils.FormatNumber(d.stats.InitialViews))
		color.Green("   Views Gained:  %s (+%d)", utils.FormatNumber(d.stats.ViewsGained), d.stats.ViewsGained)
		color.Green("   Rate: %.1f/min (%.2f/sec)", d.stats.ViewsPerMinute, d.stats.ViewsPerSecond)
	} else {
		color.Red("   Current Views: Unable to fetch")
	}
}

// displayBotStats shows bot-related statistics
func (d *Dashboard) displayBotStats() {
	color.Yellow("\n🤖 BOT STATISTICS:")
	color.White("   Total Bots: %d", len(d.stats.Bots))
	color.White("   Active: %d | Finished: %d", d.stats.ActiveBots, d.stats.FinishedBots)
	color.Green("   Bot Success: %s", utils.FormatNumber(d.stats.TotalSuccessful))
	color.Red("   Bot Failed:  %s", utils.FormatNumber(d.stats.TotalFailed))
	color.White("   Success Rate: %.1f%%", d.stats.SuccessRate)
}

// displayProgress shows progress bar and target information
func (d *Dashboard) displayProgress() {
	color.Yellow("\n🎯 PROGRESS:")

	progress := float64(0)
	if d.stats.TargetViews > 0 {
		progress = float64(d.stats.TotalSuccessful) / float64(d.stats.TargetViews) * 100
	}

	progressBar := utils.GenerateProgressBar(int(progress), 50)
	color.White("   Target: %s views", utils.FormatNumber(d.stats.TargetViews))
	color.White("   Progress: [%s] %.1f%%", progressBar, progress)
}

// displayConfigInfo shows configuration information
func (d *Dashboard) displayConfigInfo(config interface{}) {
	// This would be enhanced based on the actual config type
	// For now, we'll show basic proxy information
	if d.stats.ProxiesCount > 0 {
		color.Yellow("\n🌐 PROXY INFO:")
		color.White("   Available Proxies: %d", d.stats.ProxiesCount)
		color.White("   Mode: Proxy Enabled")
	} else {
		color.Yellow("\n⚠️  NO-PROXY MODE:")
		color.Red("   All requests from your IP address")
		color.Red("   Consider using fewer bots to avoid rate limiting")
	}
}

// displayPerformanceMetrics shows performance information
func (d *Dashboard) displayPerformanceMetrics() {
	color.Yellow("\n📈 PERFORMANCE:")

	runtime := time.Since(d.stats.StartTime)
	if runtime.Seconds() > 0 {
		requestsPerSecond := float64(d.stats.TotalSuccessful+d.stats.TotalFailed) / runtime.Seconds()
		color.White("   Requests/sec: %.2f", requestsPerSecond)
		color.White("   Total Requests: %s", utils.FormatNumber(d.stats.TotalSuccessful+d.stats.TotalFailed))
	}
}

// displayBotActivity shows recent bot errors and activity
func (d *Dashboard) displayBotActivity() {
	color.Yellow("\n🔍 RECENT BOT ACTIVITY:")

	errorCount := 0
	for _, bot := range d.stats.Bots {
		if bot.LastError != "" && errorCount < 3 {
			errorMsg := utils.TruncateString(bot.LastError, 40)
			color.Red("   Bot #%d: %s (%s)", bot.ID, bot.Status, errorMsg)
			errorCount++
		}
	}

	if errorCount == 0 {
		color.Green("   All bots running smoothly ✅")
	}
}

// Clear clears the screen
func (d *Dashboard) Clear() {
	utils.ClearScreen()
}

// ShowStartupMessage displays initial startup information
func ShowStartupMessage() {
	color.Green("🚀 KICK.COM CLIP VIEW BOT - GOLANG EDITION")
	color.Cyan("═══════════════════════════════════════════════")
}

// ShowSuccessMessage displays a success message with formatting
func ShowSuccessMessage(message string, args ...interface{}) {
	color.Green("✅ "+message, args...)
}

// ShowWarningMessage displays a warning message with formatting
func ShowWarningMessage(message string, args ...interface{}) {
	color.Yellow("⚠️  "+message, args...)
}

// ShowErrorMessage displays an error message with formatting
func ShowErrorMessage(message string, args ...interface{}) error {
	color.Red("❌ "+message, args...)
	return fmt.Errorf(message, args...)
}

// ShowInfoMessage displays an info message with formatting
func ShowInfoMessage(message string, args ...interface{}) {
	color.White("ℹ️  "+message, args...)
}
