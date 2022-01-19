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

export type AppliedMetadata = {
  name: Scalars['String'];
  values: Array<Scalars['String']>;
};

export enum ErrorCode {
  Internal = 'INTERNAL',
  NotFound = 'NOT_FOUND',
}

export type Event = {
  action: Action;
  createdAt: Scalars['Time'];
  id: EventId;
  params: Scalars['Map'];
};

export enum EventId {
  Home = 'HOME',
  QuerySuggestions = 'QUERY_SUGGESTIONS',
  Search = 'SEARCH',
  SimilarItems = 'SIMILAR_ITEMS',
}

export type Facet = {
  facetType: FacetType;
  title: Scalars['String'];
  totalCount: Scalars['Int'];
  values: Array<FacetValue>;
};

export enum FacetType {
  BrandNames = 'BRAND_NAMES',
  CategoryIds = 'CATEGORY_IDS',
  Colors = 'COLORS',
  Metadata = 'METADATA',
}

export type FacetValue = {
  count: Scalars['Int'];
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type GetSimilarItemsInput = {
  itemId: Scalars['ID'];
  page?: InputMaybe<Scalars['Int']>;
  pageSize?: InputMaybe<Scalars['Int']>;
};

export type GetSimilarItemsResponse = {
  itemConnection: ItemConnection;
  searchId: Scalars['String'];
};

export type HomeClickItemActionParams = {
  componentId: Scalars['ID'];
  itemId: Scalars['String'];
};

export type HomeComponent = {
  id: Scalars['ID'];
  payload: HomeComponentPayload;
};

export type HomeComponentPayload =
  | HomeComponentPayloadCategories
  | HomeComponentPayloadItemGroups
  | HomeComponentPayloadItems
  | HomeComponentPayloadMediaPosts;

export type HomeComponentPayloadCategories = {
  categories: Array<ItemCategory>;
  title: Scalars['String'];
};

export type HomeComponentPayloadItemGroups = {
  payload: Array<HomeComponentPayloadItems>;
  title: Scalars['String'];
};

export type HomeComponentPayloadItems = {
  items: Array<Item>;
  title: Scalars['String'];
};

export type HomeComponentPayloadMediaPosts = {
  posts: Array<MediaPost>;
  title: Scalars['String'];
};

export type HomeResponse = {
  components: Array<HomeComponent>;
};

export type Item = {
  affiliateUrl: Scalars['String'];
  averageRating: Scalars['Float'];
  categoryId: Scalars['ID'];
  colors: Array<ItemColor>;
  description: Scalars['String'];
  groupID: Scalars['ID'];
  id: Scalars['ID'];
  imageUrls: Array<Scalars['String']>;
  name: Scalars['String'];
  platform: ItemSellingPlatform;
  price: Scalars['Int'];
  reviewCount: Scalars['Int'];
  sameGroupItems: Array<Item>;
  status: ItemStatus;
  url: Scalars['String'];
};

export type ItemCategory = {
  children: Array<ItemCategory>;
  id: Scalars['ID'];
  imageUrl?: Maybe<Scalars['String']>;
  level: Scalars['Int'];
  name: Scalars['String'];
  parent?: Maybe<ItemCategory>;
  parentId?: Maybe<Scalars['ID']>;
};

export enum ItemColor {
  Beige = 'BEIGE',
  Black = 'BLACK',
  Blue = 'BLUE',
  Brown = 'BROWN',
  Gold = 'GOLD',
  Gray = 'GRAY',
  Green = 'GREEN',
  Khaki = 'KHAKI',
  Navy = 'NAVY',
  Orange = 'ORANGE',
  Pink = 'PINK',
  Purple = 'PURPLE',
  Red = 'RED',
  Silver = 'SILVER',
  Transparent = 'TRANSPARENT',
  White = 'WHITE',
  WineRed = 'WINE_RED',
  Yellow = 'YELLOW',
}

export type ItemConnection = {
  nodes: Array<Item>;
  pageInfo: PageInfo;
};

export enum ItemSellingPlatform {
  PaypayMall = 'PAYPAY_MALL',
  Rakuten = 'RAKUTEN',
  YahooShopping = 'YAHOO_SHOPPING',
}

export enum ItemStatus {
  Active = 'ACTIVE',
  Inactive = 'INACTIVE',
}

export type MediaPost = {
  categories: Array<MediaPostCategory>;
  description: Scalars['String'];
  id: Scalars['ID'];
  mainImageUrl: Scalars['String'];
  publishedAt: Scalars['Time'];
  slug: Scalars['String'];
  title: Scalars['String'];
};

export type MediaPostCategory = {
  id: Scalars['ID'];
  names: Array<Scalars['String']>;
};

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
  getSimilarItems: GetSimilarItemsResponse;
  home: HomeResponse;
  search: SearchResponse;
};

