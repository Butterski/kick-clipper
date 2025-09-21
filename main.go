package main

import (
	"bufio"
	"context"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"kick-clipper/internal/bot"
	"kick-clipper/internal/client"
	"kick-clipper/internal/config"
	"kick-clipper/internal/dashboard"
	"kick-clipper/internal/models"
	"kick-clipper/internal/proxy"
	"kick-clipper/pkg/color"
	"kick-clipper/pkg/utils"
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create application
	app := NewApplication()

	// Execute CLI
	if err := app.Execute(); err != nil {
		dashboard.ShowErrorMessage("Command execution failed: %v", err)
		os.Exit(1)
	}
}

// Application holds the main application state
type Application struct {
	config     *config.Config
	stats      *models.GlobalStats
	proxyMgr   *proxy.Manager
	botMgr     *bot.Manager
	kickClient *client.KickClient
	dashboard  *dashboard.Dashboard
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	cfg := config.DefaultConfig()
	stats := models.NewGlobalStats()
	proxyMgr := proxy.NewManager()
	dash := dashboard.NewDashboard(stats)

	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		config:    cfg,
		stats:     stats,
		proxyMgr:  proxyMgr,
		dashboard: dash,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Execute runs the application with CLI interface
func (app *Application) Execute() error {
	// Parse command line arguments manually
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--proxy-file":
			if i+1 < len(args) {
				app.config.ProxyFile = args[i+1]
				i++
			}
		case "--no-proxy":
			app.config.NoproxyMode = true
		case "--bots":
			if i+1 < len(args) {
				if count, err := strconv.Atoi(args[i+1]); err == nil {
					app.config.BotCount = count
				}
				i++
			}
		case "--views-per-bot":
			if i+1 < len(args) {
				if views, err := strconv.Atoi(args[i+1]); err == nil {
					app.config.ViewsPerBot = views
				}
				i++
			}
		case "--url":
			if i+1 < len(args) {
				app.config.ClipURL = args[i+1]
				i++
			}
		case "--no-interactive":
			app.config.NoInteractive = true
		case "--timeout":
			if i+1 < len(args) {
				if timeout, err := strconv.Atoi(args[i+1]); err == nil {
					app.config.Timeout = timeout
				}
				i++
			}
		case "--help", "-h":
			app.showHelp()
			return nil
		}
	}

	return app.run()
}

// showHelp displays usage information
func (app *Application) showHelp() {
	color.Cyan("🚀 Kick Clipper - High-performance Kick.com clip view bot")
	color.White("\nUsage:")
	color.White("  kick-clipper [flags]")
	color.White("\nFlags:")
	color.White("  --proxy-file string     Path to proxy file (default \"proxies.txt\")")
	color.White("  --no-proxy              Run without proxies")
	color.White("  --bots int              Number of bots to run")
	color.White("  --views-per-bot int     Views per bot")
	color.White("  --url string            Kick.com clip URL")
	color.White("  --no-interactive        Run without interactive prompts")
	color.White("  --timeout int           HTTP request timeout in seconds (default 15)")
	color.White("  --help, -h              Show this help message")
}

// run executes the main application logic
func (app *Application) run() error {
	// Setup signal handling
	app.setupSignalHandling()

	// Initialize application
	if err := app.initialize(); err != nil {
		return err
	}

	// Run interactive mode if needed
	if !app.config.NoInteractive {
		if err := app.runInteractiveSetup(); err != nil {
			return err
		}
	}

	// Validate configuration
	if err := app.validateConfig(); err != nil {
		return err
	}

	// Parse clip URL
	channelName, clipID, err := app.kickClient.ParseClipURL(app.config.ClipURL)
	if err != nil {
		return err
	}

	dashboard.ShowSuccessMessage("Channel: %s", channelName)
	dashboard.ShowSuccessMessage("Clip ID: %s", clipID)

	// Get initial clip views
	if err := app.getInitialViews(clipID); err != nil {
		dashboard.ShowWarningMessage("Could not fetch initial views: %v", err)
	}

	// Start the application
	return app.start(channelName, clipID)
}

// initialize sets up the application components
func (app *Application) initialize() error {
	// Initialize HTTP client
	timeout := time.Duration(app.config.Timeout) * time.Second
	app.kickClient = client.NewKickClient(timeout)

	// Load proxies if not in no-proxy mode
	if !app.config.NoproxyMode {
		if err := app.proxyMgr.LoadFromFile(app.config.ProxyFile); err != nil {
			return err
		}

		if app.proxyMgr.IsEmpty() {
			return dashboard.ShowErrorMessage("No valid proxies found in %s", app.config.ProxyFile)
		}

		dashboard.ShowSuccessMessage("Loaded %d proxies", app.proxyMgr.Count())
		app.stats.ProxiesCount = app.proxyMgr.Count()
	}

	// Initialize bot manager
	app.botMgr = bot.NewManager(app.config, app.stats, app.proxyMgr)

	return nil
}

