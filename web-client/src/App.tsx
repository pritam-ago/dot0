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
  const fileInputRef = useRef<HTMLInputElement>(null);

  const connect = () => {
    if (pin.length !== 6) {
      alert('Please enter a 6-digit PIN');
      return;
    }

    const relayUrl = `ws://localhost:8080/connect-user/${pin}`;
    const socket = new WebSocket(relayUrl);

    socket.onopen = () => {
      setIsConnected(true);
      requestFileList('/');
    };

    socket.onmessage = (event) => {
      const message: Message = JSON.parse(event.data);
      handleMessage(message);
    };

    socket.onclose = () => {
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
      setIsConnected(false);
    };

    setWs(socket);
  };

  const disconnect = () => {
    if (ws) {
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
          console.log("📤 Sending upload request for:", file.name, "to:", targetPath);
          
          ws.send(JSON.stringify({
            type: 'upload_file',
            data: {
              path: targetPath,
              content: Array.from(new Uint8Array(reader.result as ArrayBuffer)),
              filename: file.name
            }
          }));
        } else {
          console.error("❌ WebSocket not connected, cannot upload file");
          alert("Connection lost. Please reconnect and try again.");
          setIsUploading(false);
          setUploadProgress(0);
        }
        setUploadProgress(((index + 1) / files.length) * 100);
      };
      
      reader.onerror = () => {
        console.error("❌ Failed to read file:", file.name);
        alert(`Failed to read file: ${file.name}`);
        setIsUploading(false);
        setUploadProgress(0);
      };
      
      reader.readAsArrayBuffer(file);
    });
  };

  const handleUploadResponse = (data: any) => {
    console.log("📥 Received upload response:", data);
    if (data.success) {
      console.log("✅ Upload successful for:", data.path);
      setIsUploading(false);
      setUploadProgress(0);
      // Refresh the file list to show the new file
      requestFileList(currentPath);
    } else {
      console.error("❌ Upload failed:", data.error);
      setIsUploading(false);
      setUploadProgress(0);
      alert(`Upload failed: ${data.error || 'Unknown error'}`);
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
        <h1>🌐 Remote File Access</h1>
        <p>Access your files from anywhere in the world</p>
      </header>

      <main className="App-main">
        {!isConnected ? (
          <div className="connection-form">
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
              <p>📁 Drag and drop files here or click to upload</p>
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
                  <span className="file-icon">📁</span>
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
                      {file.is_directory ? '📁' : '📄'}
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
