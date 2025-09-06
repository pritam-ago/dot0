# CloudStore Mobile App

Mobile client for CloudStore built with React Native.

## 🛠️ Technologies

- React Native
- TypeScript
- React Navigation
- React Native Paper
- AsyncStorage

## 🚀 Getting Started

1. **Prerequisites**
   - Node.js 18+
   - React Native development environment
   - Android Studio / Xcode

2. **Installation**
   ```bash
   # Install dependencies
   npm install
   ```

3. **Development**
   ```bash
   # Start Metro bundler
   npm start

   # Run on Android
   npm run android

   # Run on iOS
   npm run ios
   ```

## 📁 Project Structure

```
mobile/
├── src/
│   ├── screens/     # Screen components
│   ├── components/  # Reusable components
│   ├── services/   # API and business logic
│   └── utils/      # Helper functions
├── android/        # Android native code
└── ios/           # iOS native code
```

## 🔧 Configuration

The application can be configured through:
- `.env` files - Environment variables
- `app.json` - React Native configuration
- `metro.config.js` - Metro bundler configuration
