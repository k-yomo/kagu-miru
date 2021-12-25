import React, { memo } from 'react';
import { XIcon } from '@heroicons/react/solid';
import {
  defaultSearchFilter,
  SearchActionType,
  useSearch,
} from '@src/contexts/search';
import { itemCategoryIdNameMap } from '@src/lib/itemCategories';
import { platFormText } from '@src/conv/platform';
import { colorText } from '@src/conv/color';

export default function AppliedFilterIcons() {
  const { searchState, dispatch } = useSearch();

  const onClickClear = () => {
    dispatch({
      type: SearchActionType.SET_FILTER,
      payload: defaultSearchFilter,
    });
  };

  const filter = searchState.searchInput.filter;
  const filterIcons = [];
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

  if (filter.brandNames.length > 0) {
    const categoryFilterIcons = filter.brandNames.map((brandName) => (
      <FilterIcon
        key={`brandFilter:${brandName}`}
        name={brandName}
        onClear={() =>
          dispatch({
            type: SearchActionType.SET_BRAND_FILTER,
            payload: filter.brandNames.filter((name) => name !== brandName),
          })
        }
      />
    ));
    filterIcons.push(categoryFilterIcons);
  }

  if (filter.colors.length > 0) {
    const colorFilterIcons = filter.colors.map((color) => (
      <FilterIcon
        key={`colorFilter:${color}`}
        name={colorText(color)}
        onClear={() =>
          dispatch({
            type: SearchActionType.SET_COLOR_FILTER,
            payload: filter.colors.filter((c) => c !== color),
          })
        }
      />
    ));
    filterIcons.push(colorFilterIcons);
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

  if (filterIcons.length === 0) {
    return null;
  }

  return (
    <div className="flex items-center space-x-2">
      <div className="w-[80vw] space-x-2 overflow-auto whitespace-nowrap">
        {filterIcons}
      </div>
      {filterIcons.length >= 1 && (
        <span
          className="cursor-pointer text-sm text-rose-500 font-bold"
          onClick={onClickClear}
        >
          クリア
        </span>
      )}
    </div>
  );
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
    <span className="inline-flex items-center my-0.5 px-2.5 py-1.5 rounded bg-gradient-to-r from-pink-500 dark:from-pink-600 to-rose-500 dark:to-rose-600 text-white text-xs focus:outline-none">
      {name}
      <XIcon className="w-3 h-3 ml-2 cursor-pointer" onClick={onClear} />
    </span>
  );
});
