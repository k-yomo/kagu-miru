import React, { memo } from 'react';
import Rating from '@src/components/Rating';
import { SearchActionType, useSearch } from '@src/contexts/search';

export default memo(function RatingFilter() {
  const { searchState, dispatch } = useSearch();
  const minRating = searchState.searchInput.filter.minRating;

  const onClickRating = (rating: number) => {
    dispatch({ type: SearchActionType.SET_RATING_FILTER, payload: rating });
  };

  const onClickClear = () => {
    if (minRating) {
      dispatch({
        type: SearchActionType.SET_RATING_FILTER,
        payload: undefined,
      });
    }
  };

  return (
    <div>
      <div className="space-y-3 text-xs">
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            minRating === 5 ? 'font-bold' : ''
          }`}
          onClick={() => onClickRating(5)}
        >
          <Rating size={20} rating={5} maxRating={5} />
          以上
        </div>
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            minRating === 4 ? 'font-bold' : ''
          }`}
          onClick={() => onClickRating(4)}
        >
          <Rating size={20} rating={4} maxRating={5} />
          以上
        </div>
        <div
          className={`cursor-pointer flex items-center hover:font-bold ${
            minRating === 3 ? 'font-bold' : ''
          }`}
          onClick={() => onClickRating(3)}
        >
          <Rating size={20} rating={3} maxRating={5} />
          以上
        </div>
      </div>
      <div className="flex items-center justify-end mt-2 text-sm">
        <button
          type="button"
          className="inline-flex items-center px-2.5 py-1.5 border border-black dark:border-white shadow-sm text-xs font-medium bg-white hover:bg-gray-50 dark:bg-black dark:hover:bg-gray-800 focus:outline-none"
          onClick={onClickClear}
        >
          クリア
        </button>
      </div>
    </div>
  );
});
