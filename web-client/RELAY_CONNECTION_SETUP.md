# Web Client - Relay Server Connection Setup

This document explains how the web client connects to the deployed relay server.

## Current Configuration

The web client is now configured to directly connect to the deployed relay server at `dot0-go-relay.onrender.com` using WSS (secure WebSocket) protocol.

## Connection Details

- **Relay Server**: `dot0-go-relay.onrender.com`
- **Protocol**: WSS (WebSocket Secure)
- **Endpoint**: `/connect-user/{pin}`
- **Full URL**: `wss://dot0-go-relay.onrender.com/connect-user/{pin}`

## Key Features

1. **Simplified Configuration**: No environment variables needed for basic connection
2. **Secure Connection**: Always uses WSS protocol for production
3. **Direct Connection**: Connects directly to the deployed relay server
4. **Error Handling**: Improved error handling and logging

## Code Implementation

The connection is established in `App.tsx`:

```javascript
const connect = () => {
  if (pin.length !== 6) {
    alert("Please enter a 6-digit PIN");
    return;
  }

  // Always use WSS (secure WebSocket) for the deployed relay server
  const relayHost = "dot0-go-relay.onrender.com";
  const relayUrl = `wss://${relayHost}/connect-user/${pin}`;

  console.log("Connecting to:", relayUrl);

  const socket = new WebSocket(relayUrl);
  // ... connection handling
};
```

## Environment Variables

The web client no longer requires environment variables for basic functionality. However, if you need to override the relay server URL, you can set:

```env
REACT_APP_RELAY_HOST=your-custom-relay-server.com
```

## Deployment

When deploying to any platform (Render, Vercel, Netlify, etc.):

1. **No Environment Variables Required**: The web client will work out of the box
2. **Automatic WSS**: The connection automatically uses secure WebSocket
3. **CORS Handling**: The relay server handles CORS for WebSocket connections

## Testing

To test the connection:

1. Deploy the web client to your preferred platform
2. Open the PC app and generate a PIN
3. Enter the PIN in the web client
4. Click "Connect" to establish the WebSocket connection

## Troubleshooting

If connection fails:

1. **Check Relay Server**: Ensure `dot0-go-relay.onrender.com` is accessible
2. **Check PIN**: Verify the PIN is exactly 6 digits
3. **Browser Console**: Check for WebSocket connection errors
4. **Network**: Ensure no firewall is blocking WSS connections

## Development

For local development with a local relay server, you can modify the `relayHost` variable in `App.tsx`:

```javascript
const relayHost = "localhost:8080"; // For local development
const relayUrl = `ws://${relayHost}/connect-user/${pin}`; // Use WS for local
```
