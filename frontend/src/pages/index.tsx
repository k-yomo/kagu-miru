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
  HomePageSearchItemsQuery,
  SearchItemsSortType,
  useHomePageSearchItemsLazyQuery,
} from '@src/generated/graphql';
import SEOMeta from '@src/components/SEOMeta';
import Loading from '@src/components/Loading';
import { useRouter } from 'next/router';
import PlatformBadge from '@src/components/PlatformBadge';
import Pagination from '@src/components/Pagination';

gql`
  query homePageSearchItems($input: SearchItemsInput!) {
    searchItems(input: $input) {
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
`;

const Home: NextPage = () => {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [sortType, setSortType] = useState(SearchItemsSortType.BestMatch);
  const [page, setPage] = useState<number>(1);
  const [searchItems, { data, loading, error }] =
    useHomePageSearchItemsLazyQuery({
      fetchPolicy: 'no-cache',
      nextFetchPolicy: 'no-cache',
    });

  const refreshPageWithParams = () => {
    router.push(
      `${router.pathname}?q=${searchQuery}&sort=${sortType}&page=${page}`,
      undefined,
      {
        shallow: true,
      }
    );
  };

  const onChangeSortBy = useCallback(
    (e: ChangeEvent<HTMLSelectElement>) => {
      setSortType(e.target.value as SearchItemsSortType);
    },
    [setSortType]
  );

  const onChangeSearchInput = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      setSearchQuery(e.target.value);
    },
    [setSearchQuery]
  );

  const onSearchKeyPress = useCallback(
    (e: KeyboardEvent<HTMLInputElement>) => {
      if (e.key == 'Enter') {
        e.preventDefault();
        refreshPageWithParams();
      }
    },
    [searchQuery]
  );

  useEffect(() => {
    refreshPageWithParams();
  }, [page, sortType]);

  useEffect(() => {
    const page = parseInt(router.query.page as string) || 1;
    const sortType =
      (router.query.sort as SearchItemsSortType) ||
      SearchItemsSortType.BestMatch;
    setPage(page);
    if (router.query.q) {
      const query = router.query.q as string;
      setSearchQuery(query);
      searchItems({
        variables: {
          input: {
            query,
            sortType,
            page,
          },
        },
      });
    }
  }, [router.query]);

  return (
    <div className="max-w-[1200px] mx-auto">
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で検索出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
      <div className="m-4">
        <h1 className="text-2xl text-black dark:text-white">商品検索</h1>
        <div className="flex flex-col sm:flex-row items-end justify-between gap-2 my-4 w-full">
          <div className="relative flex-1 flex-col md:mr-4 lg:mr-12 w-full text-gray-400  focus-within:text-gray-600">
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
                onChange={onChangeSearchInput}
                onKeyPress={onSearchKeyPress}
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
              value={sortType}
              onChange={onChangeSortBy}
            >
              <option value={SearchItemsSortType.BestMatch}>関連度順</option>
              <option value={SearchItemsSortType.SortByPriceAsc}>
                価格の安い順
              </option>
              <option value={SearchItemsSortType.SortByPriceDesc}>
                価格の高い順
              </option>
            </select>
          </div>
        </div>
        {loading ? <Loading /> : <></>}
        <div className="flex flex-col items-center">
          <div className="relative grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3 md:gap-4 text-sm sm:text-md">
            {data && (
              <ItemList items={data.searchItems.nodes} loading={loading} />
            )}
          </div>
          {data && (
            <div className="my-4 w-full">
              <Pagination
                page={data.searchItems.pageInfo.page}
                totalPage={data.searchItems.pageInfo.totalPage}
                onClickPage={(page) => setPage(page)}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

interface ItemListProps {
  items: HomePageSearchItemsQuery['searchItems']['nodes'];
  loading: boolean;
}

const ItemList = memo(function ItemList({ items, loading }: ItemListProps) {
  if (loading) {
    return <Loading />;
  }

  return (
    <>
      {items.map((item) => (
        <a
          key={item.id}
          href={!!item.affiliateUrl ? item.affiliateUrl : item.url}
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
            />
            <div className="py-0.5 sm:p-2">
              <PlatformBadge platform={item.platform} />
              <h4 className="my-1 break-all line-clamp-2 text-sm sm:text-md">
                {item.name}
              </h4>
              <div className="font-bold">￥{item.price}</div>
            </div>
          </div>
        </a>
      ))}
    </>
  );
});

export default Home;
