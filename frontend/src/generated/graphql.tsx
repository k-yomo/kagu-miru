import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
const defaultOptions = {};
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Item = {
  affiliateUrl: Scalars['String'];
  averageRating: Scalars['Float'];
  description: Scalars['String'];
  id: Scalars['ID'];
  imageUrls: Array<Scalars['String']>;
  name: Scalars['String'];
  platform: ItemSellingPlatform;
  price: Scalars['Int'];
  reviewCount: Scalars['Int'];
  status: ItemStatus;
  url: Scalars['String'];
};

export type ItemConnection = {
  nodes: Array<Item>;
  pageInfo: PageInfo;
};

export enum ItemSellingPlatform {
  Rakuten = 'RAKUTEN',
}

export enum ItemStatus {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
}

export type PageInfo = {
  page: Scalars['Int'];
  totalPage: Scalars['Int'];
};

export type Query = {
  getQuerySuggestions: Array<Scalars['String']>;
  searchItems: ItemConnection;
};

export type QueryGetQuerySuggestionsArgs = {
  query: Scalars['String'];
};

export type QuerySearchItemsArgs = {
  input?: Maybe<SearchItemsInput>;
};

export type SearchItemsInput = {
  page?: Maybe<Scalars['Int']>;
  pageSize?: Maybe<Scalars['Int']>;
  query: Scalars['String'];
  sortType: SearchItemsSortType;
};

export enum SearchItemsSortType {
  BestMatch = 'BEST_MATCH',
  SortByPriceAsc = 'SORT_BY_PRICE_ASC',
  SortByPriceDesc = 'SORT_BY_PRICE_DESC',
}

export type HomePageSearchItemsQueryVariables = Exact<{
  input: SearchItemsInput;
}>;

export type HomePageSearchItemsQuery = {
  searchItems: {
    pageInfo: { page: number; totalPage: number };
    nodes: Array<{
      id: string;
      name: string;
      description: string;
      status: ItemStatus;
      url: string;
      affiliateUrl: string;
      price: number;
      imageUrls: Array<string>;
      averageRating: number;
      reviewCount: number;
      platform: ItemSellingPlatform;
    }>;
  };
};

export type HomePageGetQuerySuggestionsQueryVariables = Exact<{
  query: Scalars['String'];
}>;

export type HomePageGetQuerySuggestionsQuery = {
  getQuerySuggestions: Array<string>;
};

export const HomePageSearchItemsDocument = gql`
  query homePageSearchItems($input: SearchItemsInput!) {
    searchItems(input: $input) {
      pageInfo {
        page
        totalPage
      }
      nodes {
        id
        name
        description
        status
        url
        affiliateUrl
        price
        imageUrls
        averageRating
        reviewCount
        platform
      }
    }
  }
`;

/**
 * __useHomePageSearchItemsQuery__
 *
 * To run a query within a React component, call `useHomePageSearchItemsQuery` and pass it any options that fit your needs.
 * When your component renders, `useHomePageSearchItemsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHomePageSearchItemsQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useHomePageSearchItemsQuery(
  baseOptions: Apollo.QueryHookOptions<
    HomePageSearchItemsQuery,
    HomePageSearchItemsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    HomePageSearchItemsQuery,
    HomePageSearchItemsQueryVariables
  >(HomePageSearchItemsDocument, options);
}
export function useHomePageSearchItemsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    HomePageSearchItemsQuery,
    HomePageSearchItemsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    HomePageSearchItemsQuery,
    HomePageSearchItemsQueryVariables
  >(HomePageSearchItemsDocument, options);
}
export type HomePageSearchItemsQueryHookResult = ReturnType<
  typeof useHomePageSearchItemsQuery
>;
export type HomePageSearchItemsLazyQueryHookResult = ReturnType<
  typeof useHomePageSearchItemsLazyQuery
>;
export type HomePageSearchItemsQueryResult = Apollo.QueryResult<
  HomePageSearchItemsQuery,
  HomePageSearchItemsQueryVariables
>;
export const HomePageGetQuerySuggestionsDocument = gql`
  query homePageGetQuerySuggestions($query: String!) {
    getQuerySuggestions(query: $query)
  }
`;

/**
 * __useHomePageGetQuerySuggestionsQuery__
 *
 * To run a query within a React component, call `useHomePageGetQuerySuggestionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useHomePageGetQuerySuggestionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHomePageGetQuerySuggestionsQuery({
 *   variables: {
 *      query: // value for 'query'
 *   },
 * });
 */
export function useHomePageGetQuerySuggestionsQuery(
  baseOptions: Apollo.QueryHookOptions<
    HomePageGetQuerySuggestionsQuery,
    HomePageGetQuerySuggestionsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    HomePageGetQuerySuggestionsQuery,
    HomePageGetQuerySuggestionsQueryVariables
  >(HomePageGetQuerySuggestionsDocument, options);
}
export function useHomePageGetQuerySuggestionsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    HomePageGetQuerySuggestionsQuery,
    HomePageGetQuerySuggestionsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    HomePageGetQuerySuggestionsQuery,
    HomePageGetQuerySuggestionsQueryVariables
  >(HomePageGetQuerySuggestionsDocument, options);
}
export type HomePageGetQuerySuggestionsQueryHookResult = ReturnType<
  typeof useHomePageGetQuerySuggestionsQuery
>;
export type HomePageGetQuerySuggestionsLazyQueryHookResult = ReturnType<
  typeof useHomePageGetQuerySuggestionsLazyQuery
>;
export type HomePageGetQuerySuggestionsQueryResult = Apollo.QueryResult<
  HomePageGetQuerySuggestionsQuery,
  HomePageGetQuerySuggestionsQueryVariables
>;
