# Vendor App

Flutter mobile application for vendors to manage orders and services.

## Features

- Vendor registration and onboarding
- Dashboard with today's tasks
- Order management
- Workflow status updates
- Rate card management
- Settlement tracking
- Availability management
- Analytics and reports
- Push notifications

## Setup

1. Install Flutter dependencies:
```bash
flutter pub get
```

2. Configure environment:
   - Update API URL in `lib/core/config/env.dart`
   - Add Firebase configuration files:
     - Android: `android/app/google-services.json`
     - iOS: `ios/Runner/GoogleService-Info.plist`
   - Add Razorpay API keys in environment config

3. Run the app:
```bash
flutter run
```

## Build

### Android
```bash
flutter build apk --release
flutter build appbundle --release
```

### iOS
```bash
flutter build ios --release
```

## Architecture

- **Clean Architecture** with feature-first organization
- **Riverpod** for state management
- **Go Router** for navigation
- **Dio** for HTTP requests
- **Repository Pattern** for API calls

## Folder Structure

```
lib/
├── core/           # Core functionality
│   ├── api/        # API client
│   ├── config/     # Configuration
│   ├── models/     # Shared models
│   ├── providers/  # Global providers
│   └── router/     # Navigation
├── features/       # Feature modules
│   ├── auth/
│   ├── dashboard/
│   ├── services/
│   ├── orders/
│   ├── rate_card/
│   ├── settlements/
│   └── ...
└── shared/         # Shared widgets
```

## Code Generation

Run code generation for Freezed and Riverpod:
```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

Watch mode:
```bash
flutter pub run build_runner watch
```

## Testing

Run tests:
```bash
flutter test
```
