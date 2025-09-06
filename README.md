# CloudStore - Personal Cloud Storage Platform

Turn your PC into a personal cloud storage platform. Access your files from anywhere, securely.

## 🚀 Features

- Turn any PC into a personal cloud storage server
- Access files from web, desktop, and mobile clients
- Real-time file synchronization
- Secure end-to-end encryption
- Cross-platform support (Windows, macOS, Linux)

## 🛠️ Tech Stack

- **Backend**: Go (Echo/Gin) with WebSocket support
- **PC Client**: Tauri + React
- **Mobile**: React Native
- **Database**: PostgreSQL + Redis
- **Infrastructure**: Docker + Docker Compose
- **CI/CD**: GitHub Actions

## 📁 Project Structure

```
cloudstore/
├── backend/           # Go API server
├── pc-client/         # Tauri desktop client
├── mobile/           # React Native mobile app
├── db/               # Database migrations and schemas
├── infra/            # Infrastructure configuration
├── scripts/          # Development and utility scripts
├── docs/             # Project documentation
└── .github/          # GitHub workflows and templates
```

## 🏁 Quick Start

1. **Prerequisites**

   - Docker and Docker Compose
   - Go 1.21+
   - Node.js 18+
   - Rust (for Tauri)

2. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/cloudstore.git
   cd cloudstore
   ```

3. **Set up environment**

   ```bash
   cp .env.example .env
   ./scripts/setup.sh
   ```

4. **Start the development environment**

   ```bash
   docker-compose up
   ```

5. **Access the applications**
   - Backend API: http://localhost:8080
   - PC Client: Run `cd pc-client && npm run tauri dev`
   - Mobile App: Run `cd mobile && npm run start`

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](docs/CONTRIBUTING.md) for details.

## 📚 Documentation

- [Architecture Overview](docs/ARCHITECTURE.md)
- [Contributing Guidelines](docs/CONTRIBUTING.md)
- [Code of Conduct](docs/CODE_OF_CONDUCT.md)
- [Changelog](docs/CHANGELOG.md)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

If you discover any security-related issues, please email security@cloudstore.example.com instead of using the issue tracker.

## 🙏 Acknowledgments

Thanks to all contributors who help make this project better!
