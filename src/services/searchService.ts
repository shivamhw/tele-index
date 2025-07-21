import axios from 'axios';
import { SearchRequest, SearchResult } from '../types/search';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8095';

const searchApi = axios.create({
  baseURL: API_BASE_URL,
});

export const searchTele = async (query: string, page: number = 1, size: number = 10): Promise<SearchResult[]> => {
  try {
    const from = (page - 1) * size;
    
    const requestBody: SearchRequest = {
      size,
      from,
      explain: true,
      highlight: {},
      query: {
        boost: 1,
        match: query
      },
      fields: ['*']
    };

    const response = await searchApi.post<any>('/api/tele/_search', requestBody);
    console.log('API response:', response.data);
    if (!Array.isArray(response.data.hits)) {
      return [];
    }
    return response.data.hits.map((hit: any) => ({
      id: hit.id,
      score: hit.score,
      source: hit.fields,
      highlight: hit.fragments
    }));
  } catch (error) {
    console.error('Search API error:', error);
    throw error;
  }
}; 