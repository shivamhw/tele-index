# Tele Search Web Application

A modern React-based web application for searching through the tele index using the provided API endpoint.

## Features

- 🔍 **Advanced Search**: Search through the tele index with full-text search capabilities
- 🎨 **Modern UI**: Beautiful, responsive design with gradient backgrounds and smooth animations
- 📄 **Pagination**: Navigate through search results with pagination controls
- ✨ **Highlighting**: Search terms are highlighted in the results
- 📱 **Responsive**: Works perfectly on desktop, tablet, and mobile devices
- ⚡ **Real-time**: Instant search results with loading states and error handling

## Prerequisites

- Node.js (version 14 or higher)
- npm or yarn
- The tele search API server running on `http://localhost:8095`

## Installation

1. Clone or download this repository
2. Navigate to the project directory:
   ```bash
   cd bl-search-web
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

## Usage

1. Make sure your tele search API server is running on `http://localhost:8095`

2. Start the development server:
   ```bash
   npm start
   ```

3. Open your browser and navigate to `http://localhost:3000`

4. Enter your search query in the search box and click "Search" or press Enter

## API Configuration

The application is configured to call the tele search API with the exact format from your curl command:

```bash
curl 'http://localhost:8095/api/tele/_search' \
  -H 'Content-Type: application/json;charset=UTF-8' \
  --data-raw '{"size":10,"from":0,"explain":true,"highlight":{},"query":{"boost":1,"query":"your search term"},"fields":["*"]}'
```

### API Endpoint Details

- **Base URL**: `http://localhost:8095`
- **Endpoint**: `/api/tele/_search`
- **Method**: POST
- **Content-Type**: `application/json;charset=UTF-8`

### Request Format

```json
{
  "size": 10,
  "from": 0,
  "explain": true,
  "highlight": {},
  "query": {
    "boost": 1,
    "query": "your search term"
  },
  "fields": ["*"]
}
```

## Project Structure

```
src/
├── components/
│   ├── SearchComponent.tsx    # Main search component
│   └── SearchComponent.css    # Component styles
├── services/
│   └── searchService.ts       # API service layer
├── types/
│   └── search.ts             # TypeScript type definitions
├── App.tsx                   # Main app component
├── App.css                   # App-level styles
└── index.tsx                 # Application entry point
```

## Available Scripts

- `npm start` - Runs the app in development mode
- `npm build` - Builds the app for production
- `npm test` - Launches the test runner
- `npm eject` - Ejects from Create React App (not recommended)

## Customization

### Changing the API URL

To change the API base URL, edit `src/services/searchService.ts`:

```typescript
const API_BASE_URL = 'http://your-api-server:port';
```

### Modifying Search Parameters

To modify the search request format, edit the `searchTele` function in `src/services/searchService.ts`.

### Styling

The application uses CSS for styling. Main styles are in:
- `src/components/SearchComponent.css` - Component-specific styles
- `src/App.css` - App-level styles
- `src/index.css` - Global styles

## Troubleshooting

### CORS Issues

If you encounter CORS issues, make sure your API server allows requests from `http://localhost:3000`.

### API Connection Issues

1. Verify the API server is running on `http://localhost:8095`
2. Check the browser's developer console for error messages
3. Ensure the API endpoint `/api/tele/_search` is accessible

### Build Issues

If you encounter build issues:

1. Clear node_modules and reinstall:
   ```bash
   rm -rf node_modules package-lock.json
   npm install
   ```

2. Clear npm cache:
   ```bash
   npm cache clean --force
   ```

## Technologies Used

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Axios** - HTTP client
- **CSS3** - Styling with modern features
- **Create React App** - Build tooling

## License

This project is open source and available under the MIT License.
