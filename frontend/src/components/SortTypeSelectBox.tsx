import React, { ChangeEvent, memo } from 'react';
import { SearchSortType } from '@src/generated/graphql';
import { SearchActionType, useSearch } from '@src/contexts/search';

export default memo(function SortTypeSelectBox() {
  const { searchState, dispatch, loading } = useSearch();

  const onChangeSortBy = (e: ChangeEvent<HTMLSelectElement>) => {
    dispatch({
      type: SearchActionType.CHANGE_SORT_BY,
      payload: e.target.value as SearchSortType,
    });
  };

  return (
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
        className="appearance-none mt-1 block w-full pl-3 pr-10 py-2 rounded-md text-base dark:bg-gray-800 border border-gray-700 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
        value={searchState.searchInput.sortType || SearchSortType.BestMatch}
        onChange={onChangeSortBy}
        disabled={loading}
      >
        <option value={SearchSortType.BestMatch}>関連度順</option>
        <option value={SearchSortType.PriceAsc}>価格の安い順</option>
        <option value={SearchSortType.PriceDesc}>価格の高い順</option>
        <option value={SearchSortType.ReviewCount}>レビューの件数順</option>
        <option value={SearchSortType.Rating}>レビューの評価順</option>
      </select>
    </div>
  );
});
