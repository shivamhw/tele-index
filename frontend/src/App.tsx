import React from 'react';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import SearchComponent from './components/SearchComponent';
import StreamPage from './components/StreamPage';

function App() {
  return (
    <div className="App">
      <Router>
        <Routes>
          <Route path="/" element={<SearchComponent />} />
          <Route path="/stream" element={<StreamPage />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
