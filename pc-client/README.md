# CloudStore PC Client

Desktop client for CloudStore built with Tauri and React.

## 🛠️ Technologies

- Tauri (Rust)
- React + TypeScript
- Vite
- Chakra UI
- React Router

## 🚀 Getting Started

1. **Prerequisites**
   - Node.js 18+
   - Rust
   - Tauri CLI

2. **Installation**
   ```bash
   # Install dependencies
   npm install
   ```

3. **Development**
   ```bash
   # Start development server
   npm run tauri dev
   ```

4. **Build**
   ```bash
   # Build for production
   npm run tauri build
   ```

## 📁 Project Structure

```
pc-client/
├── src/              # React frontend code
├── src-tauri/        # Tauri/Rust backend code
├── public/           # Static assets
└── dist/            # Build output
```

## 🔧 Configuration

The application can be configured through:
- `src-tauri/tauri.conf.json` - Tauri configuration
- `.env` files - Environment variables
- `vite.config.ts` - Vite configuration
