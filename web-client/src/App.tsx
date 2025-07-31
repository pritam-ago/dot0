import React, { useState, useRef } from 'react';
import './App.css';

interface FileInfo {
  name: string;
  path: string;
  is_directory: boolean;
  size?: number;
  modified?: string;
}

interface Message {
  type: string;
  data: any;
}

function App() {
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [pin, setPin] = useState('');
  const [currentPath, setCurrentPath] = useState('/');
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [isUploading, setIsUploading] = useState(false);
  const [connectionStatus, setConnectionStatus] = useState<string>('Not connected');
  const fileInputRef = useRef<HTMLInputElement>(null);

  const checkPinExists = async (pin: string): Promise<{ exists: boolean; pcConnected?: boolean; error?: string }> => {
    try {
      const response = await fetch(`https://dot0-go-relay.onrender.com/check-pin/${pin}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      
      if (response.ok) {
        const data = await response.json();
        return { 
          exists: data.valid === true, 
          pcConnected: data.pc_connected,
          error: data.error
        };
      }
      return { exists: false, error: 'Failed to check PIN' };
    } catch (error) {
      console.error('Error checking PIN:', error);
      return { exists: false, error: 'Network error' };
    }
  };

  const connect = async () => {   
    if (pin.length !== 6) {
      alert('Please enter a 6-digit PIN');
      return;
    }

    setConnectionStatus('Checking PIN...');
    
    // Check if PIN exists first
    const pinCheck = await checkPinExists(pin);
    if (!pinCheck.exists) {
      setConnectionStatus('PIN not found');
      const errorMsg = pinCheck.error || 'PIN not found';
      alert(`PIN check failed: ${errorMsg}\n\nPlease make sure:\n\n1. The PC app is running\n2. You entered the correct 6-digit PIN\n3. The PC app has registered the PIN`);
      return;
    }

    if (!pinCheck.pcConnected) {
      setConnectionStatus('PC not connected');
      alert('PIN is valid but PC is not connected!\n\nPlease make sure:\n\n1. The PC app is running and connected\n2. The PC app has selected a folder to share\n3. Try refreshing the PC app connection');
      return;
    }

    // Always use WSS (secure WebSocket) for the deployed relay server
    const relayHost = 'dot0-go-relay.onrender.com';
    const relayUrl = `wss://${relayHost}/connect-user/${pin}`;
    
    console.log('Connecting to:', relayUrl);
    
    setConnectionStatus('Connecting...');
    
    const socket = new WebSocket(relayUrl);

    socket.onopen = () => {
      console.log('WebSocket connection established');
      setConnectionStatus('Connected');
      setIsConnected(true);
      requestFileList('/');
    };

    socket.onmessage = (event) => {
      try {
        const message: Message = JSON.parse(event.data);
        handleMessage(message);
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error);
      }
    };

    socket.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
      
      let errorMessage = 'Disconnected';
      switch (event.code) {
        case 1000:
          errorMessage = 'Connection closed normally';
          break;
        case 1006:
          errorMessage = 'Connection failed - PC not connected or PIN invalid';
          break;
        case 1002:
          errorMessage = 'Protocol error';
          break;
        case 1003:
          errorMessage = 'Unsupported data type';
          break;
        case 1005:
          errorMessage = 'No status code provided';
          break;
        case 1007:
          errorMessage = 'Invalid frame payload data';
          break;
        case 1008:
          errorMessage = 'Policy violation';
          break;
        case 1009:
          errorMessage = 'Message too big';
          break;
        case 1010:
          errorMessage = 'Client terminating connection';
          break;
        case 1011:
          errorMessage = 'Server error';
          break;
        case 1012:
          errorMessage = 'Service restart';
          break;
        case 1013:
          errorMessage = 'Try again later';
          break;
        case 1014:
          errorMessage = 'Bad gateway';
          break;
        case 1015:
          errorMessage = 'TLS handshake failed';
          break;
        default:
          errorMessage = `Disconnected (Code: ${event.code})`;
      }
      
      setConnectionStatus(errorMessage);
      setIsConnected(false);
      
      // Show user-friendly error message
      if (event.code === 1006) {
        alert('Connection failed! Please make sure:\n\n1. The PC app is running and connected\n2. You entered the correct 6-digit PIN\n3. The PC app has selected a folder to share');
      }
    };

    socket.onerror = (error) => {
      console.error('WebSocket connection error:', error);
      setConnectionStatus('Connection failed');
      setIsConnected(false);
    };

    setWs(socket);
  };

  const disconnect = () => {
    if (ws) {
      setConnectionStatus('Disconnecting...');
      ws.close();
      setWs(null);
    }
  };

  const requestFileList = (path: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'list_files',
        data: { path }
      }));
    }
  };

  const handleMessage = (message: Message) => {
    switch (message.type) {
      case 'list_files':
        setFiles(message.data.files);
        setCurrentPath(message.data.path);
        break;
      case 'file_content':
        downloadFileFromData(message.data);
        break;
      case 'upload_response':
        handleUploadResponse(message.data);
        break;
    }
  };

  const downloadFileFromData = (data: any) => {
    const blob = new Blob([new Uint8Array(data.content)]);
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = data.filename || 'download';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const downloadFile = (filePath: string) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'download_file',
        data: { path: filePath }
      }));
    }
  };

  const deleteFile = (filePath: string) => {
    if (window.confirm(`Are you sure you want to delete ${filePath}?`)) {
      if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({
          type: 'delete_file',
          data: { path: filePath }
        }));
        // Refresh file list after deletion
        setTimeout(() => requestFileList(currentPath), 100);
      }
    }
  };

  const navigateToFolder = (path: string) => {
    requestFileList(path);
  };

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      uploadFiles(Array.from(files));
    }
  };

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault();
    const files = Array.from(event.dataTransfer.files);
    uploadFiles(files);
  };

  const handleDragOver = (event: React.DragEvent) => {
    event.preventDefault();
  };

  const uploadFiles = (files: File[]) => {
    setIsUploading(true);
    setUploadProgress(0);

    files.forEach((file, index) => {
      const reader = new FileReader();
      reader.onload = () => {
        if (ws && ws.readyState === WebSocket.OPEN) {
          // Create proper path for the target file
          const targetPath = currentPath === '/' ? file.name : `${currentPath}/${file.name}`;
          
          ws.send(JSON.stringify({
            type: 'upload_file',
            data: {
              path: targetPath,
              content: Array.from(new Uint8Array(reader.result as ArrayBuffer)),
              filename: file.name
            }
          }));
        } else {
          alert("Connection lost. Please reconnect and try again.");
          setIsUploading(false);
          setUploadProgress(0);
        }
        setUploadProgress(((index + 1) / files.length) * 100);
      };
      
      reader.onerror = () => {
        alert(`Failed to read file: ${file.name}`);
        setIsUploading(false);
        setUploadProgress(0);
      };
      
      reader.readAsArrayBuffer(file);
    });
  };

  const handleUploadResponse = (data: any) => {
    if (data.success) {
      setIsUploading(false);
      setUploadProgress(0);
      // Refresh the file list to show the new file
      requestFileList(currentPath);
    } else {
      alert(`Upload failed: ${data.error || 'Unknown error'}`);
      setIsUploading(false);
      setUploadProgress(0);
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Remote File Access</h1>
        <p>Access your files from anywhere in the world</p>
      </header>

      <main className="App-main">
        {!isConnected ? (
          <div className="connection-form">
            <div className="instructions">
              <h3>How to Connect:</h3>
              <ol>
                <li>Open the PC app on your computer</li>
                <li>Click "Select Folder to Share"</li>
                <li>Choose a folder and note the PIN</li>
                <li>Make sure the PC app shows "Connected" status</li>
                <li>Enter the PIN below and click Connect</li>
              </ol>
              <div className="connection-status">
                <p><strong>Connection Status:</strong> {connectionStatus}</p>
              </div>
            </div>
            <div className="input-group">
              <input
                type="text"
                value={pin}
                onChange={(e) => setPin(e.target.value)}
                placeholder="Enter 6-digit PIN"
                maxLength={6}
                className="pin-input"
              />
            </div>
            <button onClick={connect} className="connect-btn">
              Connect
            </button>
          </div>
        ) : (
          <div className="file-browser">
            <div className="browser-header">
              <div className="status-indicator">
                <span className={`status-dot ${isConnected ? 'connected' : 'disconnected'}`}></span>
                {isConnected ? 'Connected' : 'Disconnected'}
              </div>
              <button onClick={disconnect} className="disconnect-btn">
                Disconnect
              </button>
            </div>

            <div className="breadcrumb">
              <span>Current Path: {currentPath}</span>
            </div>

            <div 
              className="upload-area"
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onClick={() => fileInputRef.current?.click()}
            >
              <p>Drag and drop files here or click to upload</p>
              <input
                ref={fileInputRef}
                type="file"
                multiple
                onChange={handleFileUpload}
                style={{ display: 'none' }}
              />
            </div>

            {isUploading && (
              <div className="upload-progress">
                <div className="progress-bar">
                  <div 
                    className="progress-fill" 
                    style={{ width: `${uploadProgress}%` }}
                  ></div>
                </div>
                <p>Uploading... {Math.round(uploadProgress)}%</p>
              </div>
            )}

            <div className="file-list">
              {currentPath !== '/' && (
                <div 
                  className="file-item directory"
                  onClick={() => navigateToFolder(currentPath.split('/').slice(0, -1).join('/') || '/')}
                >
                  <span className="file-icon">DIR</span>
                  <span className="file-name">..</span>
                  <span className="file-details">Parent Directory</span>
                </div>
              )}

              {files
                .sort((a, b) => {
                  if (a.is_directory && !b.is_directory) return -1;
                  if (!a.is_directory && b.is_directory) return 1;
                  return a.name.localeCompare(b.name);
                })
                .map((file, index) => (
                  <div 
                    key={index} 
                    className={`file-item ${file.is_directory ? 'directory' : 'file'}`}
                    onClick={() => file.is_directory ? navigateToFolder(file.path) : null}
                  >
                    <span className="file-icon">
                      {file.is_directory ? 'DIR' : 'FILE'}
                    </span>
                    <span className="file-name">{file.name}</span>
                    <span className="file-details">
                      {file.size ? formatFileSize(file.size) : ''}
                    </span>
                    <div className="file-actions">
                      {!file.is_directory && (
                        <button 
                          onClick={(e) => {
                            e.stopPropagation();
                            downloadFile(file.path);
                          }}
                          className="action-btn download"
                        >
                          Download
                        </button>
                      )}
                      <button 
                        onClick={(e) => {
                          e.stopPropagation();
                          deleteFile(file.path);
                        }}
                        className="action-btn delete"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
