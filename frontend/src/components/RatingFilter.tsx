import React, { memo } from 'react';
import { SearchActionType, useSearch } from '@src/contexts/search';
import RatingSelect from '@src/components/RatingSelect';

export default memo(function RatingFilter() {
  const { searchState, dispatch } = useSearch();
  const minRating = searchState.searchInput.filter?.minRating;

  const onChangeRating = (rating?: number) => {
    dispatch({ type: SearchActionType.SET_RATING_FILTER, payload: rating });
  };

  return (
    <RatingSelect
      curMinRating={minRating || undefined}
      onChangeRating={onChangeRating}
    />
  );
});
