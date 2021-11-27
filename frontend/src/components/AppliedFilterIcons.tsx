import React, { memo } from 'react';
import { XIcon } from '@heroicons/react/solid';
import { SearchActionType, useSearch } from '@src/contexts/search';
import { itemCategoryIdNameMap } from '@src/lib/itemCategories';
import { platFormText } from '@src/conv/platform';

export default function AppliedFilterIcons() {
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

  if (filter.platforms.length > 0) {
    const platformFilterIcons = filter.platforms.map((platform) => (
      <FilterIcon
        key={`platformFilter:${platform}`}
        name={platFormText(platform)}
        onClear={() =>
          dispatch({
            type: SearchActionType.SET_PLATFORM_FILTER,
            payload: filter.platforms.filter((p) => p !== platform),
          })
        }
      />
    ));
    filterIcons.push(platformFilterIcons);
  }

  if (filter.minPrice || filter.maxPrice) {
    let name: string;
    if (filter.minPrice && filter.maxPrice) {
      name = `${filter.minPrice}円 ~ ${filter.maxPrice}円`;
    } else if (filter.minPrice) {
      name = `${filter.minPrice}円 ~`;
    } else {
      name = `~ ${filter.maxPrice}円`;
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
    <span className="inline-flex items-center px-2.5 py-1.5 rounded bg-gradient-to-r from-primary-500 dark:from-primary-600 to-rose-500 dark:to-rose-600 text-white text-xs focus:outline-none">
      {name}
      <XIcon className="w-3 h-3 ml-2 cursor-pointer" onClick={onClear} />
    </span>
  );
});
