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
  categoryId: Scalars['ID'];
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

export type SubItemCategoryFragment = {
  id: string;
  level: number;
  name: string;
  parentId?: string | null | undefined;
};

export type CategoryInputGetAllCategoriesQueryVariables = Exact<{
  [key: string]: never;
}>;

export type CategoryInputGetAllCategoriesQuery = {
  getAllItemCategories: Array<{
    id: string;
    level: number;
    name: string;
    parentId?: string | null | undefined;
    children: Array<{
      id: string;
      level: number;
      name: string;
      parentId?: string | null | undefined;
      children: Array<{
        id: string;
        level: number;
        name: string;
        parentId?: string | null | undefined;
        children: Array<{
          id: string;
          level: number;
          name: string;
          parentId?: string | null | undefined;
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
