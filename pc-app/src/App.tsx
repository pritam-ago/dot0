import { useState, useEffect } from "react";
import { invoke } from "@tauri-apps/api/core";
import { open } from "@tauri-apps/plugin-dialog";
import { join } from '@tauri-apps/api/path';
import "./App.css";

interface FileInfo {
  name: string;
  path: string;
  isDirectory: boolean;
  size?: number;
  modified?: string;
}

function App() {
  const [pin, setPin] = useState("");
  const [selectedFolder, setSelectedFolder] = useState("");
  const [isConnected, setIsConnected] = useState(false);
  const [currentPath, setCurrentPath] = useState("");
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);
  // Store current session data for immediate access
  const [sessionData, setSessionData] = useState<{pin: string, folder: string} | null>(null);

  // Helper functions for localStorage
  const saveToLocalStorage = (key: string, value: string) => {
    try {
      localStorage.setItem(key, value);
    } catch (error) {
      console.error(`Failed to save to localStorage: ${key}`, error);
    }
  };

  const getFromLocalStorage = (key: string): string | null => {
    try {
      const value = localStorage.getItem(key);
      return value;
    } catch (error) {
      console.error(`Failed to get from localStorage: ${key}`, error);
      return null;
    }
  };

  // Load saved data on app startup
  const loadSavedData = () => {
    const savedPin = getFromLocalStorage('pc_app_pin');
    const savedFolder = getFromLocalStorage('pc_app_folder');
    
    if (savedPin && savedFolder) {
      setPin(savedPin);
      setSelectedFolder(savedFolder);
      setSessionData({ pin: savedPin, folder: savedFolder });
    }
  };

  // Load saved data when component mounts
  useEffect(() => {
    loadSavedData();
  }, []);

  // Register PIN with relay server
  const registerPin = async (pin: string): Promise<boolean> => {
    try {
      // Use environment variable or fallback to deployed relay server
      const relayUrl = import.meta.env.VITE_RELAY_SERVER_URL || 'https://dot0-go-relay.onrender.com';
      console.log('Registering PIN with relay server:', relayUrl);
      
      const response = await fetch(`${relayUrl}/register-pin`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ pin }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("Failed to register PIN:", response.status, errorText);
        return false;
      }

      const result = await response.json();
      console.log('PIN registered successfully:', result);
      return true;
    } catch (error) {
      console.error("Error registering PIN:", error);
      return false;
    }
  };

  // Connect to relay server
  const connectToRelay = async (pin: string, baseDir: string) => {
    // Use environment variable or fallback to deployed relay server
    const wsRelayUrl = import.meta.env.VITE_WS_RELAY_SERVER_URL || 'wss://dot0-go-relay.onrender.com';
    const relayUrl = `${wsRelayUrl}/connect-pc/${pin}`;
    console.log('Connecting to relay server:', relayUrl);
    const socket = new WebSocket(relayUrl);
    
    socket.onopen = () => {
      console.log('WebSocket connection established with relay server');
      setIsConnected(true);
      
      // Register base directory
      socket.send(JSON.stringify({
        type: "register_base_dir",
        data: { path: baseDir }
      }));
      
      // Send a test message to verify connection
      socket.send(JSON.stringify({
        type: "test_message",
        data: { message: "PC app is ready to receive messages" }
      }));
    };

    socket.onmessage = (event) => {
      
      try {
        const message = JSON.parse(event.data);
        handleRelayMessage(message);
      } catch (error) {
        console.error("Failed to parse message from relay:", error);
      }
    };

    socket.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.error("WebSocket connection error:", error);
      setIsConnected(false);
    };

    setWs(socket);
  };

  const handleRelayMessage = (message: any) => {
    
    switch (message.type) {
      case "list_files":
        handleListFiles(message.data.path);
        break;
      case "download_file":
        handleDownloadFile(message.data.path);
        break;
      case "upload_file":
        // Use localStorage for immediate access, then session data, then state
        const storedPin = getFromLocalStorage('pc_app_pin');
        const storedFolder = getFromLocalStorage('pc_app_folder');
        const currentPin = storedPin || sessionData?.pin || pin;
        const currentFolder = storedFolder || sessionData?.folder || selectedFolder;
        
        if (!currentPin || !currentFolder) {
          if (ws) {
            ws.send(JSON.stringify({
              type: "upload_response",
              data: { 
                success: false, 
                error: "App not properly initialized. Please select a folder first." 
              }
            }));
          }
          return;
        }
        handleUploadFile(message.data, currentPin, currentFolder);
        break;
      case "delete_file":
        handleDeleteFile(message.data.path);
        break;
      case "test_message":
        break;
      default:
        break;
    }
  };

  const handleListFiles = async (path: string) => {
    try {
      // For list files, use the selectedFolder directly if path is empty or just "/"
      const fullPath = (path === "/" || path === "" || path === selectedFolder) ? selectedFolder : await join(selectedFolder, path);
      const files = await invoke("list_files", { path: fullPath }) as FileInfo[];
      
      // Update local file list
      setFiles(files);
      
      if (ws) {
        ws.send(JSON.stringify({
          type: "list_files",
          data: { files, path }
        }));
      }
      
      // Update current path if this was a navigation request
      if (path !== currentPath) {
        setCurrentPath(path);
      }
    } catch (error) {
      console.error("Failed to list files:", error);
    }
  };

  const handleDownloadFile = async (filePath: string) => {
    try {
      const fullPath = selectedFolder ? await join(selectedFolder, filePath) : filePath;
      const content = await invoke("read_file", { path: fullPath }) as number[];
      if (ws) {
        ws.send(JSON.stringify({
          type: "file_content",
          data: { 
            path: filePath, 
            content: content,
            filename: filePath.split('/').pop() || filePath
          }
        }));
      }
    } catch (error) {
      console.error("Failed to read file:", error);
    }
  };

  const handleUploadFile = async (data: any, currentPin?: string, currentFolder?: string) => {
    try {
      
      // Always get base directory from database
      let baseDir = null;
      
      // Use passed parameters or fallback to state
      const pinToUse = currentPin || pin;
      const folderToUse = currentFolder || selectedFolder;
      
      // Try to get PIN from state, or use selectedFolder as fallback
      if (pinToUse) {
        try {
          const relayUrl = import.meta.env.VITE_RELAY_SERVER_URL || 'https://dot0-go-relay.onrender.com';
          const response = await fetch(`${relayUrl}/get-base-dir/${pinToUse}`);
          if (response.ok) {
            const result = await response.json();
            if (result.base_directory) {
              baseDir = result.base_directory;
            } else {
              console.error("No base directory found in database for PIN:", pinToUse);
            }
          } else {
            console.error("Failed to fetch base directory from database:", response.status);
          }
        } catch (error) {
          console.error("Error fetching base directory from database:", error);
        }
      } else {
        if (folderToUse) {
          baseDir = folderToUse;
        } else {
          console.error("No PIN available and no selectedFolder fallback");
          throw new Error("No PIN available for fetching base directory");
        }
      }
      
      // Convert content to Uint8Array if it's an array of numbers
      let content: Uint8Array;
      if (Array.isArray(data.content)) {
        content = new Uint8Array(data.content);
      } else if (data.content instanceof ArrayBuffer) {
        content = new Uint8Array(data.content);
      } else {
        throw new Error("Invalid content format");
      }
      
      let fullPath;
      if (baseDir) {
        fullPath = await join(baseDir, data.path);
      } else {
        fullPath = data.path;
      }
      
      // Ensure we have an absolute path
      let absolutePath;
      if (fullPath.startsWith('/') || fullPath.includes(':')) {
        absolutePath = fullPath;
      } else {
        // If it's not absolute, we need to construct the proper path
        if (baseDir) {
          // Use the base directory as the base
          absolutePath = `${baseDir}\\${fullPath}`;
        } else {
          // Fallback to E: drive
          absolutePath = `E:\\${fullPath}`;
        }
      }
      
      await invoke("write_file", { 
        path: absolutePath, 
        content: Array.from(content) // Convert back to array for Tauri
      });
      
      if (ws) {
        ws.send(JSON.stringify({
          type: "upload_response",
          data: { success: true, path: data.path }
        }));
      }
      
      // Refresh file list
      handleListFiles(currentPath);
    } catch (error) {
      console.error("Failed to write file:", error);
      if (ws) {
        ws.send(JSON.stringify({
          type: "upload_response",
          data: { success: false, error: String(error) }
        }));
      }
    }
  };

  const handleDeleteFile = async (filePath: string) => {
    try {
      const fullPath = selectedFolder ? await join(selectedFolder, filePath) : filePath;
      await invoke("delete_file", { path: fullPath });
      // Refresh file list
      handleListFiles(currentPath);
    } catch (error) {
      console.error("Failed to delete file:", error);
    }
  };

  async function selectFolder() {
    try {
      const folder = await open({
        directory: true,
        multiple: false,
      });
    
      if (folder) {
        setSelectedFolder(folder);
        
        // Generate PIN and connect to relay
        const newPin = await invoke("generate_pin") as string;
        setPin(newPin);
        
        // Save to localStorage for immediate access
        saveToLocalStorage('pc_app_pin', newPin);
        saveToLocalStorage('pc_app_folder', folder);
        
        // Set session data for immediate access
        setSessionData({ pin: newPin, folder });
        
        // Start local file server
        await invoke('start_server', { folder });
        
        // Register PIN with relay server
        const registrationSuccess = await registerPin(newPin);
        if (!registrationSuccess) {
          console.error("Failed to register PIN, cannot proceed with connection");
          const relayUrl = import.meta.env.VITE_RELAY_SERVER_URL || 'https://dot0-go-relay.onrender.com';
          alert(`Failed to register PIN with relay server at ${relayUrl}. Please check if the relay server is running and try again.`);
          return;
        }
        
        // Connect to relay server
        await connectToRelay(newPin, folder);
        
        // Wait a moment for WebSocket connection to be fully established
        await new Promise(resolve => setTimeout(resolve, 100));
        
        // Set initial path and send initial file listing
        setCurrentPath(folder);
        
        if (ws) {
          ws.send(JSON.stringify({
            type: "list_files",
            data: { path: folder }
          }));
        }
      } else {
      }
    } catch (error) {
      console.error("Error in selectFolder:", error);
      alert("An error occurred while setting up the connection. Please try again.");
    }
  }

  const navigateToFolder = (path: string) => {
    if (ws) {
      ws.send(JSON.stringify({
        type: "list_files",
        data: { path }
      }));
    }
  };

  const downloadFile = (filePath: string) => {
    if (ws) {
      ws.send(JSON.stringify({
        type: "download_file",
        data: { path: filePath }
      }));
    }
  };

  return (
    <main className="container">
      <h1>Remote File Access - PC Controller</h1>

      <div className="status-section">
        <p>Status: {isConnected ? "Connected" : "Disconnected"}</p>
        {pin && <p>PIN: <strong>{pin}</strong></p>}
        {selectedFolder && <p>Sharing: {selectedFolder}</p>}
      </div>

      <div className="controls">
        <button onClick={selectFolder} className="btn primary">
          Select Folder to Share
        </button>
      </div>

      {isConnected && (
        <div className="file-browser">
          <h2>Current Path: {currentPath}</h2>
          <div className="file-list">
            {files.map((file, index) => (
              <div key={index} className="file-item">
                <span className="file-icon">
                  {file.isDirectory ? "DIR" : "FILE"}
                </span>
                <span 
                  className="file-name"
                  onClick={() => file.isDirectory ? navigateToFolder(file.path) : downloadFile(file.path)}
                >
                  {file.name}
                </span>
                {!file.isDirectory && file.size && (
                  <span className="file-size">({Math.round(file.size / 1024)}KB)</span>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </main>
  );
}

export default App;
