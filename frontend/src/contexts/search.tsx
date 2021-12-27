import {
  createContext,
  Dispatch,
  FC,
  memo,
  PropsWithChildren,
  useContext,
  useEffect,
  useReducer,
  useState,
} from 'react';
import {
  Action,
  EventId,
  ItemColor,
  ItemSellingPlatform,
  SearchFrom,
  SearchSortType,
  SearchDisplayItemsActionParams,
  SearchFilter,
  AppliedMetadata,
  SearchInput,
  SearchQuery,
  useSearchLazyQuery,
  useTrackEventMutation,
} from '@src/generated/graphql';
import { ParsedUrlQuery } from 'querystring';
import { useNextQueryParams } from '@src/lib/nextqueryparams';
import { useRouter } from 'next/router';
import gql from 'graphql-tag';

gql`
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

  mutation trackEvent($event: Event!) {
    trackEvent(event: $event)
  }
`;

export enum SearchActionType {
  CHANGE_QUERY,
  CHANGE_SORT_BY,
  CHANGE_PAGE,
  SET_FILTER,
  SET_CATEGORY_FILTER,
  SET_BRAND_FILTER,
  SET_PLATFORM_FILTER,
  SET_COLOR_FILTER,
  SET_PRICE_FILTER,
  SET_RATING_FILTER,
  SET_METADATA_FILTER,
}

type SearchAction =
  | {
      type: SearchActionType.CHANGE_QUERY;
      payload: { query: string; searchFrom: SearchFrom };
    }
  | { type: SearchActionType.CHANGE_SORT_BY; payload: SearchSortType }
  | { type: SearchActionType.CHANGE_PAGE; payload: number }
  | { type: SearchActionType.SET_FILTER; payload: SearchFilter }
  | { type: SearchActionType.SET_CATEGORY_FILTER; payload: string[] }
  | { type: SearchActionType.SET_BRAND_FILTER; payload: string[] }
  | {
      type: SearchActionType.SET_PLATFORM_FILTER;
      payload: ItemSellingPlatform[];
    }
  | {
      type: SearchActionType.SET_COLOR_FILTER;
      payload: ItemColor[];
    }
  | {
      type: SearchActionType.SET_PRICE_FILTER;
      payload: { minPrice?: number; maxPrice?: number };
    }
  | { type: SearchActionType.SET_RATING_FILTER; payload?: number }
  | { type: SearchActionType.SET_METADATA_FILTER; payload: AppliedMetadata[] };

const searchReducer = (
  state: SearchState,
  action: SearchAction
): SearchState => {
  const { searchInput } = state;
  switch (action.type) {
    case SearchActionType.CHANGE_QUERY:
      return {
        searchInput: {
          ...searchInput,
          query: action.payload.query,
          page: 0,
        },
        searchFrom: action.payload.searchFrom,
      };
    case SearchActionType.CHANGE_SORT_BY:
      return {
        searchInput: {
          ...searchInput,
          sortType: action.payload,
          page: 0,
        },
        searchFrom: SearchFrom.Search,
      };
    case SearchActionType.CHANGE_PAGE:
      return {
        searchInput: {
          ...searchInput,
          page: action.payload,
        },
        searchFrom: SearchFrom.Search,
      };
    case SearchActionType.SET_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: action.payload,
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_CATEGORY_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, categoryIds: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_BRAND_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, brandNames: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_PLATFORM_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, platforms: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_COLOR_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, colors: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_PRICE_FILTER:
      const { minPrice, maxPrice } = action.payload;
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, minPrice, maxPrice },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_RATING_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, minRating: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    case SearchActionType.SET_METADATA_FILTER:
      return {
        searchInput: {
          ...searchInput,
          filter: { ...searchInput.filter, metadata: action.payload },
          page: 0,
        },
        searchFrom: SearchFrom.Filter,
      };
    default:
      return state;
  }
};

export const defaultSearchFilter: SearchFilter = {
  categoryIds: [],
  brandNames: [],
  platforms: [],
  colors: [],
  metadata: [],
};

const SearchContext = createContext<{
  searchState: SearchState;
  searchId: string;
  items: SearchQuery['search']['itemConnection']['nodes'];
  facets: SearchQuery['search']['facets'];
  pageInfo: SearchQuery['search']['itemConnection']['pageInfo'] | undefined;
  dispatch: Dispatch<SearchAction>;
  loading: boolean;
}>({
  searchId: '',
  searchState: {
    searchInput: {
      query: '',
      filter: defaultSearchFilter,
      sortType: SearchSortType.BestMatch,
    },
    searchFrom: SearchFrom.Url,
  },
  items: [],
  facets: [],
  pageInfo: undefined,
  dispatch: () => {},
  loading: false,
});

type SearchState = {
  searchInput: SearchInput;
  searchFrom: SearchFrom;
};

