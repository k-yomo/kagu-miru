import {
  createContext,
  Dispatch,
  FC,
  memo,
  useContext,
  useEffect,
  useReducer,
  useState,
} from 'react';
import {
  Action,
  EventId,
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
      filter: { categoryIds: [] },
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

export const useSearch = () => useContext(SearchContext);

export const SearchProvider: FC = memo((props) => {
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

    const urlQuery: { [key: string]: string } = {
      q: query,
    };
    if (filter.categoryIds.length > 0)
      urlQuery.categoryIds = filter.categoryIds.join(',');
    if (filter.minPrice) urlQuery.minPrice = filter.minPrice.toString();
    if (filter.maxPrice) urlQuery.maxPrice = filter.maxPrice.toString();
    if (filter.minRating) urlQuery.minRating = filter.minRating.toString();
    if (sortType !== SearchSortType.BestMatch) urlQuery.sort = sortType;
    if (page && page >= 2) urlQuery.page = page.toString();

    router.push(
      {
        pathname: router.pathname,
        query: {
          ...urlQuery,
          searchFrom,
        },
      },
      // Exclude searchFrom to track actual searched from, since url can be shared.
      `${router.pathname}?${new URLSearchParams(urlQuery).toString()}`,
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
        {props.children}
      </SearchContext.Provider>
    </>
  );
});
