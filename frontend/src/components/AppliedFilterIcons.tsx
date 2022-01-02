import React, { memo, ReactNode } from 'react';
import { XIcon } from '@heroicons/react/solid';
import {
  defaultSearchFilter,
  SearchActionType,
  useSearch,
} from '@src/contexts/search';
import { findCategoryNameById } from '@src/lib/itemCategories';
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
  let filterIcons: ReactNode[] = [];
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
    filterIcons = [...filterIcons, ...platformFilterIcons];
  }

  if (filter.categoryIds.length > 0) {
    const categoryFilterIcons = filter.categoryIds.map((categoryId) => (
      <FilterIcon
        key={`categoryFilter:${categoryId}`}
        name={findCategoryNameById(categoryId)}
        onClear={() =>
          dispatch({
            type: SearchActionType.SET_CATEGORY_FILTER,
            payload: filter.categoryIds.filter((id) => id !== categoryId),
          })
        }
      />
    ));
    filterIcons = [...filterIcons, ...categoryFilterIcons];
  }

  if (filter.brandNames.length > 0) {
    const brandFilterIcons = filter.brandNames.map((brandName) => (
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
    filterIcons = [...filterIcons, ...brandFilterIcons];
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
    filterIcons = [...filterIcons, ...colorFilterIcons];
  }

  if (filter.minPrice || filter.maxPrice) {
    let name: string;
    if (filter.minPrice && filter.maxPrice) {
      name = `${filter.minPrice.toLocaleString()}円 ~ ${filter.maxPrice.toLocaleString()}円`;
    } else if (filter.minPrice) {
      name = `${filter.minPrice.toLocaleString()}円 ~`;
    } else {
      name = `~ ${filter.maxPrice?.toLocaleString()}円`;
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

  if (filter.metadata) {
    const metadataFilterIcons = filter.metadata.map((appliedMetadata) =>
      appliedMetadata.values.map((appliedValue) => (
        <FilterIcon
          key={`metadataFilter:${appliedMetadata.name}:${appliedValue}`}
          name={`${appliedMetadata.name}: ${appliedValue}`}
          onClear={() =>
            dispatch({
              type: SearchActionType.SET_METADATA_FILTER,
              payload: filter.metadata
                .map((m) => ({
                  name: m.name,
                  values: m.values.filter(
                    (v) => m.name !== appliedMetadata.name || v !== appliedValue
                  ),
                }))
                .filter((m) => m.values.length !== 0),
            })
          }
        />
      ))
    );
    if (metadataFilterIcons.length > 0) {
      filterIcons = [...filterIcons, metadataFilterIcons];
    }
  }

  if (filterIcons.length === 0) {
    return null;
  }

  return (
    <div className="flex items-center space-x-2">
      {/* TODO: Use flex-1 instead of using fixed width */}
      <div className="w-[80vw] sm:w-[90%] space-x-2 overflow-auto whitespace-nowrap">
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
