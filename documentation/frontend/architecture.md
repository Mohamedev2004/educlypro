# Frontend Architecture Documentation

## Vision & Core Principles

The frontend of EduclyPro is built with a focus on **Maintainability**, **Scalability**, and **Developer Velocity**. We achieve this by adhering to a strict **Layered Architecture** and the **Single Responsibility Principle (SRP)**.

### 1. Single Responsibility Principle (SRP)
Each file must have exactly ONE reason to change. If a file's purpose requires the word "and" to describe (e.g., "this hook fetches data **and** formats it"), it should be split.
- **Components**: Responsible only for UI (Props in, JSX out).
- **Hooks**: Responsible only for state management or a single side effect.
- **Services**: Responsible only for API communication.
- **Utils**: Responsible only for pure data transformation.

### 2. Layered Architecture
We enforce a strict separation of concerns through predefined layers. Dependencies should only flow "downwards" or "sideways" within specific rules (e.g., UI depends on Hooks, Hooks depend on Services).

---

## Directory Structure & Layer Responsibilities

### `src/api/`
- **`services/`**: Contains classes or objects that perform raw API calls (Axios/Fetch). They return data directly and do not manage React state.
- **`types/`**: API-specific TypeScript interfaces (Request/Response shapes).
- **`axios.ts`**: Global configuration for the Axios instance (interceptors, base URL).

### `src/components/`
- **`ui/`**: Low-level, reusable primitive components (Buttons, Inputs, Modals). Usually built on top of Radix UI or Shadcn.
- **`[domain]/`**: Feature-specific UI components (e.g., `logs/LogRow.tsx`). These are pure UI and receive data/handlers via props.
- **Global Components**: Shared layout components like `AppHeader`, `AppSidebar`, etc.

### `src/hooks/`
- Custom React hooks that encapsulate complex state logic or side effects.
- **Domain-specific hooks**: (e.g., `hooks/notifications/use-mark-read.ts`) handle one specific operation.
- **Shared hooks**: (e.g., `use-mobile.ts`) provide reusable utility state.

### `src/pages/`
- These are the entry points for routes.
- **Responsibility**: Orchestrate domain components and hooks. They are the "glue" that connects the UI to the data fetching logic.

### `src/utils/`
- Pure helper functions with NO side effects and NO state.
- They must be deterministic (same input = same output).
- Categories include `ui-utils.ts` (styling), `error-utils.ts` (error normalization), and `session-utils.ts`.

### `src/types/`
- Global UI-specific TypeScript definitions that aren't strictly tied to API responses.

### `src/constants/`
- Static, immutable values used across the app (e.g., `PER_PAGE_OPTIONS`, `ROLE_PERMISSIONS`).

---

## Developer Guidelines

### 🚫 Prohibited Actions
- **No API calls in Components**: Move them to a Service and call them via a Hook or React Query.
- **No Business Logic in Components**: Extract logic to a Util (if pure) or a Hook (if stateful).
- **No Direct State Manipulation in Services**: Services only return data; they don't know about React.
- **No Cross-Concern Mixing**: Don't put formatting logic inside a data-fetching hook.

### ✅ Best Practices
1. **File Naming**: Use kebab-case for files (e.g., `use-user-profile.ts`).
2. **Component Structure**:
   ```tsx
   // 1. Imports (External -> Internal)
   // 2. Types/Interfaces
   // 3. Component Implementation
   // 4. Sub-components (if small)
   ```
3. **State Management**:
   - Use **React Query** for server state (caching, loading states).
   - Use **Context API** for global UI state (auth, theme, language).
   - Use **Local State (`useState`)** for component-specific UI toggles.

### 🔄 Data Flow
1. **User Interaction**: User clicks a button in a `Component`.
2. **Hook Execution**: The component calls a function from a `Hook`.
3. **Service Call**: The hook triggers a `Service` method.
4. **API Response**: The service returns data from the backend.
5. **State Update**: The hook updates the `React Query` cache or local state.
6. **UI Re-render**: The component receives the new data via props or hook return value and updates.

---

## Architectural Vision
By following these rules, we ensure that:
1. **Testing is easy**: Utils and Services can be tested in isolation.
2. **Refactoring is safe**: Changing an API endpoint only requires updating the Service and its Type.
3. **Onboarding is fast**: New developers know exactly where to find (and where to put) any piece of code.
