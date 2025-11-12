# PostgreSQL Client - Web Application

This is the web application for the browser-based PostgreSQL client.

## Development

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
npm install
```

### Running the Development Server

```bash
npm run dev
```

Open [http://localhost:5173](http://localhost:5173) in your browser.

### Building for Production

```bash
npm run build
```

Preview the production build:

```bash
npm run preview
```

### Testing

Run tests:

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

Generate coverage report:

```bash
npm run test:coverage
```

### Code Quality

Check formatting and linting:

```bash
npm run lint
```

Format code:

```bash
npm run format
```

Type check:

```bash
npm run check
```

## Project Structure

```
web/
├── src/
│   ├── lib/
│   │   ├── components/    # Svelte components
│   │   ├── services/      # WebSocket client, etc.
│   │   ├── stores/        # Svelte stores
│   │   └── utils/         # Utility functions
│   ├── routes/            # SvelteKit routes
│   ├── app.html           # HTML template
│   └── app.d.ts           # TypeScript declarations
├── static/                # Static assets
├── package.json
├── svelte.config.js
├── vite.config.ts
└── tsconfig.json
```

## Tech Stack

- **Framework**: SvelteKit with static adapter
- **Language**: TypeScript
- **Styling**: TailwindCSS + shadcn-svelte
- **Editor**: Monaco Editor
- **Testing**: Vitest
- **Linting**: ESLint + Prettier
