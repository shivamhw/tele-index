export interface SearchRequest {
  size: number;
  from: number;
  explain: boolean;
  highlight: Record<string, any>;
  query: {
    boost: number;
    match: string;
  };
  fields: string[];
}

export interface SearchResponse {
  hits: {
    total: {
      value: number;
      relation: string;
    };
    hits: Array<{
      _index: string;
      _id: string;
      _score: number;
      _source: Record<string, any>;
      highlight?: Record<string, string[]>;
    }>;
  };
  took: number;
  timed_out: boolean;
}

export interface SearchResult {
  id: string;
  score: number;
  source: Record<string, any>;
  highlight?: Record<string, string[]>;
} 