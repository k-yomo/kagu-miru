import React, { memo } from 'react';
import { XIcon } from '@heroicons/react/solid';
import { SearchActionType, useSearch } from '@src/contexts/search';
import { itemCategoryIdNameMap } from '@src/lib/itemCategories';

export default function FilterIcons() {
  const { searchState, dispatch } = useSearch();
  const filter = searchState.searchInput.filter;

  const filterIcons = [];
  if (filter.categoryIds.length > 0) {
    const categoryFilterIcons = filter.categoryIds.map((categoryId) => (
      <FilterIcon
        key={`categoryFilter:${categoryId}`}
        name={itemCategoryIdNameMap[categoryId]}
        onClear={() =>
          dispatch({
            type: SearchActionType.SET_CATEGORY_FILTER,
            payload: filter.categoryIds.filter((id) => id !== categoryId),
          })
        }
      />
    ));
    filterIcons.push(categoryFilterIcons);
  }

  if (filter.minPrice || filter.maxPrice) {
    let name: string;
    if (filter.minPrice && filter.maxPrice) {
      name = `${filter.minPrice}円 ~ ${filter.maxPrice}円`;
    } else if (filter.minPrice) {
      name = `${filter.minPrice}円 ~`;
    } else {
      name = `~ ${filter.minPrice}円`;
    }
    filterIcons.push(
      <FilterIcon
        key="priceFilter"
        name={name}
        onClear={() => {
          dispatch({ type: SearchActionType.SET_PRICE_FILTER, payload: {} });
        }}
      />
    );
  }

  if (filter.minRating) {
    filterIcons.push(
      <FilterIcon
        key="ratingFilter"
        name={`評価${filter.minRating}以上`}
        onClear={() => {
          dispatch({
            type: SearchActionType.SET_RATING_FILTER,
            payload: undefined,
          });
        }}
      />
    );
  }

  return <div className="space-x-2">{filterIcons}</div>;
}

interface FilterIconProps {
  name: string;
  onClear: () => void;
}

const FilterIcon = memo(function FilterIcon({
  name,
  onClear,
}: FilterIconProps) {
  return (
    <span className="inline-flex items-center px-2.5 py-1.5 rounded bg-indigo-100 dark:bg-indigo-600 text-indigo-800 dark:text-indigo-100 text-xs focus:outline-none">
      {name}
      <XIcon className="w-3 h-3 ml-2 cursor-pointer" onClick={onClear} />
    </span>
  );
});
