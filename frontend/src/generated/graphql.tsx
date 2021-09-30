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

export type Mutation = {
  trackEvent: Scalars['Boolean'];
};

export type MutationTrackEventArgs = {
  event: Event;
};

export type PageInfo = {
  page: Scalars['Int'];
  totalPage: Scalars['Int'];
};

export type Query = {
  getQuerySuggestions: QuerySuggestionsResponse;
  search: SearchResponse;
};

export type QueryGetQuerySuggestionsArgs = {
  query: Scalars['String'];
};

export type QuerySearchArgs = {
  input?: Maybe<SearchInput>;
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

export enum SearchFrom {
  QuerySuggestion = 'QUERY_SUGGESTION',
  Search = 'SEARCH',
  Url = 'URL',
}

export type SearchInput = {
  page?: Maybe<Scalars['Int']>;
  pageSize?: Maybe<Scalars['Int']>;
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

export type HomePageSearchQueryVariables = Exact<{
  input: SearchInput;
}>;

export type HomePageSearchQuery = {
  search: {
    searchId: string;
    itemConnection: {
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
};

export type HomePageGetQuerySuggestionsQueryVariables = Exact<{
  query: Scalars['String'];
}>;

export type HomePageGetQuerySuggestionsQuery = {
  getQuerySuggestions: { query: string; suggestedQueries: Array<string> };
};

export type HomePageTrackEventMutationVariables = Exact<{
  event: Event;
}>;

export type HomePageTrackEventMutation = { trackEvent: boolean };

export const HomePageSearchDocument = gql`
  query homePageSearch($input: SearchInput!) {
    search(input: $input) {
      searchId
      itemConnection {
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
  }
`;

/**
 * __useHomePageSearchQuery__
 *
 * To run a query within a React component, call `useHomePageSearchQuery` and pass it any options that fit your needs.
 * When your component renders, `useHomePageSearchQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHomePageSearchQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useHomePageSearchQuery(
  baseOptions: Apollo.QueryHookOptions<
    HomePageSearchQuery,
    HomePageSearchQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<HomePageSearchQuery, HomePageSearchQueryVariables>(
    HomePageSearchDocument,
    options
  );
}
export function useHomePageSearchLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    HomePageSearchQuery,
    HomePageSearchQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<HomePageSearchQuery, HomePageSearchQueryVariables>(
    HomePageSearchDocument,
    options
  );
}
export type HomePageSearchQueryHookResult = ReturnType<
  typeof useHomePageSearchQuery
>;
export type HomePageSearchLazyQueryHookResult = ReturnType<
  typeof useHomePageSearchLazyQuery
>;
export type HomePageSearchQueryResult = Apollo.QueryResult<
  HomePageSearchQuery,
  HomePageSearchQueryVariables
>;
export const HomePageGetQuerySuggestionsDocument = gql`
  query homePageGetQuerySuggestions($query: String!) {
    getQuerySuggestions(query: $query) {
      query
      suggestedQueries
    }
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
export const HomePageTrackEventDocument = gql`
  mutation homePageTrackEvent($event: Event!) {
    trackEvent(event: $event)
  }
`;
export type HomePageTrackEventMutationFn = Apollo.MutationFunction<
  HomePageTrackEventMutation,
  HomePageTrackEventMutationVariables
>;

/**
 * __useHomePageTrackEventMutation__
 *
 * To run a mutation, you first call `useHomePageTrackEventMutation` within a React component and pass it any options that fit your needs.
 * When your component renders, `useHomePageTrackEventMutation` returns a tuple that includes:
 * - A mutate function that you can call at any time to execute the mutation
 * - An object with fields that represent the current status of the mutation's execution
 *
 * @param baseOptions options that will be passed into the mutation, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options-2;
 *
 * @example
 * const [homePageTrackEventMutation, { data, loading, error }] = useHomePageTrackEventMutation({
 *   variables: {
 *      event: // value for 'event'
 *   },
 * });
 */
export function useHomePageTrackEventMutation(
  baseOptions?: Apollo.MutationHookOptions<
    HomePageTrackEventMutation,
    HomePageTrackEventMutationVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useMutation<
    HomePageTrackEventMutation,
    HomePageTrackEventMutationVariables
  >(HomePageTrackEventDocument, options);
}
export type HomePageTrackEventMutationHookResult = ReturnType<
  typeof useHomePageTrackEventMutation
>;
export type HomePageTrackEventMutationResult =
  Apollo.MutationResult<HomePageTrackEventMutation>;
export type HomePageTrackEventMutationOptions = Apollo.BaseMutationOptions<
  HomePageTrackEventMutation,
  HomePageTrackEventMutationVariables
>;
