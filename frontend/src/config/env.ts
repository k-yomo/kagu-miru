export const GA_ID = process.env.NEXT_PUBLIC_GOOGLE_ANALYTICS_ID;
export const GRAPHQL_API_URL =
  process.env.NEXT_PUBLIC_GRAPHQL_API_URL ||
  'http://localhost:8000/api/graphql';

export const SANITY_DATASET_ENV =
  process.env.NEXT_PUBLIC_SANITY_DATASET_ENV || 'development';
