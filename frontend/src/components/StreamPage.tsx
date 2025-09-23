import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

const StreamPage: React.FC = () => {
  const query = useQuery();
  const url = query.get('url');
  const navigate = useNavigate();

  if (!url) {
    return <div className="stream-container"><h2 className="stream-header">No stream URL provided.</h2></div>;
  }

  // Simple file type check for video/audio
  const isVideo = url.match(/\.(mp4|webm|ogg)$/i);
  const isAudio = url.match(/\.(mp3|wav|ogg)$/i);

  return (
    <div className="stream-container">
      <div className="stream-header">
        <h1>Streaming</h1>
        <p>Enjoy your media directly in the browser</p>
      </div>
      <button className="download-button" style={{ marginBottom: '1.5rem' }} onClick={() => navigate(-1)}>
        ← Back to Search
      </button>
      <div className="stream-player-wrapper">
        <p>
          <a href={url} target="_blank" rel="noopener noreferrer" className="result-link">Direct Link</a>
        </p>
        {isVideo ? (
          <video src={url} controls style={{ maxWidth: '100%', maxHeight: 480, borderRadius: 12, boxShadow: '0 4px 20px rgba(0,0,0,0.08)' }} autoPlay />
        ) : isAudio ? (
          <audio src={url} controls autoPlay style={{ width: '100%' }} />
        ) : (
          <div>
            <p>Cannot determine file type. Attempting to play as video:</p>
            <video src={url} controls style={{ maxWidth: '100%', maxHeight: 480, borderRadius: 12, boxShadow: '0 4px 20px rgba(0,0,0,0.08)' }} autoPlay />
          </div>
        )}
      </div>
    </div>
  );
};

export default StreamPage; 