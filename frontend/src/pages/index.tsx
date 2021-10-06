import React, {
  ChangeEvent,
  memo,
  useCallback,
  useEffect,
  useState,
} from 'react';
import type { GetServerSideProps } from 'next';
import Image from 'next/image';
import gql from 'graphql-tag';
import {
  Action,
  EventId,
  HomePageSearchDocument,
  HomePageSearchQuery,
  HomePageTrackEventDocument,
  HomePageTrackEventMutation,
  SearchClickItemActionParams,
  SearchDisplayItemsActionParams,
  SearchFrom,
  SearchInput,
  SearchSortType,
  useHomePageSearchLazyQuery,
  useHomePageTrackEventMutation,
} from '@src/generated/graphql';
import SEOMeta from '@src/components/SEOMeta';
import Loading from '@src/components/Loading';
import { useRouter } from 'next/router';
import PlatformBadge from '@src/components/PlatformBadge';
import Pagination from '@src/components/Pagination';
import Rating from '@src/components/Rating';
import CategoryList from '@src/components/CategoryList';
import SearchBar from '@src/components/SearchBar';
import { useNextQueryParams } from '@src/lib/nextqueryparams';
import { ParsedUrlQuery } from 'querystring';
import apolloClient from '@src/lib/apolloClient';
import PriceFilter from '@src/components/PriceFilter';
import RatingFilter from '@src/components/RatingFilter';

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
          categoryIds
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

type SearchParams = {
  searchInput: SearchInput;
  searchFrom: SearchFrom;
};

