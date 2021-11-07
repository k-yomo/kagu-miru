import React, { memo, useCallback, useState } from 'react';
import { SearchActionType, useSearch } from '@src/contexts/search';
import itemCategories from '@src/static/itemCategories.json';
import CategoryList from '@src/components/CategoryList';

export default memo(function CategoryFilter() {
  const { searchState, items, dispatch } = useSearch();
  const [categories, setCategories] = useState(itemCategories);
  const [showMore, setShowMore] = useState(false);

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
