import axios from 'axios';

export type TMDBMediaType = 'movie' | 'tv';

export interface TMDBSuggestion {
  id: number;
  mediaType: TMDBMediaType;
  title: string;
  year?: number;
  posterPath?: string;
}

const TMDB_API_KEY = process.env.REACT_APP_TMDB_API_KEY;
const TMDB_BASE_URL = 'https://api.themoviedb.org/3';

interface TMDBMultiSearchResponse {
  results: any[];
}

export async function searchTMDB(query: string): Promise<TMDBSuggestion[]> {
  if (!TMDB_API_KEY) {
    return [];
  }

  try {
    const url = `${TMDB_BASE_URL}/search/multi`;
    const response = await axios.get<TMDBMultiSearchResponse>(url, {
      params: {
        query,
        include_adult: false,
        language: 'en-US',
        page: 1,
        api_key: TMDB_API_KEY
      }
    });

    const rawResults = Array.isArray(response.data?.results) ? response.data.results : [];
    const results: TMDBSuggestion[] = rawResults
      .filter((r: any) => r.media_type === 'movie' || r.media_type === 'tv')
      .map((r: any) => {
        if (r.media_type === 'movie') {
          const year = r.release_date ? Number(String(r.release_date).slice(0, 4)) : undefined;
          return {
            id: r.id,
            mediaType: 'movie' as const,
            title: r.title || r.original_title,
            year,
            posterPath: r.poster_path || undefined
          };
        }
        const year = r.first_air_date ? Number(String(r.first_air_date).slice(0, 4)) : undefined;
        return {
          id: r.id,
          mediaType: 'tv' as const,
          title: r.name || r.original_name,
          year,
          posterPath: r.poster_path || undefined
        };
      });

    return results;
  } catch (err) {
    console.error('TMDB search error:', err);
    return [];
  }
}