export type QueryGetItemArgs = {
  id: Scalars['ID'];
};

export type QueryGetQuerySuggestionsArgs = {
  query: Scalars['String'];
};

export type QueryGetSimilarItemsArgs = {
  input: GetSimilarItemsInput;
};

export type QuerySearchArgs = {
  input: SearchInput;
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
  brandNames: Array<Scalars['String']>;
  categoryIds: Array<Scalars['ID']>;
  colors: Array<ItemColor>;
  maxPrice?: InputMaybe<Scalars['Int']>;
  metadata: Array<AppliedMetadata>;
  minPrice?: InputMaybe<Scalars['Int']>;
  minRating?: InputMaybe<Scalars['Int']>;
  platforms: Array<ItemSellingPlatform>;
};

export enum SearchFrom {
  Filter = 'FILTER',
  Home = 'HOME',
  Media = 'MEDIA',
  OpenSearch = 'OPEN_SEARCH',
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
  facets: Array<Facet>;
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

export type SimilarItemsDisplayItemsActionParams = {
  getSimilarItemsInput: GetSimilarItemsInput;
  itemIds: Array<Scalars['ID']>;
  searchId: Scalars['String'];
};

export type ItemListItemFragmentFragment = {
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
  categoryId: string;
  platform: ItemSellingPlatform;
};

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
        categoryId: string;
        platform: ItemSellingPlatform;
      }>;
    };
    facets: Array<{
      title: string;
      facetType: FacetType;
      totalCount: number;
      values: Array<{ id: string; name: string; count: number }>;
    }>;
  };
};

export type TrackEventMutationVariables = Exact<{
  event: Event;
}>;

export type TrackEventMutation = { trackEvent: boolean };

export type HomeQueryVariables = Exact<{ [key: string]: never }>;

export type HomeQuery = {
  home: {
    components: Array<{
      id: string;
      payload:
        | {
            __typename: 'HomeComponentPayloadCategories';
            title: string;
            categories: Array<{
              id: string;
              name: string;
              imageUrl?: string | null | undefined;
            }>;
          }
        | {
            __typename: 'HomeComponentPayloadItemGroups';
            title: string;
            payload: Array<{
              title: string;
              items: Array<{
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
                categoryId: string;
                platform: ItemSellingPlatform;
              }>;
            }>;
          }
        | {
            __typename: 'HomeComponentPayloadItems';
            title: string;
            items: Array<{
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
              categoryId: string;
              platform: ItemSellingPlatform;
            }>;
          }
        | {
            __typename: 'HomeComponentPayloadMediaPosts';
            title: string;
            posts: Array<{
              id: string;
              slug: string;
              title: string;
              description: string;
              mainImageUrl: string;
              publishedAt: any;
              categories: Array<{ id: string; names: Array<string> }>;
            }>;
          };
    }>;
  };
};

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
    categoryId: string;
    platform: ItemSellingPlatform;
    sameGroupItems: Array<{
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
      categoryId: string;
      platform: ItemSellingPlatform;
    }>;
  };
};

export type ItemDetailPageGetSimilarItemsQueryVariables = Exact<{
  input: GetSimilarItemsInput;
}>;

export type ItemDetailPageGetSimilarItemsQuery = {
  getSimilarItems: {
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
        categoryId: string;
        platform: ItemSellingPlatform;
      }>;
    };
  };
};

