import { useState } from "react";
import { invoke } from "@tauri-apps/api/core";
import { open } from "@tauri-apps/plugin-dialog";
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

  // Register PIN with relay server
  const registerPin = async (pin: string): Promise<boolean> => {
    try {
      console.log("🔍 Registering PIN with relay server:", pin);
      const response = await fetch('http://localhost:8080/register-pin', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ pin }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        console.error("❌ Failed to register PIN:", response.status, errorText);
        return false;
      }

      const result = await response.json();
      console.log("✅ PIN registered successfully:", result);
      return true;
    } catch (error) {
      console.error("❌ Error registering PIN:", error);
      return false;
    }
  };

  // Connect to relay server
  const connectToRelay = async (pin: string, baseDir: string) => {
    console.log("🔍 Connecting to relay server with PIN:", pin);
    const relayUrl = `ws://localhost:8080/connect-pc/${pin}`;
    const socket = new WebSocket(relayUrl);
    
    socket.onopen = () => {
      console.log("✅ Connected to relay server");
      setIsConnected(true);
      
      // Register base directory
      console.log("📂 Registering base directory:", baseDir);
      socket.send(JSON.stringify({
        type: "register_base_dir",
        data: { path: baseDir }
      }));
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      handleRelayMessage(message);
    };

    socket.onclose = () => {
      console.log("🔌 Disconnected from relay server");
      setIsConnected(false);
    };

    socket.onerror = (error) => {
      console.error("❌ WebSocket error:", error);
    };

    setWs(socket);
  };

  const handleRelayMessage = (message: any) => {
    console.log("📨 Received relay message:", message.type, message.data);
    switch (message.type) {
      case "list_files":
        handleListFiles(message.data.path);
        break;
      case "download_file":
        handleDownloadFile(message.data.path);
        break;
      case "upload_file":
        handleUploadFile(message.data);
        break;
      case "delete_file":
        handleDeleteFile(message.data.path);
        break;
    }
  };

  const handleListFiles = async (path: string) => {
    try {
      const files = await invoke("list_files", { path }) as FileInfo[];
      if (ws) {
        ws.send(JSON.stringify({
          type: "list_files",
          data: { files, path }
        }));
      }
    } catch (error) {
      console.error("Failed to list files:", error);
    }
  };

  const handleDownloadFile = async (filePath: string) => {
    try {
      const content = await invoke("read_file", { path: filePath }) as number[];
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

  const handleUploadFile = async (data: any) => {
    try {
      console.log("📤 Processing upload request for:", data.path);
      await invoke("write_file", { 
        path: data.path, 
        content: data.content 
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
      await invoke("delete_file", { path: filePath });
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
        console.log("📁 Selected folder path:", folder);
        setSelectedFolder(folder);
        
        // Generate PIN and connect to relay
        const newPin = await invoke("generate_pin") as string;
        console.log("🔍 Generated PIN:", newPin);
        setPin(newPin);
        
        // Start local file server
        console.log("🌐 Starting local file server...");
        await invoke('start_server', { folder });
        console.log("✅ Local file server started");
        
        // Register PIN with relay server
        console.log("🔍 Attempting to register PIN with relay server...");
        const registrationSuccess = await registerPin(newPin);
        if (!registrationSuccess) {
          console.error("❌ Failed to register PIN, cannot proceed with connection");
          alert("Failed to register PIN with relay server. Please check if the relay server is running on localhost:8080 and try again.");
          return;
        }
        
        console.log("✅ PIN registered successfully, connecting to relay...");
        // Connect to relay server
        await connectToRelay(newPin, folder);
        
        // Send initial file listing
        if (ws) {
          console.log("📋 Sending initial file listing...");
          ws.send(JSON.stringify({
            type: "list_files",
            data: { path: folder }
          }));
        }
      } else {
        console.log("User cancelled folder selection");
      }
    } catch (error) {
      console.error("❌ Error in selectFolder:", error);
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

  const uploadFile = (file: File, targetPath: string) => {
    if (ws) {
      const reader = new FileReader();
      reader.onload = () => {
        ws.send(JSON.stringify({
          type: "upload_file",
          data: {
            path: targetPath,
            content: reader.result,
            filename: file.name
          }
        }));
      };
      reader.readAsArrayBuffer(file);
    }
  };

  return (
    <main className="container">
      <h1>Remote File Access - PC Controller</h1>

      <div className="status-section">
        <p>Status: {isConnected ? "🟢 Connected" : "🔴 Disconnected"}</p>
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
                  {file.isDirectory ? "📁" : "📄"}
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
