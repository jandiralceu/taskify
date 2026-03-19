# Taskify Dashboard (Frontend)

The Taskify frontend is a cutting-edge web application built using **Svelte 5** and its innovative **Runes** API. It provides a highly responsive, modern, and interactive interface for managing tasks, users, and teams.

## 🌟 Key Technologies

- **Svelte 5 (Runes)**: Next-generation reactive UI development.
- **TypeScript**: Robust types for components and API integration.
- **TanStack Svelte Query**: Advanced server-state management with caching and optimistic updates.
- **Vite**: Rapid-fast development and build pipeline.
- **Custom CSS**: Premium, handcrafted styles for a unique and high-end aesthetic.

## 🚀 Key Features

- **Dynamic Kanban Board**: Interactive task movement with live status updates.
- **Secure Authentication**: Integrated sign-in, sign-up, and automated token rotation.
- **Task Collaboration**: Create and review comments, attachments, and blockers.
- **User Management**: Profile customization (including avatar uploads) and user role administration.
- **Responsive Design**: Silky-smooth adaptation to desktop and mobile screens.

## 🛠️ Getting Started

1. **Install Dependencies**:
   ```bash
   npm install
   ```
2. **Setup Environment**:
   Ensure you have configured the backend URL (usually `http://localhost:8080/api/v1`).
3. **Run Dev Server**:
   ```bash
   npm run dev
   ```
4. **Build for Production**:
   ```bash
   npm run build
   ```

## 🧪 Testing

- **Unit/Component Tests**: `npm run test:unit` for isolated business logic and component testing via Vitest.
- **End-to-End Tests**: `npm run test:e2e` for complete browser-based workflow validation (Playwright).
