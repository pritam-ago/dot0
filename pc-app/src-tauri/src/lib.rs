// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
mod server;
use rand;


#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
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

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet, generate_pin, start_server])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
