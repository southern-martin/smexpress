---
paths:
  - "frontend/**/*.ts"
  - "frontend/**/*.tsx"
  - "frontend/**/*.css"
---
# Frontend Conventions

## Stack

React 18.3, Vite, TypeScript, Tailwind CSS, TanStack Query, Zustand, React Router 6

## Structure

```
frontend/
├── apps/
│   ├── admin/         # Admin dashboard (port 3000)
│   └── customer/      # Customer portal (port 3001)
└── packages/
    ├── ui/            # @smexpress/ui — shared components (Button, Card, Modal, DataTable, PageHeader)
    ├── auth/          # @smexpress/auth — AuthProvider, ProtectedRoute, useAuth
    ├── api-client/    # @smexpress/api-client — typed API client, TanStack Query hooks
    ├── config/        # @smexpress/config — app configuration
    └── i18n/          # @smexpress/i18n — internationalization
```

## Patterns

- Use shared `@smexpress/ui` components before creating page-specific ones
- API calls go through `@smexpress/api-client` with TanStack Query hooks
- State: server state in TanStack Query, client state in Zustand
- Auth flow: `AuthProvider` -> `ProtectedRoute` (redirect to /login) -> `LoginGuard` (redirect to / if authed)
- Layout: sidebar nav (dark) + top bar with user/logout + `<Outlet/>` content area

## Commands

```bash
cd frontend && pnpm install      # install
cd frontend && pnpm dev          # dev servers
cd frontend && pnpm build        # build all
cd frontend && pnpm test         # test all
```
