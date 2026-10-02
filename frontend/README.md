# Repono Frontend

React + Vite frontend for Go DFS. This version runs independently in mock/demo mode with browser-local persistence. It is not yet connected to the Go API.

## Run

```bash
npm install
npm run dev
```

Open the local URL printed by Vite. Build with `npm run build`.

## Structure

- `pages/`: route-level screens
- `components/`: reusable UI, dashboard layout, files, sharing and mascot
- `context/`: shared auth and file state
- `services/`: mock service boundary to be replaced/extended for API integration
- `utils/`: validation and formatting
- `styles/`: design tokens, global styles and page/animation styles

Demo authentication is intentionally local and must not be treated as secure authentication. File actions use demo state; backend integration will provide real identity, persistence, authorization and file transfer.
