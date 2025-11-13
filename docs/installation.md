# Installation Guide

This guide will walk you through installing and setting up PostgresMaster on your system.

## System Requirements

### Local Proxy Binary

- **Operating System**: Windows 10+, macOS 10.15+, or Linux (any modern distribution)
- **Architecture**: x64 (amd64) or ARM64 (Apple Silicon)
- **Memory**: 50MB RAM (typical usage)
- **Disk Space**: 10MB for the binary
- **Network**: Outbound connection to your PostgreSQL server

### Web Application

- **Browser** (one of):
  - Chrome/Edge 90+
  - Firefox 88+
  - Safari 14+
  - Opera 76+
- **JavaScript**: Must be enabled
- **WebSocket Support**: Required (all modern browsers support this)

### PostgreSQL Server

- **Version**: PostgreSQL 10.0 or later
- **Network Access**: The proxy must be able to reach your PostgreSQL server
- **Authentication**: Valid PostgreSQL credentials (username/password)

## Download Instructions

### Option 1: Download from GitHub Releases (Recommended)

1. Visit the [Releases page](https://github.com/MPJHorner/PostgresMaster/releases)
2. Download the latest release for your platform:

| Platform | Architecture | Filename |
|----------|-------------|----------|
| Windows | x64 | `postgres-proxy-windows-amd64.exe` |
| macOS | Intel | `postgres-proxy-darwin-amd64` |
| macOS | Apple Silicon (M1/M2/M3) | `postgres-proxy-darwin-arm64` |
| Linux | x64 | `postgres-proxy-linux-amd64` |
| Linux | ARM64 | `postgres-proxy-linux-arm64` |

3. (Optional) Download `checksums.txt` to verify file integrity

### Option 2: Build from Source

If you prefer to build from source or need to customize the proxy:

```bash
# Clone the repository
git clone https://github.com/MPJHorner/PostgresMaster.git
cd PostgresMaster/proxy

# Build for your current platform
make build

# Or build for all platforms
make build-all

# Binaries will be in the bin/ directory
```

**Prerequisites for building**:
- Go 1.21 or later
- Make (optional, but recommended)

## Platform-Specific Installation

### Windows

#### Installation Steps

1. Download `postgres-proxy-windows-amd64.exe`
2. Move the file to a permanent location (e.g., `C:\Program Files\PostgresMaster\`)
3. (Optional) Add the directory to your system PATH:
   - Right-click "This PC" → Properties → Advanced system settings
   - Click "Environment Variables"
   - Under "System variables", select "Path" and click "Edit"
   - Click "New" and add your proxy directory path
   - Click "OK" to save

#### Verify Installation

Open Command Prompt or PowerShell:

```cmd
# If added to PATH
postgres-proxy --version

# Or use the full path
"C:\Program Files\PostgresMaster\postgres-proxy-windows-amd64.exe" --version
```

#### Security Note

Windows may show a "Windows protected your PC" warning when first running the binary. This is because the binary is not signed with a Windows code signing certificate. Click "More info" → "Run anyway" to proceed.

### macOS

#### Installation Steps

1. Download the appropriate binary:
   - Intel Macs: `postgres-proxy-darwin-amd64`
   - Apple Silicon (M1/M2/M3): `postgres-proxy-darwin-arm64`

2. Open Terminal and navigate to your Downloads folder:
   ```bash
   cd ~/Downloads
   ```

3. Make the binary executable:
   ```bash
   chmod +x postgres-proxy-darwin-*
   ```

4. (Optional) Move to a standard location:
   ```bash
   sudo mv postgres-proxy-darwin-* /usr/local/bin/postgres-proxy
   ```

#### Verify Installation

```bash
postgres-proxy --version
```

#### Security Note - Gatekeeper

On first run, macOS Gatekeeper may block the binary with a message like "cannot be opened because the developer cannot be verified."

**Solution**:

1. Go to System Preferences → Security & Privacy → General
2. You should see a message about postgres-proxy being blocked
3. Click "Open Anyway"

Or use the command line:

```bash
# Remove the quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/postgres-proxy
```

### Linux

#### Installation Steps

1. Download the appropriate binary:
   - x64 systems: `postgres-proxy-linux-amd64`
   - ARM64 systems: `postgres-proxy-linux-arm64`

2. Open a terminal and navigate to your downloads:
   ```bash
   cd ~/Downloads
   ```

3. Make the binary executable:
   ```bash
   chmod +x postgres-proxy-linux-*
   ```

4. (Recommended) Move to a standard location:
   ```bash
   sudo mv postgres-proxy-linux-* /usr/local/bin/postgres-proxy
   ```

5. (Optional) Verify the binary with checksums:
   ```bash
   sha256sum postgres-proxy-linux-amd64
   # Compare with checksums.txt from the release
   ```

#### Verify Installation

```bash
postgres-proxy --version
```

#### Linux Distribution Notes

- **Ubuntu/Debian**: Works out of the box
- **RHEL/CentOS/Fedora**: Works out of the box
- **Arch Linux**: Works out of the box
- **Alpine Linux**: The binary is dynamically linked to glibc. If you're using Alpine (which uses musl), you may need to install the `libc6-compat` package

## Verifying Your Installation

After installation, verify everything works:

```bash
# Check version
postgres-proxy --version

# Check help
postgres-proxy --help
```

Expected output:
```
PostgresMaster Proxy v0.1.0
A WebSocket-to-PostgreSQL proxy for browser-based database clients
```

## First Connection

### Using a Connection String

```bash
postgres-proxy "postgres://username:password@hostname:5432/database"
```

Example:
```bash
postgres-proxy "postgres://admin:secret@db.example.com:5432/myapp"
```

### Using Interactive Mode

Simply run the proxy without arguments:

```bash
postgres-proxy
```

You'll be prompted for:
- **Host**: The PostgreSQL server hostname or IP address
- **Port**: The PostgreSQL port (default: 5432)
- **Database**: The database name to connect to
- **Username**: Your PostgreSQL username
- **Password**: Your PostgreSQL password (hidden)
- **SSL Mode**: SSL connection mode (default: prefer)

### SSL Modes Explained

| Mode | Description |
|------|-------------|
| `disable` | No SSL encryption |
| `allow` | Try non-SSL first, then SSL if that fails |
| `prefer` | Try SSL first, then non-SSL if that fails (default) |
| `require` | Only SSL connections (fails if SSL unavailable) |
| `verify-ca` | Require SSL and verify server certificate against CA |
| `verify-full` | Require SSL, verify CA, and verify hostname matches certificate |

**Recommendation**: Use `require` or higher for production databases.

## What Happens After Connection?

On successful connection, you'll see:

```
Connecting to postgres://db.example.com:5432/myapp...
✓ Connected!

→ Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6789...
   Proxy running on :8080
   Press Ctrl+C to stop
```

**Important**: The secret in the URL is unique to this session and required for the web app to connect to the proxy. Keep it confidential and don't share it.

## Troubleshooting

### Connection Issues

#### Problem: "connection refused"

**Possible causes**:
- PostgreSQL server is not running
- Firewall is blocking the connection
- Wrong host or port

**Solutions**:
1. Verify PostgreSQL is running: `pg_isready -h hostname -p 5432`
2. Check your firewall settings
3. Verify the hostname and port are correct
4. Try connecting with `psql` first to rule out proxy issues

#### Problem: "authentication failed"

**Possible causes**:
- Wrong username or password
- PostgreSQL server not configured to accept connections from your IP
- User doesn't have access to the specified database

**Solutions**:
1. Verify credentials with `psql`
2. Check PostgreSQL's `pg_hba.conf` for connection permissions
3. Ensure the user has CONNECT privilege on the database

#### Problem: "SSL connection required"

**Possible causes**:
- Server requires SSL but you're using `disable` or `allow` mode

**Solutions**:
1. Use SSL mode `require` or higher
2. Check with your DBA about SSL requirements

### Proxy Issues

#### Problem: "address already in use"

**Cause**: Another process is using port 8080

**Solutions**:
1. Stop the other process using port 8080
2. (Future feature) Use the `--port` flag to specify a different port

To find what's using port 8080:

**Windows**:
```cmd
netstat -ano | findstr :8080
```

**macOS/Linux**:
```bash
lsof -i :8080
# or
sudo netstat -tlnp | grep 8080
```

#### Problem: Browser can't connect to proxy

**Possible causes**:
- Proxy is not running
- Browser is blocking WebSocket connections
- Wrong secret in URL

**Solutions**:
1. Verify proxy is running (check terminal output)
2. Check browser console for errors (F12 → Console)
3. Verify the complete URL including secret
4. Try a different browser
5. Disable browser extensions that might block WebSockets

### Performance Issues

#### Problem: Slow query execution

**Not a proxy issue if**:
- Queries are also slow in `psql` or other clients
- This indicates a database performance issue

**Might be a proxy issue if**:
- Queries are fast in `psql` but slow in the web app
- Check proxy memory and CPU usage
- Check network latency to the database

#### Problem: High memory usage

**Normal**: The proxy maintains a connection pool and caches some data

**Abnormal**: If memory usage exceeds 500MB

**Solutions**:
1. Restart the proxy
2. Report as a bug if it persists

## Frequently Asked Questions

### Is my password secure?

Yes. Your credentials are only stored in the proxy process memory and never sent to any third-party service. The proxy runs entirely on your local machine and connects directly to your PostgreSQL server.

### Do I need to install PostgreSQL locally?

No. You only need the lightweight proxy binary. You can connect to any remote PostgreSQL server.

### Can I connect to multiple databases?

Currently, each proxy instance connects to one database. To connect to multiple databases, run multiple proxy instances on different ports (feature coming soon).

### Does this work with cloud PostgreSQL services?

Yes! It works with:
- AWS RDS PostgreSQL
- Google Cloud SQL for PostgreSQL
- Azure Database for PostgreSQL
- DigitalOcean Managed Databases
- Heroku Postgres
- Supabase
- Neon
- Any other PostgreSQL-compatible service

Just make sure your connection string or credentials are correct and the service allows connections from your IP.

### Can I use this in production?

The proxy is designed for development and administration tasks. While it's secure and stable, consider your use case:

- ✅ **Good for**: Local development, database administration, data exploration
- ⚠️ **Consider carefully**: Production operations, automated scripts
- ❌ **Not recommended for**: Embedding in customer-facing applications

### What data is sent over the network?

- **To PostgreSQL**: Only SQL queries and authentication credentials (via secure connection if SSL enabled)
- **To browser**: Query results via local WebSocket (localhost only)
- **To internet**: Nothing. The web app is static files, and all communication stays local

### How do I update to a new version?

1. Download the new binary from releases
2. Stop the old proxy (Ctrl+C)
3. Replace the old binary with the new one
4. Restart the proxy

Your web app will automatically update if you're using the hosted version, or you can rebuild it locally.

### Can I run this on a server?

While technically possible, this is not the intended use case. The proxy is designed to run on your local machine. Running it on a server would:
- Expose your database credentials
- Create unnecessary security risks
- Defeat the "zero backend infrastructure" design goal

If you need multi-user access, consider traditional database management tools designed for that purpose.

### Does this support PostgreSQL extensions?

Yes! The proxy works with any PostgreSQL extensions. The autocomplete system will also introspect custom functions and types from extensions.

### What about connection pooling?

The proxy maintains a connection pool with:
- Maximum 5 connections
- Minimum 1 connection
- Automatic connection lifecycle management

This provides good performance for typical usage patterns.

## Getting Help

If you encounter issues not covered here:

1. Check the [GitHub Issues](https://github.com/MPJHorner/PostgresMaster/issues) for similar problems
2. Review the [Usage Guide](usage.md) for detailed usage instructions
3. Check the [Architecture Documentation](architecture.md) for technical details
4. Open a new issue with:
   - Your operating system and version
   - Proxy version (`postgres-proxy --version`)
   - Complete error message
   - Steps to reproduce

## Next Steps

- Read the [Usage Guide](usage.md) to learn about features and workflows
- Explore the [Architecture Documentation](architecture.md) to understand how it works
- Check out the [Contributing Guide](../CONTRIBUTING.md) if you want to contribute

---

**Ready to start querying?** Head back to the [Quick Start Guide](../README.md#-quick-start)!
