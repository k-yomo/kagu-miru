import React, { memo, useCallback } from 'react';
import { useRouter } from 'next/router';
import { useClipboard } from 'use-clipboard-copy';
import {
  Action,
  EventId,
  SearchClickItemActionParams,
  SearchFrom,
  useTrackEventMutation,
} from '@src/generated/graphql';
import {
  SearchActionType,
  SearchProvider,
  useSearch,
} from '@src/contexts/search';
import { useToast } from '@src/contexts/toast';
import SEOMeta from '@src/components/SEOMeta';
import Loading from '@src/components/Loading';
import Pagination from '@src/components/Pagination';
import SearchBar from '@src/components/SearchBar';
import CategoryFilter from '@src/components/CategoryFilter';
import ColorFilter from '@src/components/ColorFilter';
import PriceFilter from '@src/components/PriceFilter';
import RatingFilter from '@src/components/RatingFilter';
import SortTypeSelectBox from '@src/components/SortTypeSelectBox';
import AppliedFilterIcons from '@src/components/AppliedFilterIcons';
import MobileSearchFilterModal from '@src/components/MobileSearchFilterModal';
import PlatformFilter from '@src/components/PlatformFilter';
import SearchPageScreenImg from '@public/images/search_screen.jpeg';
import ItemList from '@src/components/ItemList';
import Facets from '@src/components/Facets';

export default function TopPage() {
  const router = useRouter();
  return (
    <>
      <SEOMeta
        title="カグミル - 家具検索サービス"
        excludeSiteTitle
        description="カグミルはオンラインで買える家具を横断で一括検索・比較出来るサービスです。"
        img={{ srcPath: SearchPageScreenImg.src }}
        path={router.asPath}
      />
      <SearchProvider>
        <TopPageInner />
      </SearchProvider>
    </>
  );
}

export const TopPageInner = memo(function TopPageInner({
  isAdmin,
}: {
  isAdmin?: boolean;
}) {
  const { searchState, searchId, items, pageInfo, loading, dispatch } =
    useSearch();
  const toast = useToast();
  const clipboard = useClipboard();
  const [trackEvent] = useTrackEventMutation();

  const onSubmitQuery = useCallback(
    (query: string, searchFrom: SearchFrom) => {
      dispatch({
        type: SearchActionType.CHANGE_QUERY,
        payload: { query, searchFrom },
      });
    },
    [dispatch]
  );

  const onClickItem = useCallback(
    (itemId: string) => {
      if (isAdmin) {
        clipboard.copy(itemId);
        toast(`商品ID ${itemId} をコピーしました`, { type: 'success' });
        return;
      }
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
    },
    [searchId]
  );

  return (
    <div className="flex max-w-[1200px] mx-auto mt-3 mb-6">
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
            searchState.searchInput.filter.platforms.length > 0
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">ECサイト</h3>
          <PlatformFilter />
        </div>
        <div
          className={`mt-6 p-3 ${
            searchState.searchInput.filter.colors.length > 0
              ? 'bg-gray-100 dark:bg-gray-800'
              : ''
          }`}
        >
          <h3 className="mb-2 text-md font-bold">カラー</h3>
          <ColorFilter />
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
          <SearchBar
            query={searchState.searchInput.query}
            onSubmit={onSubmitQuery}
          />
          <SortTypeSelectBox />
        </div>
        <div className="sticky top-0 z-10 py-2 space-y-2 bg-white dark:bg-black">
          <Facets />
          <AppliedFilterIcons />
        </div>
        {loading ? <Loading /> : <></>}
        {pageInfo && (
          <div className="my-2 text-sm">
            検索結果: {pageInfo.totalCount.toLocaleString()}件{' '}
            {pageInfo.totalCount === 10_000 && '以上'}
          </div>
        )}
        <div className="flex flex-col items-center">
          <div className="relative grid grid-cols-2 sm:grid-cols-4 md:grid-cols-5 gap-3 md:gap-4 text-sm sm:text-md">
            {items.length > 0 && (
              <ItemList
                isAdmin={isAdmin!!}
                items={items}
                onClickItem={onClickItem}
              />
            )}
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
