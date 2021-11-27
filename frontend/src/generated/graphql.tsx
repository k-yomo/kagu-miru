import { gql } from '@apollo/client';
import * as Apollo from '@apollo/client';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
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
  Map: any;
  Time: any;
};

export enum Action {
  ClickItem = 'CLICK_ITEM',
  Display = 'DISPLAY',
}

export type Event = {
  action: Action;
  createdAt: Scalars['Time'];
  id: EventId;
  params: Scalars['Map'];
};

export enum EventId {
  QuerySuggestions = 'QUERY_SUGGESTIONS',
  Search = 'SEARCH',
}

export type Item = {
  affiliateUrl: Scalars['String'];
  averageRating: Scalars['Float'];
  categoryIds: Array<Scalars['ID']>;
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

export type ItemCategory = {
  Parent?: Maybe<ItemCategory>;
  children: Array<ItemCategory>;
  id: Scalars['ID'];
  level: Scalars['Int'];
  name: Scalars['String'];
  parentId?: Maybe<Scalars['ID']>;
};

export type ItemConnection = {
  nodes: Array<Item>;
  pageInfo: PageInfo;
};

export enum ItemSellingPlatform {
  Rakuten = 'RAKUTEN',
  YahooShopping = 'YAHOO_SHOPPING',
}

export enum ItemStatus {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
}

export type Mutation = {
  trackEvent: Scalars['Boolean'];
};

export type MutationTrackEventArgs = {
  event: Event;
};

export type PageInfo = {
  page: Scalars['Int'];
  totalCount: Scalars['Int'];
  totalPage: Scalars['Int'];
};

export type Query = {
  getAllItemCategories: Array<ItemCategory>;
  getItem: Item;
  getQuerySuggestions: QuerySuggestionsResponse;
  search: SearchResponse;
};

export type QueryGetItemArgs = {
  id: Scalars['ID'];
};

export type QueryGetQuerySuggestionsArgs = {
  query: Scalars['String'];
};

export type QuerySearchArgs = {
  input?: InputMaybe<SearchInput>;
};

export type QuerySuggestionsDisplayActionParams = {
  query: Scalars['String'];
  suggestedQueries: Array<Scalars['String']>;
};

export type QuerySuggestionsResponse = {
  query: Scalars['String'];
  suggestedQueries: Array<Scalars['String']>;
};

export type SearchClickItemActionParams = {
  itemId: Scalars['String'];
  searchId: Scalars['String'];
};

export type SearchDisplayItemsActionParams = {
  itemIds: Array<Scalars['ID']>;
  searchFrom: SearchFrom;
  searchId: Scalars['String'];
  searchInput: SearchInput;
};

export type SearchFilter = {
  categoryIds: Array<Scalars['ID']>;
  maxPrice?: InputMaybe<Scalars['Int']>;
  minPrice?: InputMaybe<Scalars['Int']>;
  minRating?: InputMaybe<Scalars['Int']>;
  platforms: Array<ItemSellingPlatform>;
};

export enum SearchFrom {
  Filter = 'FILTER',
  QuerySuggestion = 'QUERY_SUGGESTION',
  Search = 'SEARCH',
  Url = 'URL',
}

export type SearchInput = {
  filter: SearchFilter;
  page?: InputMaybe<Scalars['Int']>;
  pageSize?: InputMaybe<Scalars['Int']>;
  query: Scalars['String'];
  sortType: SearchSortType;
};

export type SearchResponse = {
  itemConnection: ItemConnection;
  searchId: Scalars['String'];
};

export enum SearchSortType {
  BestMatch = 'BEST_MATCH',
  PriceAsc = 'PRICE_ASC',
  PriceDesc = 'PRICE_DESC',
  Rating = 'RATING',
  ReviewCount = 'REVIEW_COUNT',
}

export type GetQuerySuggestionsQueryVariables = Exact<{
  query: Scalars['String'];
}>;

export type GetQuerySuggestionsQuery = {
  getQuerySuggestions: { query: string; suggestedQueries: Array<string> };
};

export type SearchQueryVariables = Exact<{
  input: SearchInput;
}>;

export type SearchQuery = {
  search: {
    searchId: string;
    itemConnection: {
      pageInfo: { page: number; totalPage: number; totalCount: number };
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
        categoryIds: Array<string>;
        platform: ItemSellingPlatform;
      }>;
    };
  };
};

export type TrackEventMutationVariables = Exact<{
  event: Event;
}>;

export type TrackEventMutation = { trackEvent: boolean };

export type ItemDetailPageGetItemQueryVariables = Exact<{
  id: Scalars['ID'];
}>;

export type ItemDetailPageGetItemQuery = {
  getItem: {
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
    categoryIds: Array<string>;
    platform: ItemSellingPlatform;
  };
};

export const GetQuerySuggestionsDocument = gql`
  query getQuerySuggestions($query: String!) {
    getQuerySuggestions(query: $query) {
      query
      suggestedQueries
    }
  }
`;

