import React, { memo, useState } from 'react';
import { SearchActionType, useSearch } from '@src/contexts/search';

export default memo(function PriceFilter() {
  const { searchState, dispatch } = useSearch();
  const [minPrice, setMinPrice] = useState(
    searchState.searchInput.filter.minPrice
  );
  const [maxPrice, setMaxPrice] = useState(
    searchState.searchInput.filter.maxPrice
  );

  const onClickClear = () => {
    if (minPrice || maxPrice) {
      dispatch({ type: SearchActionType.SET_PRICE_FILTER, payload: {} });
      setMinPrice(undefined);
      setMaxPrice(undefined);
    }
  };

  const onClickApply = () => {
    dispatch({
      type: SearchActionType.SET_PRICE_FILTER,
      payload: {
        minPrice: minPrice && !isNaN(minPrice) ? minPrice : undefined,
        maxPrice: maxPrice && !isNaN(maxPrice) ? maxPrice : undefined,
      },
    });
  };

  return (
    <div>
      <div className="flex items-center text-xs">
        <div>
          <input
            type="text"
            inputMode="numeric"
            value={minPrice || ''}
            onChange={(e) => setMinPrice(parseInt(e.target.value))}
            className="w-[5rem] mr-1 p-1 bg-white dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
          />
          円
        </div>
        {'　'}〜{'　'}
        <div>
          <input
            type="text"
            inputMode="numeric"
            value={maxPrice || ''}
            onChange={(e) => setMaxPrice(parseInt(e.target.value))}
            className="w-[5rem] mr-1 p-1 bg-white dark:bg-gray-800 border border-gray-700 leading-5 text-gray-900 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none focus:ring-1 focus:ring-black dark:focus:ring-gray-400"
          />
          円
        </div>
      </div>
      <div className="flex items-center justify-end mt-4 space-x-2 text-sm">
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 border border-black dark:border-white shadow-sm text-xs font-medium bg-white hover:bg-gray-50 dark:bg-black dark:hover:bg-gray-800 focus:outline-none"
          onClick={onClickClear}
        >
          クリア
        </button>
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 bg-gradient-to-r from-primary-500 dark:from-primary-600 to-rose-500 dark:to-rose-600 text-white focus:outline-none"
          onClick={onClickApply}
        >
          適用
        </button>
      </div>
    </div>
  );
});
