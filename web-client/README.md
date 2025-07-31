# Web Client - Remote File Access

A React-based web client for accessing files remotely through a WebSocket connection to the Go relay server.

## Overview

This web client connects to the deployed Go relay server at `dot0-go-relay.onrender.com` to provide remote file access functionality. Users can browse, download, upload, and delete files on their remote computer through a simple web interface.

## Features

- **Remote File Browsing**: Navigate through directories and view files
- **File Download**: Download files from the remote computer
- **File Upload**: Upload files to the remote computer with drag-and-drop support
- **File Deletion**: Delete files and directories
- **Real-time Connection**: WebSocket-based real-time communication
- **Secure Connection**: Uses WSS (WebSocket Secure) protocol

## Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- A running PC app instance with an active PIN

## Installation

1. Clone the repository
2. Navigate to the web-client directory:
   ```bash
   cd web-client
   ```
3. Install dependencies:
   ```bash
   npm install
   ```

## Development

### `npm start`

Runs the app in development mode on port 3001.
Open [http://localhost:3001](http://localhost:3001) to view it in the browser.

The page will reload if you make edits.
You will also see any lint errors in the console.

### `npm test`

Launches the test runner in interactive watch mode.

### `npm run build`

Builds the app for production to the `build` folder.
The build is optimized and ready for deployment.

## Usage

1. **Start the PC App**: Open the PC application on your computer and select a folder to share
2. **Get the PIN**: Note the 6-digit PIN displayed by the PC app
3. **Connect**: Enter the PIN in the web client and click "Connect"
4. **Browse Files**: Navigate through your files and perform operations

## Connection Details

- **Relay Server**: `dot0-go-relay.onrender.com`
- **Protocol**: WSS (WebSocket Secure)
- **Endpoint**: `/connect-user/{pin}`

## Deployment

The web client can be deployed to any static hosting platform:

- **Render**: Deploy as a static site
- **Vercel**: Connect your repository for automatic deployments
- **Netlify**: Drag and drop the build folder
- **GitHub Pages**: Use GitHub Actions for deployment

No environment variables are required for basic functionality.

## Configuration

The web client is configured to connect directly to the deployed relay server. If you need to use a different relay server, you can modify the `relayHost` variable in `src/App.tsx`.

## Troubleshooting

- **Connection Failed**: Ensure the relay server is running and accessible
- **Invalid PIN**: Verify the PIN is exactly 6 digits
- **File Operations Fail**: Check that the PC app is still running and connected

## Learn More

- [React documentation](https://reactjs.org/)
- [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)
- [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started)
