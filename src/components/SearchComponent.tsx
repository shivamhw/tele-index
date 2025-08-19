import React, { useMemo, useState } from 'react';
import { searchTele } from '../services/searchService';
import { SearchResult } from '../types/search';
import './SearchComponent.css';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { searchTMDB, TMDBSuggestion } from '../services/tmdbService';

const DOWNLOAD_BASE_URL = process.env.REACT_APP_DOWNLOAD_BASE_URL || 'http://localhost:8080';

// Helper to format bytes to MB/GB
function formatBytes(bytes: number): string {
  if (bytes >= 1024 * 1024 * 1024) {
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
  } else if (bytes >= 1024 * 1024) {
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  } else if (bytes >= 1024) {
    return (bytes / 1024).toFixed(2) + ' KB';
  } else {
    return bytes + ' B';
  }
}

const SearchComponent: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialQuery = searchParams.get('q') || '';
  const initialPage = parseInt(searchParams.get('page') || '1', 10);
  const [query, setQuery] = useState(initialQuery);
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const [suggestions, setSuggestions] = useState<TMDBSuggestion[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selected, setSelected] = useState<{
    mediaType: 'movie' | 'tv' | null;
    title: string | null;
    year?: number;
    posterPath?: string;
  }>({ mediaType: null, title: null, year: undefined, posterPath: undefined });
  const [movieYear, setMovieYear] = useState<number | ''>('');
  const [tvSeason, setTvSeason] = useState<number | ''>('');
  const [tvEpisode, setTvEpisode] = useState<number | ''>('');
  const [hasSearched, setHasSearched] = useState(false);
  const [lastSearchedQuery, setLastSearchedQuery] = useState('');
  const navigate = useNavigate();

  // Do not auto-search on mount; user must click Search
  React.useEffect(() => {
    // Keep query string in input if provided via URL, but don't trigger search automatically
    // eslint-disable-next-line
  }, []);

  // Debounced TMDB autocomplete
  React.useEffect(() => {
    const q = query.trim();
    if (!q || q.length < 2) {
      setSuggestions([]);
      return;
    }
    const id = setTimeout(async () => {
      const res = await searchTMDB(q);
      setSuggestions(res);
      setShowSuggestions(true);
    }, 300);
    return () => clearTimeout(id);
  }, [query]);

  const pad2 = (n: number) => String(n).padStart(2, '0');

  const composedQuery = useMemo(() => {
    if (selected.mediaType === 'movie' && selected.title) {
      const year = movieYear || selected.year;
      return year ? `${selected.title} ${year}` : `${selected.title}`;
    }
    if (selected.mediaType === 'tv' && selected.title) {
      const s = typeof tvSeason === 'number' ? tvSeason : Number(tvSeason);
      const e = typeof tvEpisode === 'number' ? tvEpisode : Number(tvEpisode);
      if (s && e) return `${selected.title} s${pad2(s)}e${pad2(e)}`;
      if (s) return `${selected.title} s${pad2(s)}`;
      return selected.title;
    }
    return query;
  }, [selected, movieYear, tvSeason, tvEpisode, query]);

  const posterUrl = useMemo(() => {
    if (!selected.posterPath) return '';
    return `https://image.tmdb.org/t/p/w154${selected.posterPath}`;
  }, [selected.posterPath]);

  const handleSearch = async (searchQuery: string = composedQuery || query, page: number = 1) => {
    if (!searchQuery.trim()) {
      setError('Please enter a search query');
      setResults([]);
      setCurrentPage(1);
      setSearchParams({});
      return;
    }
    setLoading(true);
    setError(null);
    setHasSearched(true);
    setLastSearchedQuery(searchQuery);
    setSearchParams({ q: searchQuery, page: String(page) });
    try {
      const searchResults = await searchTele(searchQuery, page);
      setResults(searchResults);
      setCurrentPage(page);
    } catch (err) {
      setError('Failed to fetch search results. Please try again.');
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    handleSearch(composedQuery || query, 1);
  };

  const handlePageChange = (newPage: number) => {
    handleSearch(lastSearchedQuery || composedQuery || query, newPage);
  };

  // Updated renderResult for your API's hit structure
  const renderResult = (result: SearchResult) => {
    const fromID = result.source.ChatId;
    const fileID = result.source.ID;
    const downloadUrl =
      fromID && fileID ? `${DOWNLOAD_BASE_URL}/${fromID}/${fileID}` : null;
    console.log(downloadUrl);
    return (
      <div key={result.id} className="search-result">
        <div className="result-header">
          <h3 className="result-title">
            {result.source.Tokens || result.source.File || result.id}
          </h3>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <span className="result-score">Score: {result.score.toFixed(2)}</span>
            {downloadUrl && (
              <>
              <a
                href={downloadUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="download-button"
              > 
                Download
              </a>
              
              <button
                className="download-button"
                onClick={() => navigate(`/stream?url=${encodeURIComponent(downloadUrl)}`)}
                style={{ cursor: 'pointer' }}
              >
                Stream
              </button>
              </>
            )}
          </div>
        </div>
        <div className="result-content">
          {/* Show highlight if available, else show Tokens */}
          {result.highlight?.Tokens && result.highlight.Tokens.length > 0 ? (
            result.highlight.Tokens.map((frag, idx) => (
              <div key={idx} dangerouslySetInnerHTML={{ __html: frag }} />
            ))
          ) : (
            <div>{result.source.Tokens}</div>
          )}
        </div>
        <div className="result-meta">
          <span className="result-id">ID: {result.id}</span>
          {result.source.Size && (
            <span className="result-size">Size: {formatBytes(Number(result.source.Size))}</span>
          )}
          {result.source.File && (
            <span className="result-file">File: {result.source.File}</span>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="search-container">
      <div className="search-header">
        <h1>
          <a
            href="/"
            className="brand-link"
            onClick={(e) => {
              e.preventDefault();
              setQuery('');
              setResults([]);
              setHasSearched(false);
              setLastSearchedQuery('');
              setError(null);
              setCurrentPage(1);
              setSelected({ mediaType: null, title: null, year: undefined, posterPath: undefined });
              setMovieYear('');
              setTvSeason('');
              setTvEpisode('');
              setSuggestions([]);
              setShowSuggestions(false);
              setSearchParams({});
              navigate('/', { replace: true });
            }}
          >
            Tele Search
          </a>
        </h1>
        <p>Search through the tele index with advanced querying</p>
      </div>

      <form onSubmit={handleSubmit} className="search-form" autoComplete="off">
        <div className="search-input-group">
          <input
            type="text"
            value={query}
            onChange={(e) => {
              setQuery(e.target.value);
              setSelected({ mediaType: null, title: null, year: undefined, posterPath: undefined });
              setMovieYear('');
              setTvSeason('');
              setTvEpisode('');
              setHasSearched(false);
            }}
            onFocus={() => suggestions.length && setShowSuggestions(true)}
            placeholder="Search movies or series (e.g., 'Wanted', 'Game of Thrones')"
            className="search-input"
            disabled={loading}
          />
          <button type="submit" className="search-button" disabled={loading}>
            {loading ? 'Searching...' : 'Search'}
          </button>
        </div>

        {showSuggestions && suggestions.length > 0 && (
          <div className="suggestions" onMouseLeave={() => setShowSuggestions(false)}>
            {suggestions.slice(0, 8).map((sug) => (
              <div
                key={`${sug.mediaType}-${sug.id}`}
                className="suggestion-item"
                onClick={() => {
                  setQuery(sug.title);
                  setSelected({ mediaType: sug.mediaType, title: sug.title, year: sug.year, posterPath: sug.posterPath });
                  if (sug.mediaType === 'movie') {
                    setMovieYear(sug.year || '');
                  } else {
                    setTvSeason('');
                    setTvEpisode('');
                  }
                  setShowSuggestions(false);
                }}
              >
                {sug.posterPath ? (
                  <img
                    className="suggestion-thumb"
                    src={`https://image.tmdb.org/t/p/w92${sug.posterPath}`}
                    alt={sug.title}
                  />
                ) : (
                  <span className="suggestion-type">{sug.mediaType === 'movie' ? '🎬' : '📺'}</span>
                )}
                <span className="suggestion-title">{sug.title}</span>
                {sug.year ? <span className="suggestion-year">({sug.year})</span> : null}
              </div>
            ))}
          </div>
        )}

        {selected.mediaType === 'movie' && (
          <div className="refine-row">
            <label className="refine-label">Release year</label>
            <input
              type="number"
              min={1900}
              max={3000}
              placeholder="e.g., 2008"
              value={movieYear}
              onChange={(e) => setMovieYear(e.target.value ? Number(e.target.value) : '')}
              className="refine-input"
            />
            <div className="composed-preview">Query: {composedQuery}</div>
          </div>
        )}

        {selected.mediaType === 'tv' && (
          <div className="refine-row">
            <label className="refine-label">Season</label>
            <input
              type="number"
              min={1}
              placeholder="e.g., 1"
              value={tvSeason}
              onChange={(e) => setTvSeason(e.target.value ? Number(e.target.value) : '')}
              className="refine-input"
            />
            <label className="refine-label">Episode</label>
            <input
              type="number"
              min={1}
              placeholder="e.g., 2"
              value={tvEpisode}
              onChange={(e) => setTvEpisode(e.target.value ? Number(e.target.value) : '')}
              className="refine-input"
            />
            <div className="composed-preview">Query: {composedQuery}</div>
          </div>
        )}
      </form>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      {loading && (
        <div className="loading">
          <div className="spinner"></div>
          <p>Searching...</p>
        </div>
      )}

      {results.length > 0 && (
        <div className="results-container">
          <div className="results-header">
            <h2>Search Results ({results.length})</h2>
            <div className="pagination">
              <button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage <= 1}
                className="page-button"
              >
                Previous
              </button>
              <span className="page-info">Page {currentPage}</span>
              <button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={results.length < 10}
                className="page-button"
              >
                Next
              </button>
            </div>
          </div>
          
          <div className="results-list">
            {results.map(renderResult)}
          </div>
        </div>
      )}

      {!loading && !error && hasSearched && results.length === 0 && (
        <div className="no-results">
          <p>No results found for "{lastSearchedQuery}"</p>
          <p>Try adjusting your search terms or check your spelling.</p>
        </div>
      )}
    </div>
  );
};

export default SearchComponent; 