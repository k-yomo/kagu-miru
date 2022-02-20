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
const defaultOptions = {} as const;
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
  mainImageUrl: Scalars['String'];
  publishedAt: Scalars['Time'];
  slug: Scalars['ID'];
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
  brandNames?: InputMaybe<Array<Scalars['String']>>;
  categoryIds?: InputMaybe<Array<Scalars['ID']>>;
  colors?: InputMaybe<Array<ItemColor>>;
  maxPrice?: InputMaybe<Scalars['Int']>;
  metadata?: InputMaybe<Array<AppliedMetadata>>;
  minPrice?: InputMaybe<Scalars['Int']>;
  minRating?: InputMaybe<Scalars['Int']>;
  platforms?: InputMaybe<Array<ItemSellingPlatform>>;
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
  filter?: InputMaybe<SearchFilter>;
  page?: InputMaybe<Scalars['Int']>;
  pageSize?: InputMaybe<Scalars['Int']>;
  query: Scalars['String'];
  sortType?: InputMaybe<SearchSortType>;
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

export type SubItemCategoryFragment = {
  id: string;
  level: number;
  name: string;
  parentId?: string | null;
};

export type CategoryInputGetAllCategoriesQueryVariables = Exact<{
  [key: string]: never;
}>;

export type CategoryInputGetAllCategoriesQuery = {
  getAllItemCategories: Array<{
    id: string;
    level: number;
    name: string;
    parentId?: string | null;
    children: Array<{
      id: string;
      level: number;
      name: string;
      parentId?: string | null;
      children: Array<{
        id: string;
        level: number;
        name: string;
        parentId?: string | null;
        children: Array<{
          id: string;
          level: number;
          name: string;
          parentId?: string | null;
        }>;
      }>;
    }>;
  }>;
};

export type ItemPreviewGetItemQueryVariables = Exact<{
  id: Scalars['ID'];
}>;

export type ItemPreviewGetItemQuery = {
  getItem: {
    id: string;
    name: string;
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
};

export const SubItemCategoryFragmentDoc = gql`
  fragment subItemCategory on ItemCategory {
    id
    level
    name
    parentId
  }
`;
export const CategoryInputGetAllCategoriesDocument = gql`
  query categoryInputGetAllCategories {
    getAllItemCategories {
      ...subItemCategory
      children {
        ...subItemCategory
        children {
          ...subItemCategory
          children {
            ...subItemCategory
          }
        }
      }
    }
  }
  ${SubItemCategoryFragmentDoc}
`;

/**
 * __useCategoryInputGetAllCategoriesQuery__
 *
 * To run a query within a React component, call `useCategoryInputGetAllCategoriesQuery` and pass it any options that fit your needs.
 * When your component renders, `useCategoryInputGetAllCategoriesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useCategoryInputGetAllCategoriesQuery({
 *   variables: {
 *   },
 * });
 */
export function useCategoryInputGetAllCategoriesQuery(
  baseOptions?: Apollo.QueryHookOptions<
    CategoryInputGetAllCategoriesQuery,
    CategoryInputGetAllCategoriesQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    CategoryInputGetAllCategoriesQuery,
    CategoryInputGetAllCategoriesQueryVariables
  >(CategoryInputGetAllCategoriesDocument, options);
}
export function useCategoryInputGetAllCategoriesLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    CategoryInputGetAllCategoriesQuery,
    CategoryInputGetAllCategoriesQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    CategoryInputGetAllCategoriesQuery,
    CategoryInputGetAllCategoriesQueryVariables
  >(CategoryInputGetAllCategoriesDocument, options);
}
export type CategoryInputGetAllCategoriesQueryHookResult = ReturnType<
  typeof useCategoryInputGetAllCategoriesQuery
>;
export type CategoryInputGetAllCategoriesLazyQueryHookResult = ReturnType<
  typeof useCategoryInputGetAllCategoriesLazyQuery
>;
export type CategoryInputGetAllCategoriesQueryResult = Apollo.QueryResult<
  CategoryInputGetAllCategoriesQuery,
  CategoryInputGetAllCategoriesQueryVariables
>;
export const ItemPreviewGetItemDocument = gql`
  query itemPreviewGetItem($id: ID!) {
    getItem(id: $id) {
      id
      name
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
  }
`;

/**
 * __useItemPreviewGetItemQuery__
 *
 * To run a query within a React component, call `useItemPreviewGetItemQuery` and pass it any options that fit your needs.
 * When your component renders, `useItemPreviewGetItemQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useItemPreviewGetItemQuery({
 *   variables: {
 *      id: // value for 'id'
 *   },
 * });
 */
export function useItemPreviewGetItemQuery(
  baseOptions: Apollo.QueryHookOptions<
    ItemPreviewGetItemQuery,
    ItemPreviewGetItemQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useQuery<
    ItemPreviewGetItemQuery,
    ItemPreviewGetItemQueryVariables
  >(ItemPreviewGetItemDocument, options);
}
export function useItemPreviewGetItemLazyQuery(
  baseOptions?: Apollo.LazyQueryHookOptions<
    ItemPreviewGetItemQuery,
    ItemPreviewGetItemQueryVariables
  >
) {
  const options = { ...defaultOptions, ...baseOptions };
  return Apollo.useLazyQuery<
    ItemPreviewGetItemQuery,
    ItemPreviewGetItemQueryVariables
  >(ItemPreviewGetItemDocument, options);
}
export type ItemPreviewGetItemQueryHookResult = ReturnType<
  typeof useItemPreviewGetItemQuery
>;
export type ItemPreviewGetItemLazyQueryHookResult = ReturnType<
  typeof useItemPreviewGetItemLazyQuery
>;
export type ItemPreviewGetItemQueryResult = Apollo.QueryResult<
  ItemPreviewGetItemQuery,
  ItemPreviewGetItemQueryVariables
>;
