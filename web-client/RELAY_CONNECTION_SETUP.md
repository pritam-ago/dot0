# Web Client - Relay Server Connection Setup

This document explains how to configure the web client to connect to the deployed relay server.

## The Problem

The web client was trying to connect to `wss://dot0-web-client.onrender.com:8080/connect-user/302597`, which is incorrect because:

- It's connecting to the web client's own domain instead of the relay server
- It's using port 8080 with WSS (secure WebSocket), which is invalid

## The Solution

The web client should connect to the relay server domain: `wss://dot0-go-relay.onrender.com/connect-user/302597`

## Environment Configuration

### Production Environment (env.production)

```env
REACT_APP_WEBSOCKET_HOST=dot0-go-relay.onrender.com
REACT_APP_WEBSOCKET_PORT=
```

### Local Development (env.local)

```env
REACT_APP_WEBSOCKET_HOST=dot0-go-relay.onrender.com
REACT_APP_WEBSOCKET_PORT=
```

## Key Changes Made

1. **Updated WebSocket Host**: Changed from `your-app-name.onrender.com` to `dot0-go-relay.onrender.com`
2. **Removed Port**: Set `REACT_APP_WEBSOCKET_PORT=` (empty) for HTTPS/WSS connections
3. **Updated App.tsx**: Modified URL construction to handle empty port values

## URL Construction Logic

The web client now builds the WebSocket URL as follows:

```javascript
const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
const wsHost = process.env.REACT_APP_WEBSOCKET_HOST || window.location.hostname;
const wsPort = process.env.REACT_APP_WEBSOCKET_PORT || "8080";

// Build relay URL - don't include port if it's empty (for HTTPS/WSS)
const relayUrl = wsPort
  ? `${protocol}//${wsHost}:${wsPort}/connect-user/${pin}`
  : `${protocol}//${wsHost}/connect-user/${pin}`;
```

## Result

- **Development**: `ws://localhost:8080/connect-user/{pin}` (if using local relay)
- **Production**: `wss://dot0-go-relay.onrender.com/connect-user/{pin}` (no port specified)

## Deployment Notes

When deploying to Render, make sure to:

1. Set the environment variables in the Render dashboard:

   - `REACT_APP_WEBSOCKET_HOST=dot0-go-relay.onrender.com`
   - `REACT_APP_WEBSOCKET_PORT=` (leave empty)

2. Ensure your relay server (`dot0-go-relay.onrender.com`) is running and accessible

3. Verify that the relay server accepts WebSocket connections on the `/connect-user/{pin}` endpoint
