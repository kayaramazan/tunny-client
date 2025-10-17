# Tunny - Fast & Simple HTTP Tunnel Client

**Tunny** is a lightweight CLI tool that creates secure tunnels to expose your local web server to the internet. Perfect for demos, webhooks testing, and quick prototyping.

## ✨ Features

- 🚀 **Zero Configuration** - Works out of the box
- 🔒 **Secure** - HTTPS endpoints by default  
- ⚡ **Fast** - Minimal latency overhead
- 🎯 **Simple** - One command to start tunneling
- 🌐 **Cross-Platform** - macOS, Linux, Windows

## 📦 Installation

### Via NPM (Recommended)

```bash
npm install -g tunny-tunnel
```

### Via Go

```bash
go install github.com/kayaramazan/tunny/cmd/tunny@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/kayaramazan/tunny/releases)

## 🚀 Quick Start

### 1. Start your local server

```bash
# Example: Start a simple web server on port 3000
python3 -m http.server 3000
```

### 2. Create a tunnel

```bash
tunny connect localhost:3000 --token YOUR_TOKEN --subdomain myapp
```

### 3. Access your app

```
🌐 Public URL: https://your-server.com/abc123/
```

Your local server is now accessible from the internet!

## 📖 Usage

### Basic Command

```bash
tunny connect [target] [flags]
```

### Examples

```bash
# Basic usage with token
tunny connect localhost:3000 --token dev --subdomain myapp

# Custom tunnel ID
tunny connect localhost:8080 --token dev --subdomain api --id my-tunnel

# With environment variables
export TUNNY_TOKEN=dev
export TUNNY_SUBDOMAIN=myapp
tunny connect localhost:3000

# Dev mode with verbose logging
tunny connect localhost:3000 --token dev --subdomain test --dev
```

### Configuration

Set default values to avoid repeating flags:

```bash
tunny init --token YOUR_TOKEN --subdomain myapp
```

Config is saved to `~/.tunny/config.json`

### Environment Variables

```bash
export TUNNY_TOKEN="your-token"
export TUNNY_SUBDOMAIN="your-subdomain"
```

### List Active Tunnels

```bash
tunny list
```

## 🔧 Commands

| Command | Description |
|---------|-------------|
| `tunny connect [target]` | Create a tunnel to your local server |
| `tunny list` | List all active tunnels |
| `tunny init` | Initialize configuration |
| `tunny --version` | Show version info |
| `tunny --help` | Show help |

## 🎯 Use Cases

- **Webhook Development** - Test webhooks from services like Stripe, GitHub, etc.
- **Mobile App Testing** - Test your backend API from mobile devices
- **Client Demos** - Share your local development with clients
- **IoT Callbacks** - Receive callbacks from IoT devices
- **SSH Access** - Tunnel SSH connections (with TCP support)
- **Quick Prototypes** - Share prototypes without deployment

## 🛠️ Flags

```
Flags:
  -d, --dev                Enable development mode with verbose logging
  -h, --help               Help for connect
  -i, --id string          Custom tunnel ID (auto-generated if not provided)
      --subdomain string   Subdomain for the tunnel
  -t, --token string       Authentication token
```

## 🔐 Authentication

Tunny requires an authentication token. Contact the server administrator to get your token, or use `dev` for development servers.

## 📊 Examples

### React Development

```bash
# Start React app
npm start  # Running on localhost:3000

# In another terminal
tunny connect localhost:3000 --token dev --subdomain react-app
```

### Python Flask

```bash
# Start Flask app
flask run --port 5000

# Create tunnel
tunny connect localhost:5000 --token dev --subdomain flask-api
```

### Node.js Express

```bash
# Start Express
node server.js  # Listening on port 8080

# Tunnel it
tunny connect localhost:8080 --token dev --subdomain node-api
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details

## 🔗 Links

- **GitHub**: https://github.com/kayaramazan/tunny
- **NPM**: https://www.npmjs.com/package/tunny-tunnel
- **Issues**: https://github.com/kayaramazan/tunny/issues

## ⚠️ Note

This is the **client-only** package. The tunnel server is hosted separately. For self-hosting the server, contact the maintainer.

---

Made with ❤️ by the Tunny team

