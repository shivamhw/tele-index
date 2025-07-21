import React, { useState } from 'react';
import { searchTele } from '../services/searchService';
import { SearchResult } from '../types/search';
import './SearchComponent.css';
import { useNavigate, useSearchParams } from 'react-router-dom';

const DOWNLOAD_BASE_URL = process.env.REACT_APP_DOWNLOAD_BASE_URL || 'http://localhost:8080';

const SearchComponent: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const initialQuery = searchParams.get('q') || '';
  const initialPage = parseInt(searchParams.get('page') || '1', 10);
  const [query, setQuery] = useState(initialQuery);
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [currentPage, setCurrentPage] = useState(initialPage);
  const navigate = useNavigate();

  // On mount, if query param exists, trigger search
  React.useEffect(() => {
    if (initialQuery) {
      handleSearch(initialQuery, initialPage);
    }
    // eslint-disable-next-line
  }, []);

  const handleSearch = async (searchQuery: string = query, page: number = 1) => {
    if (!searchQuery.trim()) {
      setError('Please enter a search query');
      setResults([]);
      setCurrentPage(1);
      setSearchParams({});
      return;
    }
    setLoading(true);
    setError(null);
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
    handleSearch(query, 1);
  };

  const handlePageChange = (newPage: number) => {
    handleSearch(query, newPage);
  };

  // Updated renderResult for your API's hit structure
  const renderResult = (result: SearchResult) => {
    const fromID = result.source.PeerID;
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
            <span className="result-size">Size: {result.source.Size}</span>
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
        <h1>Tele Search</h1>
        <p>Search through the tele index with advanced querying</p>
      </div>

      <form onSubmit={handleSubmit} className="search-form">
        <div className="search-input-group">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Enter your search query (e.g., 'housefull 5a')"
            className="search-input"
            disabled={loading}
          />
          <button type="submit" className="search-button" disabled={loading}>
            {loading ? 'Searching...' : 'Search'}
          </button>
        </div>
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

      {!loading && !error && results.length === 0 && query && (
        <div className="no-results">
          <p>No results found for "{query}"</p>
          <p>Try adjusting your search terms or check your spelling.</p>
        </div>
      )}
    </div>
  );
};

export default SearchComponent; 