// runInteractiveSetup handles interactive user input
func (app *Application) runInteractiveSetup() error {
	dashboard.ShowStartupMessage()

	// Get clip URL if not provided
	if app.config.ClipURL == "" {
		if url, err := app.promptForClipURL(); err != nil {
			return err
		} else {
			app.config.ClipURL = url
		}
	}

	// Auto-calculate bot count if not provided
	if app.config.BotCount == 0 {
		autoCalculated := app.calculateOptimalBotCount()
		if count, err := app.promptForBotCount(autoCalculated); err != nil {
			return err
		} else {
			app.config.BotCount = count
		}
	}

	// Get views per bot if not provided
	if app.config.ViewsPerBot == 0 {
		if views, err := app.promptForViewsPerBot(); err != nil {
			return err
		} else {
			app.config.ViewsPerBot = views
		}
	}

	// Show configuration summary
	app.stats.TargetViews = app.config.BotCount * app.config.ViewsPerBot
	dashboard.ShowSuccessMessage("Using %d bots", app.config.BotCount)
	dashboard.ShowInfoMessage("Target total views: %s", utils.FormatNumber(app.stats.TargetViews))

	// Warning for no-proxy mode
	if app.config.NoproxyMode {
		return app.showNoproxyWarning()
	}

	return nil
}

// validateConfig ensures the configuration is valid
func (app *Application) validateConfig() error {
	if app.config.ClipURL == "" {
		return dashboard.ShowErrorMessage("Clip URL is required")
	}

	if app.config.BotCount <= 0 {
		return dashboard.ShowErrorMessage("Bot count must be positive")
	}

	if app.config.ViewsPerBot <= 0 {
		return dashboard.ShowErrorMessage("Views per bot must be positive")
	}

	return nil
}

// getInitialViews fetches and stores initial clip view count
func (app *Application) getInitialViews(clipID string) error {
	dashboard.ShowInfoMessage("Getting initial clip views...")

	useProxy := !app.config.NoproxyMode
	views, err := app.kickClient.GetClipViews(app.ctx, clipID, app.proxyMgr, useProxy)
	if err != nil {
		return err
	}

	app.stats.InitialViews = views
	app.stats.CurrentViews = views
	dashboard.ShowSuccessMessage("Initial views: %s", utils.FormatNumber(views))

	return nil
}

// start begins the bot execution and dashboard
func (app *Application) start(channelName, clipID string) error {
	dashboard.ShowInfoMessage("Starting view bot army...")
	color.Cyan("═══════════════════════════════════════════════")

	// Start dashboard in background
	go app.runDashboard(clipID)

	// Start bots
	if err := app.botMgr.StartBots(app.config.ClipURL, channelName, clipID); err != nil {
		return err
	}

	dashboard.ShowSuccessMessage("All bots finished!")
	return nil
}

// runDashboard runs the dashboard update loop
func (app *Application) runDashboard(clipID string) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-app.ctx.Done():
			return
		case <-ticker.C:
			// Update current views
			useProxy := !app.config.NoproxyMode
			if views, err := app.kickClient.GetClipViews(app.ctx, clipID, app.proxyMgr, useProxy); err == nil {
				app.stats.UpdateGlobalMetrics(views)
			}

			// Display dashboard
			app.dashboard.Clear()
			app.dashboard.Display(clipID, app.config)

			// Check if all bots finished
			_, _, _, finished, _ := app.stats.GetStats()
			if finished >= app.config.BotCount {
				color.Green("\n🏁 ALL BOTS FINISHED!")
				return
			}
		}
	}
}

// setupSignalHandling configures graceful shutdown
func (app *Application) setupSignalHandling() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		dashboard.ShowWarningMessage("Stopping bots...")
		app.cancel()
		if app.botMgr != nil {
			app.botMgr.Stop()
		}
		dashboard.ShowInfoMessage("Bot army stopped!")
		os.Exit(0)
	}()
}

// Interactive prompt functions
func (app *Application) promptForClipURL() (string, error) {
	color.Cyan("📎 Enter Kick.com clip URL: ")
	reader := bufio.NewReader(os.Stdin)
	clipURL, _ := reader.ReadString('\n')
	return strings.TrimSpace(clipURL), nil
}

func (app *Application) calculateOptimalBotCount() int {
	if !app.config.NoproxyMode && app.proxyMgr.Count() > 0 {
		optimal := app.proxyMgr.Count() * 5
		dashboard.ShowInfoMessage("Auto-calculated bots: %d (proxies × 5)", optimal)
		return optimal
	}
	return 10 // Default for no-proxy mode
}

func (app *Application) promptForBotCount(autoCalculated int) (int, error) {
	color.Cyan("🤖 Number of bots (press Enter for %d): ", autoCalculated)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return autoCalculated, nil
	}

	if count, err := utils.ValidatePositiveInt(input); err == nil {
		if !app.config.NoproxyMode && count > autoCalculated {
			dashboard.ShowWarningMessage("%d bots with %d proxies may cause issues", count, app.proxyMgr.Count())
		}
		return count, nil
	}

	dashboard.ShowWarningMessage("Invalid input. Using auto-calculated value.")
	return autoCalculated, nil
}

func (app *Application) promptForViewsPerBot() (int, error) {
	color.Cyan("👁️  Views per bot: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	views, err := strconv.Atoi(input)
	if err != nil || views <= 0 {
		dashboard.ShowErrorMessage("Invalid input. Please enter a positive number.")
		return app.promptForViewsPerBot()
	}

	return views, nil
}

func (app *Application) showNoproxyWarning() error {
	dashboard.ShowWarningMessage("All requests will come from your IP address!")
	dashboard.ShowWarningMessage("Consider using fewer bots to avoid rate limiting.")

	color.Cyan("Continue? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')

	if strings.ToLower(strings.TrimSpace(confirm)) != "y" {
		dashboard.ShowInfoMessage("Cancelled")
		os.Exit(0)
	}

	return nil
}
