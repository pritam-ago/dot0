# Environment Configuration for PC App

This document explains how to configure the PC App to use the deployed relay server at `dot0-go-relay.onrender.com`.

## Environment Variables

The PC App uses the following environment variables:

- `VITE_RELAY_SERVER_URL`: The HTTP URL of your deployed relay server (e.g., `https://dot0-go-relay.onrender.com`)
- `VITE_WS_RELAY_SERVER_URL`: The WebSocket URL of your deployed relay server (e.g., `wss://dot0-go-relay.onrender.com`)

## Current Configuration

The PC app is now configured to connect to the deployed relay server at `dot0-go-relay.onrender.com` by default. If environment variables are not set, it will automatically use:

- HTTP: `https://dot0-go-relay.onrender.com`
- WebSocket: `wss://dot0-go-relay.onrender.com`

## Setup Instructions

1. **Copy the environment file:**

   ```bash
   cp env.local .env.local
   ```

2. **Edit `.env.local` and update the URLs:**

   ```env
   # Replace with your actual deployed URL
   VITE_RELAY_SERVER_URL=https://your-deployed-domain.com
   VITE_WS_RELAY_SERVER_URL=wss://your-deployed-domain.com
   ```

3. **Important Notes:**
   - Use `https://` for HTTP URLs and `wss://` for WebSocket URLs in production
   - Use `http://` for HTTP URLs and `ws://` for WebSocket URLs in development
   - The `.env.local` file is already in `.gitignore` and won't be committed to version control

## Example Configurations

### Current Production Configuration (Default)

```env
VITE_RELAY_SERVER_URL=https://dot0-go-relay.onrender.com
VITE_WS_RELAY_SERVER_URL=wss://dot0-go-relay.onrender.com
```

### Development (Localhost)

```env
VITE_RELAY_SERVER_URL=http://localhost:8080
VITE_WS_RELAY_SERVER_URL=ws://localhost:8080
```

### Custom Production Server

```env
VITE_RELAY_SERVER_URL=https://your-custom-relay-server.com
VITE_WS_RELAY_SERVER_URL=wss://your-custom-relay-server.com
```

## Building for Production

When building the app for production, make sure to:

1. Set the correct environment variables in your deployment platform
2. Ensure your relay server is accessible from the deployed URL
3. Configure CORS on your relay server to allow requests from your deployed domain

## Troubleshooting

- If the app can't connect, check that the relay server is running and accessible at the configured URL
- Ensure your relay server has CORS configured to allow requests from your domain
- For WebSocket connections, make sure your server supports WSS (secure WebSocket) in production
