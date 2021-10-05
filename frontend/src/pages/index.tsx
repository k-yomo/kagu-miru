import { ChangeEvent, memo, useCallback, useState } from 'react';
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

export const getServerSideProps: GetServerSideProps = async (ctx) => {
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

function queryParamsToSearchParams(queryParams: ParsedUrlQuery): SearchParams {
  return {
    searchInput: {
      query: (queryParams.q as string) || '',
      categoryIds: ((queryParams.categoryIds as string) || '')
        .split(',')
        .filter((s) => s),
      sortType:
        (queryParams.sort as SearchSortType) || SearchSortType.BestMatch,
      page: parseInt(queryParams.page as string) || 1,
    },
    searchFrom: (queryParams.searchFrom as SearchFrom) || SearchFrom.Url,
  };
}

const Home = memo<Props>(({ itemConnection }: Props) => {
  const router = useRouter();
  const queryParams = useNextQueryParams();
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

  const updateSearchParams = useCallback(
    ({ searchInput, searchFrom }: typeof searchParams) => {
      setSearchParams({ searchInput, searchFrom });
      setItems([]);
      setPageInfo(undefined);

      const { query, categoryIds, sortType, page } = searchInput;
      search({
        variables: {
          input: {
            ...searchInput,
            query: query.trim(),
          },
        },
      });

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
        // Exclude searchFrom to track actual searched from, since url can be shared.
        `${router.pathname}?${new URLSearchParams(urlQuery).toString()}`,
        {
          shallow: true,
        }
      );
    },
    []
  );

  const onChangeSortBy = (e: ChangeEvent<HTMLSelectElement>) => {
    updateSearchParams({
      searchInput: {
        ...searchParams.searchInput,
        sortType: e.target.value as SearchSortType,
      },
      searchFrom: SearchFrom.Search,
    });
  };

  const onClickPage = useCallback(
    (page: number) => {
      updateSearchParams({
        searchInput: { ...searchParams.searchInput, page },
        searchFrom: SearchFrom.Search,
      });
    },
    [updateSearchParams, searchParams]
  );

  const onSubmitQuery = useCallback(
    (query: string, searchFrom: SearchFrom) => {
      updateSearchParams({
        searchInput: { ...searchParams.searchInput, query, page: 1 },
        searchFrom,
      });
    },
    [updateSearchParams]
  );

  const onClickCategory = useCallback(
    (categoryId: string) => {
      updateSearchParams({
        searchInput: { ...searchParams.searchInput, categoryIds: [categoryId] },
        searchFrom: SearchFrom.Search,
      });
    },
    [updateSearchParams]
  );

  const onClearCategory = useCallback(() => {
    updateSearchParams({
      searchInput: { ...searchParams.searchInput, categoryIds: [] },
      searchFrom: SearchFrom.Search,
    });
  }, [updateSearchParams]);

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
          displayedItemTopLevelCategoryIds={
            items ? items.map((item) => item.categoryIds[0]) : []
          }
          selectedCategoryId={
            searchParams.searchInput.categoryIds.length > 0
              ? searchParams.searchInput.categoryIds[0]
              : undefined
          }
          onClickCategory={onClickCategory}
          onClearCategory={onClearCategory}
        />
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
