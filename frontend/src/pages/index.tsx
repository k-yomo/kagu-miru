import React, { memo } from 'react';
import Image from 'next/image';
import {
  Action,
  EventId,
  SearchClickItemActionParams,
  SearchQuery,
  useTrackEventMutation,
} from '@src/generated/graphql';
import { SearchProvider, useSearch } from '@src/contexts/search';
import SEOMeta from '@src/components/SEOMeta';
import Loading from '@src/components/Loading';
import PlatformBadge from '@src/components/PlatformBadge';
import Pagination from '@src/components/Pagination';
import Rating from '@src/components/Rating';
import SearchBar from '@src/components/SearchBar';
import CategoryFilter from '@src/components/CategoryFilter';
import PriceFilter from '@src/components/PriceFilter';
import RatingFilter from '@src/components/RatingFilter';
import SortTypeSelectBox from '@src/components/SortTypeSelectBox';
import AppliedFilterIcons from '@src/components/AppliedFilterIcons';
import MobileSearchFilterModal from '@src/components/MobileSearchFilterModal';
import PlatformFilter from '@src/components/PlatformFilter';

export default function TopPage() {
  return (
    <SearchProvider>
      <TopPageInner />
    </SearchProvider>
  );
}

const TopPageInner = memo(function TopPageInner() {
  const { searchState, searchId, items, pageInfo, loading } = useSearch();
  const [trackEvent] = useTrackEventMutation();

  const onClickItem = (itemId: string) => {
    const params: SearchClickItemActionParams = {
      searchId: searchId,
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
        description="カグミルはオンラインで買える家具を横断で一括検索・比較出来るサービスです。"
        // img={{ srcPath: TopImg.src }}
      />
      <MobileSearchFilterModal />
      <div className="my-8 mx-2 lg:mx-4 lg:min-w-[300px] hidden md:block">
        <div
          className={`mt-1 p-3 ${
            searchState.searchInput.filter.categoryIds.length > 0
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="my-2 text-md font-bold">カテゴリー</h3>
          <CategoryFilter />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchState.searchInput.filter.platform
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">ECサイト</h3>
          <PlatformFilter />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchState.searchInput.filter.minPrice ||
            searchState.searchInput.filter.maxPrice
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">価格</h3>
          <PriceFilter />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchState.searchInput.filter.minRating
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">レビュー評価</h3>
          <RatingFilter />
        </div>
      </div>
      <div className="flex-1 mx-2">
        <div className="flex flex-col sm:flex-row items-end justify-between my-2 gap-2 w-full">
          <SearchBar />
          <SortTypeSelectBox />
        </div>
        <div className="mb-4">
          <AppliedFilterIcons />
        </div>
        {loading ? <Loading /> : <></>}
        <div className="flex flex-col items-center">
          <div className="relative grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3 md:gap-4 text-sm sm:text-md">
            {items && <ItemList items={items} onClickItem={onClickItem} />}
          </div>
          {pageInfo && (
            <div className="my-4 w-full">
              <Pagination />
            </div>
          )}
        </div>
      </div>
    </div>
  );
});

interface ItemListProps {
  items: SearchQuery['search']['itemConnection']['nodes'];
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
              className="w-20 h-20 rounded-t-md"
              unoptimized
            />
            <div className="py-0.5 sm:p-2">
              <PlatformBadge platform={item.platform} />
              <div className="flex items-center">
                <Rating rating={item.averageRating} maxRating={5} />
                <div className="ml-1 text-xs text-gray-600 dark:text-gray-300">
                  {item.reviewCount}
                </div>
              </div>
              <h4 className="mt-1 break-all line-clamp-2 text-sm sm:text-md">
                {item.name}
              </h4>
              <div className=" text-lg font-bold">￥{item.price}</div>
            </div>
          </div>
        </a>
      ))}
    </>
  );
});
