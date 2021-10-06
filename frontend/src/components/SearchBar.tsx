import React, {
  ChangeEvent,
  KeyboardEvent,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
import { SearchIcon } from '@heroicons/react/solid';
import QuerySuggestionsDropdown from '@src/components/QuerySuggestionsDropdown';
import {
  Action,
  EventId,
  QuerySuggestionsDisplayActionParams,
  SearchFrom,
  useGetQuerySuggestionsLazyQuery,
  useTrackEventMutation,
} from '@src/generated/graphql';
import { SearchActionType, useSearch } from '@src/contexts/search';
import gql from 'graphql-tag';

gql`
  query getQuerySuggestions($query: String!) {
    getQuerySuggestions(query: $query) {
      query
      suggestedQueries
    }
  }
`;

export default memo(function SearchBar() {
  const { searchState, dispatch, loading } = useSearch();
  const [searchQuery, setSearchQuery] = useState(searchState.searchInput.query);
  const [suggestedQueries, setSuggestedQueries] = useState<string[]>([]);
  const [showQuerySuggestions, setShowQuerySuggestions] = useState(false);
  const [trackEvent] = useTrackEventMutation();
  const [getQuerySuggestions, { data: getQuerySuggestionsData }] =
    useGetQuerySuggestionsLazyQuery({
      fetchPolicy: 'no-cache',
      nextFetchPolicy: 'no-cache',
      onCompleted: (data) => {
        const params: QuerySuggestionsDisplayActionParams = {
          query: data.getQuerySuggestions.query,
          suggestedQueries: data.getQuerySuggestions.suggestedQueries,
        };
        trackEvent({
          variables: {
            event: {
              id: EventId.QuerySuggestions,
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

  const onSubmit = useCallback(
    (query: string, searchFrom: SearchFrom) => () => {
      dispatch({
        type: SearchActionType.CHANGE_QUERY,
        payload: { query, searchFrom },
      });
    },
    [dispatch]
  );

  const onChangeSearchQuery = (e: ChangeEvent<HTMLInputElement>) => {
    const query = e.target.value as string;
    setSearchQuery(query);
    getQuerySuggestions({ variables: { query: query.trim() } });
  };

  const onSearchKeyPress = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key == 'Enter') {
      e.preventDefault();
      setShowQuerySuggestions(false);
      onSubmit(searchQuery, SearchFrom.Search);
    }
  };

  const onClickSuggestedQuery = useCallback(
    (query: string) => {
      setSearchQuery(query);
      setShowQuerySuggestions(false);
      onSubmit(query, SearchFrom.QuerySuggestion);
    },
    [onSubmit]
  );

  useEffect(() => {
    if (getQuerySuggestionsData?.getQuerySuggestions) {
      setSuggestedQueries(
        getQuerySuggestionsData.getQuerySuggestions.suggestedQueries
      );
    }
  }, [getQuerySuggestionsData?.getQuerySuggestions]);

  return (
    <div className="z-10 relative flex-1 flex-col md:mr-4 lg:mr-12 w-full text-gray-400 focus-within:text-gray-600">
      <div className="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center">
        <SearchIcon className="h-5 w-5" aria-hidden="true" />
      </div>
      <form action=".">
        <input
          id="search"
          className="appearance-none lock w-full bg-white py-3 pl-10 pr-3 dark:bg-gray-800 border border-gray-700 rounded-md leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
          placeholder="Search"
          type="search"
          name="search"
          value={searchQuery}
          onChange={onChangeSearchQuery}
          onKeyPress={onSearchKeyPress}
          onFocus={() => setShowQuerySuggestions(true)}
          onBlur={() => {
            setTimeout(() => {
              setShowQuerySuggestions(false);
            }, 100);
          }}
          disabled={loading}
        />
        <QuerySuggestionsDropdown
          show={showQuerySuggestions && suggestedQueries.length > 0}
          suggestedQueries={suggestedQueries}
          onClickQuery={onClickSuggestedQuery}
        />
      </form>
    </div>
  );
});
