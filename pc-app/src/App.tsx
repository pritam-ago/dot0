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

  // Connect to relay server
  const connectToRelay = async (pin: string, baseDir: string) => {
    const relayUrl = `ws://localhost:8080/connect-pc/${pin}`;
    const socket = new WebSocket(relayUrl);
    
    socket.onopen = () => {
      console.log("Connected to relay server");
      setIsConnected(true);
      
      // Register base directory
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
      console.log("Disconnected from relay server");
      setIsConnected(false);
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
    const folder = await open({
      directory: true,
      multiple: false,
    });
  
    if (folder) {
      console.log("Selected folder path:", folder);
      setSelectedFolder(folder);
      
      // Generate PIN and connect to relay
      const newPin = await invoke("generate_pin") as string;
      setPin(newPin);
      
      // Start local file server
      await invoke('start_server', { folder });
      
      // Connect to relay server
      await connectToRelay(newPin, folder);
      
      // Send initial file listing
      if (ws) {
        ws.send(JSON.stringify({
          type: "list_files",
          data: { path: folder }
        }));
      }
    } else {
      console.log("User cancelled folder selection");
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