function queryParamsToSearchParams(queryParams: ParsedUrlQuery): SearchParams {
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

export const getServerSideProps: GetServerSideProps = async (ctx) => {
  ctx.res.setHeader('Cache-Control', 'public, max-age=3600, stale-while-revalidate=60');
  const searchParams = queryParamsToSearchParams(ctx.query);
  const { data, errors } = await apolloClient.query<HomePageSearchQuery>({
    query: HomePageSearchDocument,
    variables: {
      input: searchParams.searchInput,
    },
  });

  if (errors) {
    throw errors;
  }

  const params: SearchDisplayItemsActionParams = {
    searchId: data.search.searchId,
    searchInput: searchParams.searchInput,
    searchFrom: searchParams.searchFrom,
    itemIds: data.search.itemConnection.nodes.map((item) => item.id),
  };

  apolloClient
    .mutate<HomePageTrackEventMutation>({
      mutation: HomePageTrackEventDocument,
      variables: {
        event: {
          id: EventId.Search,
          action: Action.Display,
          createdAt: new Date(),
          params,
        },
      },
    })
    .catch(() => {
      // do nothing
    });

  return {
    props: {
      itemConnection: data.search.itemConnection,
    },
  };
};

interface Props {
  itemConnection: NonNullable<HomePageSearchQuery['search']['itemConnection']>;
}

const Home = memo<Props>(({ itemConnection }: Props) => {
  const router = useRouter();
  const queryParams = useNextQueryParams();
  const [isMount, setIsMount] = useState(false);
  const [items, setItems] = useState(itemConnection.nodes);
  const [pageInfo, setPageInfo] = useState<
    typeof itemConnection.pageInfo | undefined
  >(itemConnection.pageInfo);
  const [searchParams, setSearchParams] = useState<{
    searchInput: SearchInput;
    searchFrom: SearchFrom;
  }>(queryParamsToSearchParams(queryParams));
  const [trackEvent] = useHomePageTrackEventMutation();
  const [search, { data, loading, error }] = useHomePageSearchLazyQuery({
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache',
    onCompleted: (data) => {
      setItems(data.search.itemConnection.nodes);
      setPageInfo(data.search.itemConnection.pageInfo);
      const params: SearchDisplayItemsActionParams = {
        searchId: data.search.searchId,
        searchInput: searchParams.searchInput,
        searchFrom: searchParams.searchFrom,
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

  const onChangeSortBy = (e: ChangeEvent<HTMLSelectElement>) => {
    setSearchParams(({ searchInput }) => ({
      searchInput: {
        ...searchInput,
        sortType: e.target.value as SearchSortType,
      },
      searchFrom: SearchFrom.Search,
    }));
  };

  const onClickPage = useCallback(
    (page: number) => {
      setSearchParams(({ searchInput }) => ({
        searchInput: { ...searchInput, page },
        searchFrom: SearchFrom.Search,
      }));
    },
    [setSearchParams]
  );

  const onSubmitQuery = useCallback(
    (query: string, searchFrom: SearchFrom) => {
      setSearchParams(({ searchInput }) => ({
        searchInput: { ...searchInput, query, page: 1 },
        searchFrom,
      }));
    },
    [setSearchParams]
  );

  const onClickCategory = useCallback(
    (categoryId: string) => {
      setSearchParams(({ searchInput }) => {
        const { filter, ...rest } = searchInput;
        return {
          searchInput: {
            ...rest,
            filter: { ...filter, categoryIds: [categoryId] },
          },
          searchFrom: SearchFrom.Search,
        };
      });
    },
    [setSearchParams]
  );

  const onClearCategory = useCallback(() => {
    setSearchParams(({ searchInput }) => {
      const { filter, ...rest } = searchInput;
      return {
        searchInput: { ...rest, filter: { ...filter, categoryIds: [] } },
        searchFrom: SearchFrom.Search,
      };
    });
  }, [setSearchParams]);

  const onSubmitPriceFilter = useCallback(
    (minPrice?: number, maxPrice?: number) => {
      setSearchParams(({ searchInput }) => {
        const { filter, ...rest } = searchInput;
        return {
          searchInput: {
            ...rest,
            filter: { ...filter, minPrice, maxPrice },
          },
          searchFrom: SearchFrom.Search,
        };
      });
    },
    [setSearchParams]
  );

  const onClearPriceFilter = useCallback(() => {
    setSearchParams(({ searchInput }) => {
      const { filter, ...rest } = searchInput;
      return {
        searchInput: {
          ...rest,
          filter: { ...filter, minPrice: undefined, maxPrice: undefined },
        },
        searchFrom: SearchFrom.Search,
      };
    });
  }, [setSearchParams]);

  const onSubmitRatingFilter = useCallback(
    (minRating?: number) => {
      setSearchParams(({ searchInput }) => {
        const { filter, ...rest } = searchInput;
        return {
          searchInput: {
            ...rest,
            filter: { ...filter, minRating },
          },
          searchFrom: SearchFrom.Search,
        };
      });
    },
    [setSearchParams]
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
    if (!isMount) {
      setIsMount(true);
      return;
    }
    console.log('fired');
    setItems([]);
    setPageInfo(undefined);

    const { searchInput, searchFrom } = searchParams;
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
      }
    );
  }, [searchParams]);

  return (
    <div className="flex max-w-[1200px] mx-auto my-3">
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で検索出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
      <div className="my-8 mx-2 lg:mx-4 lg:min-w-[300px] hidden md:block">
        <div
          className={`mt-1 p-3 ${
            searchParams.searchInput.filter.categoryIds.length > 0
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="my-2 text-md font-bold">カテゴリー</h3>
          <CategoryList
            displayedItemTopLevelCategoryIds={
              items ? items.map((item) => item.categoryIds[0]) : []
            }
            selectedCategoryId={
              searchParams.searchInput.filter.categoryIds.length > 0
                ? searchParams.searchInput.filter.categoryIds[0]
                : undefined
            }
            onClickCategory={onClickCategory}
            onClearCategory={onClearCategory}
          />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchParams.searchInput.filter.minPrice ||
            searchParams.searchInput.filter.maxPrice
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">価格</h3>
          <PriceFilter
            defaultMinPrice={
              searchParams.searchInput.filter.minPrice || undefined
            }
            defaultMaxPrice={
              searchParams.searchInput.filter.maxPrice || undefined
            }
            onSubmit={onSubmitPriceFilter}
            onClear={onClearPriceFilter}
          />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchParams.searchInput.filter.minRating
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">レビュー評価</h3>
          <RatingFilter
            minRating={searchParams.searchInput.filter.minRating || undefined}
            onSubmit={onSubmitRatingFilter}
          />
        </div>
      </div>
      <div className="flex-1 mx-2">
        <div className="flex flex-col sm:flex-row items-end justify-between my-4 gap-2 w-full">
          <SearchBar
            defaultQuery={searchParams.searchInput.query}
            loading={loading}
            onSubmit={onSubmitQuery}
          />
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
              value={searchParams.searchInput.sortType}
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
            {items && <ItemList items={items} onClickItem={onClickItem} />}
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
});

interface ItemListProps {
  items: HomePageSearchQuery['search']['itemConnection']['nodes'];
  onClickItem: (itemId: string) => void;
}

const ItemList = memo(function ItemList({ items, onClickItem }: ItemListProps) {
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
              className="w-20 h-20 rounded-md"
              unoptimized
            />
            <div className="py-0.5 sm:p-2">
              <div className="flex justify-between">
                <PlatformBadge platform={item.platform} />
                <div className="flex items-center">
                  <Rating rating={item.averageRating} maxRating={5} />
                  <div className="ml-1 text-xs text-gray-600 dark:text-gray-300">
                    {item.reviewCount}
                  </div>
                </div>
              </div>
              <h4 className="my-1 break-all line-clamp-2 text-sm sm:text-md">
                {item.name}
              </h4>
              <div className=" my-1 text-lg font-bold">￥{item.price}</div>
            </div>
          </div>
        </a>
      ))}
    </>
  );
});

export default Home;