/**
 * __useGetQuerySuggestionsQuery__
 *
 * To run a query within a React component, call `useGetQuerySuggestionsQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetQuerySuggestionsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetQuerySuggestionsQuery({
 *   variables: {
 *      query: // value for 'query'
 *   },
 * });
 */
export function useGetQuerySuggestionsQuery(
  baseOptions: Apollo.QueryHookOptions<
    GetQuerySuggestionsQuery,
    GetQuerySuggestionsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    GetQuerySuggestionsQuery,
    GetQuerySuggestionsQueryVariables
  >(GetQuerySuggestionsDocument, options);
}
export function useGetQuerySuggestionsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    GetQuerySuggestionsQuery,
    GetQuerySuggestionsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    GetQuerySuggestionsQuery,
    GetQuerySuggestionsQueryVariables
  >(GetQuerySuggestionsDocument, options);
}
export type GetQuerySuggestionsQueryHookResult = ReturnType<
  typeof useGetQuerySuggestionsQuery
>;
export type GetQuerySuggestionsLazyQueryHookResult = ReturnType<
  typeof useGetQuerySuggestionsLazyQuery
>;
export type GetQuerySuggestionsQueryResult = Apollo.QueryResult<
  GetQuerySuggestionsQuery,
  GetQuerySuggestionsQueryVariables
>;
export const SearchDocument = gql`
  query search($input: SearchInput!) {
    search(input: $input) {
      searchId
      itemConnection {
        pageInfo {
          page
          totalPage
          totalCount
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
          categoryIds
          platform
        }
      }
    }
  }
`;

/**
 * __useSearchQuery__
 *
 * To run a query within a React component, call `useSearchQuery` and pass it any options that fit your needs.
 * When your component renders, `useSearchQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useSearchQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useSearchQuery(
  baseOptions: Apollo.QueryHookOptions<SearchQuery, SearchQueryVariables>
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<SearchQuery, SearchQueryVariables>(
    SearchDocument,
    options
  );
}
export function useSearchLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<SearchQuery, SearchQueryVariables>
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<SearchQuery, SearchQueryVariables>(
    SearchDocument,
    options
  );
}
export type SearchQueryHookResult = ReturnType<typeof useSearchQuery>;
export type SearchLazyQueryHookResult = ReturnType<typeof useSearchLazyQuery>;
export type SearchQueryResult = Apollo.QueryResult<
  SearchQuery,
  SearchQueryVariables
>;
export const TrackEventDocument = gql`
  mutation trackEvent($event: Event!) {
    trackEvent(event: $event)
  }
`;
export type TrackEventMutationFn = Apollo.MutationFunction<
  TrackEventMutation,
  TrackEventMutationVariables
>;

/**
 * __useTrackEventMutation__
 *
 * To run a mutation, you first call `useTrackEventMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useTrackEventMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [trackEventMutation, { data, loading, error }] = useTrackEventMutation({
 *   variables: {
 *      event: // value for 'event'
 *   },
 * });
 */
export function useTrackEventMutation(
  baseOptions?: Apollo.MutationHookOptions<
    TrackEventMutation,
    TrackEventMutationVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useMutation<TrackEventMutation, TrackEventMutationVariables>(
    TrackEventDocument,
    options
  );
}
export type TrackEventMutationHookResult = ReturnType<
  typeof useTrackEventMutation
>;
export type TrackEventMutationResult =
  Apollo.MutationResult<TrackEventMutation>;
export type TrackEventMutationOptions = Apollo.BaseMutationOptions<
  TrackEventMutation,
  TrackEventMutationVariables
>;
export const ItemDetailPageGetItemDocument = gql`
  query itemDetailPageGetItem($id: ID!) {
    getItem(id: $id) {
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
      categoryIds
      platform
    }
  }
`;

/**
 * __useItemDetailPageGetItemQuery__
 *
 * To run a query within a React component, call `useItemDetailPageGetItemQuery` and pass it any options that fit your needs.
 * When your component renders, `useItemDetailPageGetItemQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useItemDetailPageGetItemQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useItemDetailPageGetItemQuery(
  baseOptions: Apollo.QueryHookOptions<
    ItemDetailPageGetItemQuery,
    ItemDetailPageGetItemQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    ItemDetailPageGetItemQuery,
    ItemDetailPageGetItemQueryVariables
  >(ItemDetailPageGetItemDocument, options);
}
export function useItemDetailPageGetItemLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    ItemDetailPageGetItemQuery,
    ItemDetailPageGetItemQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    ItemDetailPageGetItemQuery,
    ItemDetailPageGetItemQueryVariables
  >(ItemDetailPageGetItemDocument, options);
}
export type ItemDetailPageGetItemQueryHookResult = ReturnType<
  typeof useItemDetailPageGetItemQuery
>;
export type ItemDetailPageGetItemLazyQueryHookResult = ReturnType<
  typeof useItemDetailPageGetItemLazyQuery
>;
export type ItemDetailPageGetItemQueryResult = Apollo.QueryResult<
  ItemDetailPageGetItemQuery,
  ItemDetailPageGetItemQueryVariables
>;