export function queryParamsToSearchParams(
  queryParams: ParsedUrlQuery
): SearchState {
  return {
    searchInput: {
      query: (queryParams.q as string) || '',
      filter: {
        categoryIds: ((queryParams.categoryIds as string) || '')
          .split(',')
          .filter((s) => s),
        brandNames: ((queryParams.brandNames as string) || '')
          .split(',')
          .filter((s) => s),
        platforms: ((queryParams.platforms as string) || '')
          .split(',')
          .flatMap((s: string) => {
            if (Object.values<string>(ItemSellingPlatform).includes(s)) {
              return s as ItemSellingPlatform;
            } else {
              return [];
            }
          }),
        colors: ((queryParams.colors as string) || '')
          .split(',')
          .flatMap((s: string) => {
            if (Object.values<string>(ItemColor).includes(s)) {
              return s as ItemColor;
            } else {
              return [];
            }
          }),
        metadata: Object.keys(queryParams)
          .filter((key) => key.startsWith('m:'))
          .map((key) => {
            const metadataName = key.slice(2);
            return {
              name: metadataName,
              values: (queryParams[key] as string).split(',').filter((s) => s),
            };
          }),
        minPrice: parseInt(queryParams.minPrice as string) || undefined,
        maxPrice: parseInt(queryParams.maxPrice as string) || undefined,
        minRating: parseInt(queryParams.minRating as string) || undefined,
      },
      sortType:
        (queryParams.sort as SearchSortType) || SearchSortType.BestMatch,
      page: parseInt(queryParams.page as string) || 1,
    },
    searchFrom: (queryParams.searchFrom as SearchFrom) || SearchFrom.Url,
  };
}

export function buildSearchUrlQuery(
  query: string,
  filter: SearchFilter,
  searchFrom: SearchFrom,
  sortType?: SearchSortType,
  page?: number
) {
  const urlQuery: { [key: string]: string } = {
    q: query,
    searchFrom: searchFrom.toString(),
  };
  if (filter.categoryIds.length > 0)
    urlQuery.categoryIds = filter.categoryIds.join(',');
  if (filter.brandNames.length > 0)
    urlQuery.brandNames = filter.brandNames.join(',');
  if (filter.platforms.length > 0)
    urlQuery.platforms = filter.platforms.join(',');
  if (filter.colors.length > 0) urlQuery.colors = filter.colors.join(',');
  if (filter.minPrice) urlQuery.minPrice = filter.minPrice.toString();
  if (filter.maxPrice) urlQuery.maxPrice = filter.maxPrice.toString();
  if (filter.minRating) urlQuery.minRating = filter.minRating.toString();
  if (filter.metadata.length > 0) {
    filter.metadata.forEach((m) => {
      urlQuery[`m:${m.name}`] = m.values.join(',');
    });
  }
  if (sortType && sortType !== SearchSortType.BestMatch)
    urlQuery.sort = sortType;
  if (page && page >= 2) urlQuery.page = page.toString();

  return urlQuery;
}

export const useSearch = () => useContext(SearchContext);

type Props = PropsWithChildren<{ isAdmin?: boolean }>;

export const SearchProvider: FC<Props> = memo(function SearchProvider({
  isAdmin,
  children,
}: Props) {
  const router = useRouter();
  const queryParams = useNextQueryParams();
  const [searchId, setSearchId] = useState<string>('');
  const [items, setItems] = useState<
    SearchQuery['search']['itemConnection']['nodes']
  >([]);
  const [facets, setFacets] = useState<SearchQuery['search']['facets']>([]);
  const [pageInfo, setPageInfo] = useState<
    SearchQuery['search']['itemConnection']['pageInfo'] | undefined
  >();
  const [searchState, dispatch] = useReducer(
    searchReducer,
    queryParamsToSearchParams(queryParams)
  );

  const [trackEvent] = useTrackEventMutation();
  const [search, { loading }] = useSearchLazyQuery({
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache',
    onCompleted: (data) => {
      setSearchId(data.search.searchId);
      setItems(data.search.itemConnection.nodes);
      setFacets(data.search.facets);
      setPageInfo(data.search.itemConnection.pageInfo);

      if (isAdmin) {
        return;
      }
      const params: SearchDisplayItemsActionParams = {
        searchId: data.search.searchId,
        searchInput: searchState.searchInput,
        searchFrom: searchState.searchFrom,
        itemIds: data.search.itemConnection.nodes.map((item) => item.id),
      };
      trackEvent({
        variables: {
          event: {
            id: EventId.Search,
            action: Action.Display,
            createdAt: new Date(),
            params,
          },
        },
      }).catch(() => {
        // do nothing
      });
    },
  });

  useEffect(() => {
    setItems([]);
    setPageInfo(undefined);

    const { searchInput, searchFrom } = searchState;
    const { query, filter, sortType, page } = searchInput;
    search({
      variables: {
        input: {
          ...searchInput,
          query: query.trim(),
        },
      },
    });

    const urlQuery = buildSearchUrlQuery(
      query,
      filter,
      searchFrom,
      sortType,
      page || undefined
    );
    // Exclude searchFrom to track actual searched from, since url can be shared.
    const urlQueryWithoutSearchFrom = { ...urlQuery };
    delete urlQueryWithoutSearchFrom.searchFrom;

    router.push(
      {
        pathname: router.pathname,
        query: urlQuery,
      },
      `${router.pathname}?${new URLSearchParams(
        urlQueryWithoutSearchFrom
      ).toString()}`,
      {
        shallow: true,
        scroll: true,
      }
    );
  }, [searchState]);

  return (
    <>
      <SearchContext.Provider
        value={{
          searchState,
          searchId,
          items,
          facets,
          pageInfo,
          dispatch,
          loading,
        }}
      >
        {children}
      </SearchContext.Provider>
    </>
  );
});
