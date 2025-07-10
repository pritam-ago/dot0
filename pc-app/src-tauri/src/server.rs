use warp;
use std::path::PathBuf;

pub async fn start_file_server(folder: PathBuf) {
    let routes = warp::fs::dir(folder);

    println!("🌐 Serving on http://localhost:3000");
    warp::serve(routes)
        .run(([127, 0, 0, 1], 3000))
        .await;
}
