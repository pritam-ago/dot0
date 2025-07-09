// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
use rand::{Rng, rng};

#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn generate_pin() -> String {
    let pin: u32 = rng().gen_range(100_000..1_000_000);
    pin.to_string()
}


#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![greet, generate_pin])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