export const ItemListItemFragmentFragmentDoc = gql`
  fragment itemListItemFragment on Item {
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
    categoryId
    platform
  }
`;
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
          ...itemListItemFragment
        }
      }
      facets {
        title
        facetType
        values {
          id
          name
          count
        }
        totalCount
      }
    }
  }
  ${ItemListItemFragmentFragmentDoc}
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
export const HomeDocument = gql`
  query home {
    home {
      components {
        id
        payload {
          __typename
          ... on HomeComponentPayloadItemGroups {
            title
            payload {
              ... on HomeComponentPayloadItems {
                title
                items {
                  ...itemListItemFragment
                }
              }
            }
          }
          ... on HomeComponentPayloadItems {
            title
            items {
              ...itemListItemFragment
            }
          }
          ... on HomeComponentPayloadCategories {
            title
            categories {
              id
              name
              imageUrl
            }
          }
          ... on HomeComponentPayloadMediaPosts {
            title
            posts {
              id
              slug
              title
              description
              mainImageUrl
              publishedAt
              categories {
                id
                names
              }
            }
          }
        }
      }
    }
  }
  ${ItemListItemFragmentFragmentDoc}
`;

/**
 * __useHomeQuery__
 *
 * To run a query within a React component, call `useHomeQuery` and pass it any options that fit your needs.
 * When your component renders, `useHomeQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useHomeQuery({
 *   variables: {
 *   },
 * });
 */
export function useHomeQuery(
  baseOptions?: Apollo.QueryHookOptions<HomeQuery, HomeQueryVariables>
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<HomeQuery, HomeQueryVariables>(HomeDocument, options);
}
export function useHomeLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<HomeQuery, HomeQueryVariables>
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<HomeQuery, HomeQueryVariables>(
    HomeDocument,
    options
  );
}
export type HomeQueryHookResult = ReturnType<typeof useHomeQuery>;
export type HomeLazyQueryHookResult = ReturnType<typeof useHomeLazyQuery>;
export type HomeQueryResult = Apollo.QueryResult<HomeQuery, HomeQueryVariables>;
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
      categoryId
      platform
      sameGroupItems {
        ...itemListItemFragment
      }
    }
  }
  ${ItemListItemFragmentFragmentDoc}
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
export const ItemDetailPageGetSimilarItemsDocument = gql`
  query itemDetailPageGetSimilarItems($input: GetSimilarItemsInput!) {
    getSimilarItems(input: $input) {
      searchId
      itemConnection {
        pageInfo {
          page
          totalPage
          totalCount
        }
        nodes {
          ...itemListItemFragment
        }
      }
    }
  }
  ${ItemListItemFragmentFragmentDoc}
`;

/**
 * __useItemDetailPageGetSimilarItemsQuery__
 *
 * To run a query within a React component, call `useItemDetailPageGetSimilarItemsQuery` and pass it any options that fit your needs.
 * When your component renders, `useItemDetailPageGetSimilarItemsQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useItemDetailPageGetSimilarItemsQuery({
 *   variables: {
 *      input: // value for 'input'
 *   },
 * });
 */
export function useItemDetailPageGetSimilarItemsQuery(
  baseOptions: Apollo.QueryHookOptions<
    ItemDetailPageGetSimilarItemsQuery,
    ItemDetailPageGetSimilarItemsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    ItemDetailPageGetSimilarItemsQuery,
    ItemDetailPageGetSimilarItemsQueryVariables
  >(ItemDetailPageGetSimilarItemsDocument, options);
}
export function useItemDetailPageGetSimilarItemsLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    ItemDetailPageGetSimilarItemsQuery,
    ItemDetailPageGetSimilarItemsQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    ItemDetailPageGetSimilarItemsQuery,
    ItemDetailPageGetSimilarItemsQueryVariables
  >(ItemDetailPageGetSimilarItemsDocument, options);
}
export type ItemDetailPageGetSimilarItemsQueryHookResult = ReturnType<
  typeof useItemDetailPageGetSimilarItemsQuery
>;
export type ItemDetailPageGetSimilarItemsLazyQueryHookResult = ReturnType<
  typeof useItemDetailPageGetSimilarItemsLazyQuery
>;
export type ItemDetailPageGetSimilarItemsQueryResult = Apollo.QueryResult<
  ItemDetailPageGetSimilarItemsQuery,
  ItemDetailPageGetSimilarItemsQueryVariables
>;
