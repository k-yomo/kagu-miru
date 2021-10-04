import {
  ChangeEvent,
  KeyboardEvent,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
import type { NextPage } from 'next';
import Image from 'next/image';
import gql from 'graphql-tag';
import { SearchIcon } from '@heroicons/react/solid';
import {
  Action,
  EventId,
  HomePageSearchQuery,
  QuerySuggestionsDisplayActionParams,
  SearchClickItemActionParams,
  SearchDisplayItemsActionParams,
  SearchFrom,
  SearchInput,
  SearchSortType,
  useHomePageGetQuerySuggestionsLazyQuery,
  useHomePageSearchLazyQuery,
  useHomePageTrackEventMutation,
} from '@src/generated/graphql';
import SEOMeta from '@src/components/SEOMeta';
import Loading from '@src/components/Loading';
import { useRouter } from 'next/router';
import PlatformBadge from '@src/components/PlatformBadge';
import Pagination from '@src/components/Pagination';
import QuerySuggestionsDropdown from '@src/components/QuerySuggestionsDropdown';
import Rating from '@src/components/Rating';
import CategoryList from '@src/components/CategoryList';

gql`
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

  query homePageGetQuerySuggestions($query: String!) {
    getQuerySuggestions(query: $query) {
      query
      suggestedQueries
    }
  }

  #  query homePageGetAllCategories {
  #      getAllCategories {
  #          id
  #          name
  #          children {
  #              id
  #              name
  #              children {
  #                  id
  #                  name
  #                  children {
  #                      id
  #                      name
  #                      # deepest categories are 4 levels
  #                      children {
  #                          id
  #                          name
  #                      }
  #                  }
  #              }
  #          }
  #      }
  #  }

  mutation homePageTrackEvent($event: Event!) {
    trackEvent(event: $event)
  }
`;

const Home: NextPage = () => {
  const router = useRouter();
  const [searchInput, setSearchInput] = useState<{
    input: SearchInput;
    searchFrom: SearchFrom;
  }>({
    input: {
      query: '',
      categoryIds: [],
      sortType: SearchSortType.BestMatch,
      page: 1,
    },
    searchFrom: SearchFrom.Url,
  });
  const [searchQuery, setSearchQuery] = useState('');
  const [suggestedQueries, setSuggestedQueries] = useState<string[]>([]);
  const [showQuerySuggestions, setShowQuerySuggestions] = useState(false);
  const [search, { data, loading, error }] = useHomePageSearchLazyQuery({
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache',
    onCompleted: (data) => {
      const params: SearchDisplayItemsActionParams = {
        searchId: data.search.searchId,
        searchInput: searchInput.input,
        searchFrom: searchInput.searchFrom,
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
  const [trackEvent] = useHomePageTrackEventMutation();
  const items = data?.search?.itemConnection?.nodes;
  const pageInfo = data?.search?.itemConnection?.pageInfo;
  const [getQuerySuggestions, { data: getQuerySuggestionsData }] =
    useHomePageGetQuerySuggestionsLazyQuery({
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

  const updateSearchInput = useCallback(
    ({ input, searchFrom }: typeof searchInput) => {
      const { query, categoryIds, sortType, page } = input;
      const urlQuery = {
        q: query,
        categoryIds: categoryIds.join(','),
        sort: sortType,
        page: page ? page.toString() : '',
      };
      router.push(
        {
          pathname: router.pathname,
          query: {
            ...urlQuery,
            searchFrom,
          },
        },
        `${router.pathname}?${new URLSearchParams(urlQuery).toString()}`,
        {
          shallow: true,
        }
      );
    },
    [router]
  );

  const onChangeSortBy = useCallback(
    (e: ChangeEvent<HTMLSelectElement>) => {
      updateSearchInput({
        input: {
          ...searchInput.input,
          sortType: e.target.value as SearchSortType,
        },
        searchFrom: SearchFrom.Search,
      });
    },
    [searchInput, updateSearchInput]
  );

  const onChangeSearchQuery = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const query = e.target.value as string;
      setSearchQuery(query);
      getQuerySuggestions({ variables: { query: query.trim() } });
    },
    [setSearchQuery, getQuerySuggestions]
  );

  const onClickSuggestedQuery = (query: string) => {
    updateSearchInput({
      input: { ...searchInput.input, query, page: 1 },
      searchFrom: SearchFrom.QuerySuggestion,
    });
  };

  const onClickPage = (page: number) => {
    updateSearchInput({
      input: { ...searchInput.input, page },
      searchFrom: SearchFrom.Search,
    });
  };

  const onClickCategory = (categoryId: string) => {
    updateSearchInput({
      input: { ...searchInput.input, categoryIds: [categoryId] },
      searchFrom: SearchFrom.Search,
    });
  };

  const onClearCategory = () => {
    updateSearchInput({
      input: { ...searchInput.input, categoryIds: [] },
      searchFrom: SearchFrom.Search,
    });
  };

  const onSearchKeyPress = useCallback(
    (e: KeyboardEvent<HTMLInputElement>) => {
      if (e.key == 'Enter') {
        e.preventDefault();
        setShowQuerySuggestions(false);
        updateSearchInput({
          input: { ...searchInput.input, query: searchQuery, page: 1 },
          searchFrom: SearchFrom.Search,
        });
      }
    },
    [searchInput, searchQuery, updateSearchInput]
  );

  const onClickItem = (itemId: string) => {
    const params: SearchClickItemActionParams = {
      searchId: data!.search.searchId,
      itemId,
    };
    trackEvent({
      variables: {
        event: {
          id: EventId.Search,
          action: Action.ClickItem,
          createdAt: new Date(),
          params,
        },
      },
    }).catch(() => {
      // do nothing
    });
  };

  useEffect(() => {
    const { query } = searchInput.input;
    if (!query) {
      return;
    }
    search({
      variables: {
        input: {
          ...searchInput.input,
          query: query.trim(),
        },
      },
    });
  }, [searchInput, search]);

  useEffect(() => {
    const categoryIds = ((router.query.categoryIds as string) || '')
      .split(',')
      .filter((s) => s); // remove empty values
    const page = parseInt(router.query.page as string) || 1;
    const sortType =
      (router.query.sort as SearchSortType) || SearchSortType.BestMatch;
    const searchFrom =
      (router.query.searchFrom as SearchFrom) || SearchFrom.Url;
    if (router.query.q) {
      const query = router.query.q as string;
      setSearchQuery(query);
      setSearchInput({
        input: {
          query,
          categoryIds,
          sortType,
          page,
        },
        searchFrom,
      });
    }
  }, [router.query]);

  useEffect(() => {
    if (getQuerySuggestionsData?.getQuerySuggestions) {
      setSuggestedQueries(
        getQuerySuggestionsData.getQuerySuggestions.suggestedQueries
      );
    }
  }, [getQuerySuggestionsData?.getQuerySuggestions]);

  return (
    <div className="flex max-w-[1200px] mx-auto my-3">
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で検索出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
      <div className="my-8 mx-2 lg:mx-4 lg:min-w-[300px] hidden md:block">
        <h2 className="my-2 text-md font-bold">カテゴリー</h2>
        <CategoryList
          selectedCategoryId={
            searchInput.input.categoryIds.length > 0
              ? searchInput.input.categoryIds[0]
              : undefined
          }
          onClickCategory={onClickCategory}
          onClearCategory={onClearCategory}
        />
      </div>
      <div className="flex-1 mx-2">
        <div className="flex flex-col sm:flex-row items-end justify-between my-4 gap-2 w-full">
          <div className="z-10 relative flex-1 flex-col md:mr-4 lg:mr-12 w-full text-gray-400  focus-within:text-gray-600">
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
          <div>
            <label
              htmlFor="location"
              className="block text-sm font-medium text-gray-700 dark:text-gray-200"
            >
              並び替え
            </label>
            <select
              id="location"
              name="location"
              className="appearance-none mt-1 block w-full pl-3 pr-10 py-2 rounded-md text-base border border-gray-700 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
              value={searchInput.input.sortType}
              onChange={onChangeSortBy}
              disabled={loading}
            >
              <option value={SearchSortType.BestMatch}>関連度順</option>
              <option value={SearchSortType.PriceAsc}>価格の安い順</option>
              <option value={SearchSortType.PriceDesc}>価格の高い順</option>
              <option value={SearchSortType.ReviewCount}>
                レビューの件数順
              </option>
              <option value={SearchSortType.Rating}>レビューの評価順</option>
            </select>
          </div>
        </div>
        {loading ? <Loading /> : <></>}
        <div className="flex flex-col items-center">
          <div className="relative grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3 md:gap-4 text-sm sm:text-md">
            {items && (
              <ItemList
                items={items}
                loading={loading}
                onClickItem={onClickItem}
              />
            )}
          </div>
          {pageInfo && (
            <div className="my-4 w-full">
              <Pagination
                page={pageInfo.page}
                totalPage={pageInfo.totalPage}
                onClickPage={onClickPage}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

interface ItemListProps {
  items: HomePageSearchQuery['search']['itemConnection']['nodes'];
  onClickItem: (itemId: string) => void;
  loading: boolean;
}

const ItemList = memo(function ItemList({
  items,
  loading,
  onClickItem,
}: ItemListProps) {
  if (loading) {
    return <Loading />;
  }

  return (
    <>
      {items.map((item) => (
        <a
          key={item.id}
          href={!!item.affiliateUrl ? item.affiliateUrl : item.url}
          onClick={() => onClickItem(item.id)}
        >
          <div className="rounded-md sm:shadow">
            <Image
              src={item.imageUrls[0] || 'https://via.placeholder.com/300'}
              alt={item.name}
              width={300}
              height={300}
              layout="responsive"
              objectFit="cover"
              className="w-20 h-20"
              unoptimized
            />
            <div className="py-0.5 sm:p-2">
              <PlatformBadge platform={item.platform} />
              <h4 className="my-1 break-all line-clamp-2 text-sm sm:text-md">
                {item.name}
              </h4>
              <div className=" my-1 text-lg font-bold">￥{item.price}</div>
              <div className="flex items-center">
                <Rating rating={item.averageRating} maxRating={5} />
                <div className="ml-1 text-xs text-gray-600 dark:text-gray-300">
                  {item.reviewCount}
                </div>
              </div>
            </div>
          </div>
        </a>
      ))}
    </>
  );
});

export default Home;
