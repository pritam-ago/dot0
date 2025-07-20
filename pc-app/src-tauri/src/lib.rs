// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
mod server;
use rand;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

#[derive(Serialize, Deserialize)]
struct FileInfo {
    name: String,
    path: String,
    is_directory: bool,
    size: Option<u64>,
    modified: Option<String>,
}

#[tauri::command]
fn generate_pin() -> String {
    let pin: u32 = rand::random_range(100_000..1_000_000);
    pin.to_string()
}

#[tauri::command]
fn start_server(folder: String) -> String {
    let folder_path = std::path::PathBuf::from(&folder);

    std::thread::spawn(move || {
        let rt = tokio::runtime::Runtime::new().unwrap();
        rt.block_on(crate::server::start_file_server(folder_path));
    });

    folder // Return the folder path to the frontend
}

#[tauri::command]
fn list_files(path: String) -> Result<Vec<FileInfo>, String> {
    let path_buf = PathBuf::from(&path);

    if !path_buf.exists() {
        return Err("Path does not exist".to_string());
    }

    let entries =
        fs::read_dir(&path_buf).map_err(|e| format!("Failed to read directory: {}", e))?;

    let mut files = Vec::new();

    for entry in entries {
        if let Ok(entry) = entry {
            let file_type = entry
                .file_type()
                .map_err(|e| format!("Failed to get file type: {}", e))?;

            let metadata = entry
                .metadata()
                .map_err(|e| format!("Failed to get metadata: {}", e))?;

            let name = entry.file_name().to_string_lossy().to_string();

            let path = entry.path().to_string_lossy().to_string();

            let size = if file_type.is_file() {
                Some(metadata.len())
            } else {
                None
            };

            let modified = metadata
                .modified()
                .ok()
                .and_then(|time| time.duration_since(std::time::UNIX_EPOCH).ok())
                .map(|duration| duration.as_secs().to_string());

            files.push(FileInfo {
                name,
                path,
                is_directory: file_type.is_dir(),
                size,
                modified,
            });
        }
    }

    Ok(files)
}

#[tauri::command]
fn read_file(path: String) -> Result<Vec<u8>, String> {
    fs::read(&path).map_err(|e| format!("Failed to read file: {}", e))
}

#[tauri::command]
fn write_file(path: String, content: Vec<u8>) -> Result<(), String> {
    fs::write(&path, content).map_err(|e| format!("Failed to write file: {}", e))
}

#[tauri::command]
fn create_directory(path: String) -> Result<(), String> {
    fs::create_dir_all(&path).map_err(|e| format!("Failed to create directory: {}", e))
}

#[tauri::command]
fn delete_file(path: String) -> Result<(), String> {
    let path_buf = PathBuf::from(&path);

    if path_buf.is_dir() {
        fs::remove_dir_all(&path_buf).map_err(|e| format!("Failed to delete directory: {}", e))
    } else {
        fs::remove_file(&path_buf).map_err(|e| format!("Failed to delete file: {}", e))
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
        .invoke_handler(tauri::generate_handler![
            generate_pin,
            start_server,
            list_files,
            read_file,
            write_file,
            create_directory,
            delete_file
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
