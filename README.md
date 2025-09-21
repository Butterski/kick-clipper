# 🚀 Kick.com Clip View Bot

A high-performance Go-based bot for generating views on Kick.com clips with superior performance and Cloudflare bypass capabilities.

## 📁 Project Structure

```
kick-clipper-go/
├── internal/            # Internal packages
│   ├── bot/            # Bot management
│   ├── client/         # Kick.com client
│   ├── config/         # Configuration handling
│   ├── dashboard/      # Terminal dashboard
│   ├── models/         # Data models
│   └── proxy/          # Proxy management
├── pkg/                # Public packages
│   ├── color/          # Color utilities
│   └── utils/          # Utility functions
├── main.go             # Main application entry point
├── proxies.txt         # Proxy list file (required)
├── Dockerfile          # Container build instructions
├── docker-compose.yml  # Container orchestration
├── env.example         # Environment configuration template
└── README.md           # This file
```

## 🚀 Quick Start

### Using Docker (Easiest)

```bash
docker-compose up --build
```

### Manual Installation

1. **Install Go 1.21+**
   ```bash
   # Download from https://golang.org/dl/
   ```

2. **Build and Run**
   ```bash
   go build -o kick-clipper.exe .
   ./kick-clipper.exe -clip <CLIP_ID> -workers 100 -delay 5 -time 300
   ```

## ✨ Features
- **High Performance**: Native compiled binary with efficient goroutine-based concurrency
- **Advanced TLS Fingerprinting**: Chrome-like TLS configuration bypasses Cloudflare
- **Low Memory Usage**: Optimized resource utilization
- **Cross-platform**: Single binary deployment
- **Real-time Dashboard**: Enhanced colorized terminal interface with live metrics

### 🛠️ Configuration Options

```bash
./kick-clipper.exe [OPTIONS]

Options:
  -clip string     Clip ID or full URL (required)
  -workers int     Number of concurrent bots (default: 100)
  -delay int       Delay between requests in seconds (default: 5)  
  -time int        Runtime in seconds (default: 300)
  -target int      Target view count (default: 10000)
  -proxy           Enable proxy rotation (default: true)
```

### 📊 Dashboard Features
- Real-time view statistics
- Bot success/failure rates
- Proxy rotation status
- Performance metrics (requests/sec)
- Progress tracking with visual indicators

## 🐳 Docker Deployment

### Quick Start with Docker Compose

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd kick-clipper
   cp .env.example .env
   # Edit .env with your clip ID and settings
   ```

2. **Run with Docker Compose:**
   ```bash
   # Production run
   docker-compose up --build
   
   # Development mode with custom settings
   docker-compose --profile dev up kick-clipper-dev
   ```

3. **Custom configuration:**
   ```bash
   # Override environment variables
   CLIP_ID=your_clip_id WORKERS=200 docker-compose up
   ```

### Manual Docker Build

```bash
# Build image
docker build -t kick-clipper .

# Run container
docker run -d \
  --name kick-clipper-bot \
  -v $(pwd)/proxies.txt:/app/proxies.txt:ro \
  kick-clipper \
  -clip your_clip_id -workers 100 -delay 5 -time 300
```

## 🛠️ Manual Installation

1. **Prepare proxy list:**
   Create a `proxies.txt` file with one proxy per line in format:
   ```
   ip:port:username:password
   ip:port:username:password
   ...
   ```

2. **Build and run the application:**
   ```bash
   go build -o kick-clipper.exe .
   ./kick-clipper.exe -clip <CLIP_ID> -workers 100 -delay 5 -time 300
   ```

3. **Monitor the real-time dashboard:**
   - View real-time statistics
   - Track bot success/failure rates
   - Monitor proxy rotation status
   - See performance metrics (requests/sec)

## 📊 Dashboard Features

Both scripts include a real-time dashboard that displays:

- **Runtime Information:** Current session duration
- **View Statistics:** Current views, initial views, views gained, rate per minute/second
- **Bot Statistics:** Active/finished bots, success/failure rates
- **Progress Tracking:** Visual progress bar towards target views
- **Proxy Information:** Available proxy count (main.py only)

## ⚙️ Configuration

### Proxy Format (main.py)
```
ip:port:username:password
```
Example:
```
192.168.1.100:8080:user1:pass1
10.0.0.50:3128:user2:pass2
```

### Bot Configuration
- **main.py:** Auto-calculated as `proxy_count × 5` (user can override for smaller machines)
- **main_no_proxy.py:** User-defined (recommended: 5-20)

### Timing Configuration
- **Proxy version:** 2-8 seconds between requests
- **No-proxy version:** 5-15 seconds between requests (longer to avoid rate limiting)

## 🔗 Supported URL Format

Both scripts support Kick.com clip URLs in the format:
```
https://kick.com/CHANNEL/clips/CLIP_ID
```

Example:
```
https://kick.com/streamer123/clips/abc123def456
```

## ⚠️ Important Notes

### Legal and Ethical Considerations
- This tool is for educational purposes only
- Ensure compliance with Kick.com's Terms of Service
- Use responsibly and respect platform guidelines
- Consider the impact on content creators and platform integrity

### Security and Privacy
- **Proxy version:** Distributes requests across multiple IP addresses
- **No-proxy version:** All requests originate from your IP address
- Monitor your network usage and respect rate limits

### Rate Limiting
- **Proxy version:** Less likely to hit rate limits due to IP distribution
- **No-proxy version:** Higher risk of rate limiting from single IP
- Both versions include built-in delays between requests

## 🛠️ Troubleshooting

### Common Issues

1. **"No proxies found" error:**
   - Ensure `proxies.txt` exists and contains valid proxies
   - Check proxy format: `ip:port:username:password`

2. **"Invalid URL format" error:**
   - Verify the clip URL follows the correct format
   - Ensure the URL is a valid Kick.com clip link

3. **Connection timeouts:**
   - Check your internet connection
   - Verify proxy credentials (if using proxy version)
   - Try reducing the number of concurrent bots

4. **High failure rate:**
   - Validate proxy list quality
   - Reduce request frequency
   - Check if target clip still exists

### Performance Optimization

- **Proxy version:** More proxies = better performance and lower detection risk
- **No-proxy version:** Use fewer bots (5-10) to avoid overwhelming your connection
- Monitor success rates and adjust bot count accordingly

## 📈 Best Practices

1. **Start Small:** Begin with fewer bots and gradually increase
2. **Monitor Performance:** Watch success rates and adjust accordingly
3. **Respect Limits:** Don't overwhelm the target server
4. **Use Quality Proxies:** Better proxies = better success rates
5. **Test First:** Use the no-proxy version for initial testing

## 🔄 Updates and Maintenance

- Regularly update proxy lists for optimal performance
- Monitor script output for any errors or warnings
- Keep the `curl_cffi` library updated
- Check Kick.com's Terms of Service for any changes

## 📞 Support

If you encounter issues:
1. Check the troubleshooting section above
2. Verify all requirements are met
3. Ensure proper file formats and configurations
4. Monitor console output for specific error messages

---

**Disclaimer:** This tool is provided for educational purposes only. Users are responsible for ensuring compliance with all applicable terms of service and laws. Use at your own risk.