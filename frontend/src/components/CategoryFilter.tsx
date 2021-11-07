import React, { memo, useCallback } from 'react';
import { SearchActionType, useSearch } from '@src/contexts/search';
import CategoryList from '@src/components/CategoryList';

export default memo(function CategoryFilter() {
  const { searchState, dispatch } = useSearch();

  const onClickCategory = useCallback(
    (categoryId: string) => {
      dispatch({
        type: SearchActionType.SET_CATEGORY_FILTER,
        payload: [categoryId],
      });
    },
    [dispatch]
  );

  const onClearCategory = useCallback(() => {
    dispatch({ type: SearchActionType.SET_CATEGORY_FILTER, payload: [] });
  }, [dispatch]);

  return (
    <CategoryList
      categoryIds={searchState.searchInput.filter.categoryIds}
      onClickCategory={onClickCategory}
      onClearCategory={onClearCategory}
    />
  );
});
