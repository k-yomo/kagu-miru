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
  SearchDisplayItemsActionParams,
  SearchFilter,
  SearchFrom,
  SearchInput,
  SearchQuery,
  SearchSortType,
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
  SET_PLATFORM_FILTER,
  SET_COLOR_FILTER,
  SET_PRICE_FILTER,
  SET_RATING_FILTER,
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
  | { type: SearchActionType.SET_RATING_FILTER; payload?: number };

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
    default:
      return state;
  }
};

export const defaultSearchFilter: SearchFilter = {
  categoryIds: [],
  platforms: [],
  colors: [],
};

const SearchContext = createContext<{
  searchState: SearchState;
  searchId: string;
  items: SearchQuery['search']['itemConnection']['nodes'];
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
  if (filter.platforms.length > 0)
    urlQuery.platforms = filter.platforms.join(',');
  if (filter.colors.length > 0) urlQuery.colors = filter.colors.join(',');
  if (filter.minPrice) urlQuery.minPrice = filter.minPrice.toString();
  if (filter.maxPrice) urlQuery.maxPrice = filter.maxPrice.toString();
  if (filter.minRating) urlQuery.minRating = filter.minRating.toString();
  if (sortType && sortType !== SearchSortType.BestMatch)
    urlQuery.sort = sortType;
  if (page && page >= 2) urlQuery.page = page.toString();

  return urlQuery;
}

export const useSearch = () => useContext(SearchContext);

type Props = PropsWithChildren<{ isAdmin?: boolean }>;

export const SearchProvider: FC<Props> = memo(
  ({ isAdmin, children }: Props) => {
    const router = useRouter();
    const queryParams = useNextQueryParams();
    const [searchId, setSearchId] = useState<string>('');
    const [items, setItems] = useState<
      SearchQuery['search']['itemConnection']['nodes']
    >([]);
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
            pageInfo,
            dispatch,
            loading,
          }}
        >
          {children}
        </SearchContext.Provider>
      </>
    );
  }
);
