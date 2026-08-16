# EduclyPro Frontend

## Architecture

This project follows a strict **Layered Architecture** and **Single Responsibility Principle (SRP)**.

For a detailed explanation of our directory structure, layer responsibilities, and developer guidelines, please refer to the [Architecture Documentation](./ARCHITECTURE.md).

## Getting Started

### Prerequisites
- Node.js (Latest LTS)
- npm or yarn

### Installation
1. Clone the repository.
2. Navigate to the `frontend` directory.
3. Install dependencies:
   ```bash
   npm install
   ```
4. Start the development server:
   ```bash
   npm run dev
   ```

## Component Management

### Adding UI Components
We use **shadcn/ui** for primitive components. To add a new component:
```bash
npx shadcn@latest add [component-name]
```
This will place the UI components in the `src/components/ui` directory.

### Importing Components
Use absolute paths for imports:
```tsx
import { Button } from "@/components/ui/button"
```
