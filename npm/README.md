# Tunny - Fast HTTP Tunnel CLI

**Tunny** is a simple, fast CLI tool for creating secure tunnels to expose your local web server to the internet.

## 🚀 Quick Start

### Installation

```bash
npm install -g tunny-tunnel
```

### Usage

```bash
# Start your local server
npm run dev  # or any local server on localhost:3000

# Create a tunnel (requires token from server admin)
tunny connect localhost:3000 --token YOUR_TOKEN --subdomain myapp
```

**That's it!** Your local server is now accessible via a public URL.

## 📖 Commands

```bash
# Connect with token and subdomain
tunny connect localhost:3000 --token dev --subdomain myapp

# With custom tunnel ID
tunny connect localhost:8080 --token dev --subdomain api --id my-tunnel

# List active tunnels
tunny list

# Initialize config (save defaults)
tunny init --token YOUR_TOKEN --subdomain myapp

# After init, just run:
tunny connect localhost:3000
```

## ⚙️ Configuration

### Option 1: Environment Variables

```bash
export TUNNY_TOKEN="your-token"
export TUNNY_SUBDOMAIN="your-subdomain"

tunny connect localhost:3000
```

### Option 2: Config File

```bash
tunny init --token YOUR_TOKEN --subdomain myapp
```

Config saved to `~/.tunny/config.json`

## 🎯 Use Cases

- **Webhook Testing** - Test webhooks from Stripe, GitHub, etc.
- **Mobile App Development** - Test your local API from mobile devices
- **Client Demos** - Share your work-in-progress
- **IoT Development** - Receive callbacks from IoT devices
- **Quick Prototyping** - No deployment needed

## 🔐 Authentication

You need an authentication token to use Tunny. The server URL is pre-configured and points to the production server.

**Get your token:**
- Contact your server administrator
- Or use `dev` token for development servers

## 📊 Examples

### React App

```bash
npm start  # Runs on localhost:3000
tunny connect localhost:3000 --token dev --subdomain my-react-app
```

### Node.js API

```bash
node server.js  # Listening on port 8080
tunny connect localhost:8080 --token dev --subdomain my-api
```

### Python Flask

```bash
flask run  # Runs on localhost:5000
tunny connect localhost:5000 --token dev --subdomain flask-app
```

## ⚠️ Important Notes

- This package contains **client-only** code
- The server URL is **hardcoded** to production
- You cannot self-host using this NPM package
- For self-hosting, see the main GitHub repository

## 🔗 Links

- **GitHub**: https://github.com/kayaramazan/tunny
- **Issues**: https://github.com/kayaramazan/tunny-client/issues
- **Documentation**: https://github.com/kayaramazan/tunny#readme

## 📄 License

MIT

---

**Note**: This package automatically downloads the appropriate binary for your platform during installation.

