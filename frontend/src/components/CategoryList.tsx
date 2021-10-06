import React, { memo, useCallback, useEffect, useState } from 'react';
import itemCategories from '@src/static/itemCategories.json';
import { ChevronDownIcon, ChevronUpIcon } from '@heroicons/react/solid';
import { SearchActionType, useSearch } from '@src/contexts/search';

function findSelectedCategoryIdPath(
  categories: typeof itemCategories,
  selectedCategoryId?: string
): string[] {
  if (!selectedCategoryId) {
    return [];
  }

  const children = categories.map((category) => {
    if (category.id === selectedCategoryId) {
      return [category.id];
    }
    const childCategoryIds = findSelectedCategoryIdPath(
      category.children,
      selectedCategoryId
    );
    if (childCategoryIds.length > 0) {
      return [category.id, ...childCategoryIds];
    }
    return [];
  });

  return children.find((categories) => categories.length > 0) || [];
}

function buildCategoryIdCountMap(displayedCategoryIds: string[]): {
  [key: string]: number;
} {
  const categoryIdDisplayCountMap: { [key: string]: number } = {};
  displayedCategoryIds.forEach((id) => {
    if (categoryIdDisplayCountMap[id]) {
      categoryIdDisplayCountMap[id] += 1;
    } else {
      categoryIdDisplayCountMap[id] = 1;
    }
  });
  return categoryIdDisplayCountMap;
}

export default memo(function CategoryList() {
  const { searchState, items, dispatch } = useSearch();
  const [categories, setCategories] = useState(itemCategories);
  const [showMore, setShowMore] = useState(false);
  const selectedCategoryIdPath = findSelectedCategoryIdPath(
    itemCategories,
    searchState.searchInput.filter.categoryIds.length > 0
      ? searchState.searchInput.filter.categoryIds[0]
      : undefined
  );

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

  useEffect(() => {
    const categoryIdCountMap = buildCategoryIdCountMap(
      items.map((item) => item.categoryIds[0])
    );
    const sortedCategories = itemCategories.sort(function (a, b) {
      const aCount = categoryIdCountMap[a.id] ? categoryIdCountMap[a.id] : 0;
      const bCount = categoryIdCountMap[b.id] ? categoryIdCountMap[b.id] : 0;
      if (aCount > bCount) return -1;
      if (aCount < bCount) return 1;
      return 0;
    });
    setCategories([...sortedCategories]);
  }, [items]);

  return (
    <>
      <span
        className="cursor-pointer hover:underline border-text-primary hover:border-b-[1px] text-sm"
        onClick={onClearCategory}
      >
        ALL
      </span>
      <div>
        {(showMore ? categories : categories.slice(0, 10)).map((category) => (
          <Category
            key={category.id}
            category={category}
            selectedCategoryIdPath={selectedCategoryIdPath}
            onClick={onClickCategory}
          />
        ))}
        <div
          className="cursor-pointer hover:underline text-sm"
          onClick={() => setShowMore(!showMore)}
        >
          {showMore ? (
            <div className="flex items-center">
              表示を少なく
              <ChevronUpIcon className="h-5 w-5" />
            </div>
          ) : (
            <div className="flex items-center">
              続きを見る
              <ChevronDownIcon className="h-5 w-5" />
            </div>
          )}
        </div>
      </div>
    </>
  );
});

interface CategoryProps {
  category: typeof itemCategories[number];
  selectedCategoryIdPath: string[];
  onClick: (categoryId: string) => void;
}

const Category = memo(function Category({
  category,
  selectedCategoryIdPath,
  onClick,
}: CategoryProps) {
  const hasChildren = category.children.length > 0;
  const isSelected =
    selectedCategoryIdPath.length > 0 &&
    selectedCategoryIdPath[selectedCategoryIdPath.length - 1] === category.id;
  const [showChildren, setShowChildren] = useState(false);

  useEffect(() => {
    setShowChildren((prevState) => {
      // don't close if already open
      if (prevState) return prevState;
      return selectedCategoryIdPath.includes(category.id);
    });
  }, [setShowChildren, selectedCategoryIdPath, category.id]);

  return (
    <div className="ml-2">
      <div
        className={`flex items-center justify-between cursor-pointer ${
          hasChildren ? 'hover:bg-gray-50 dark:hover:bg-gray-800' : 'py-2'
        }`}
      >
        <span
          className={`hover:underline text-sm ${isSelected ? 'font-bold' : ''}`}
          onClick={() => onClick(category.id)}
        >
          {category.name}
        </span>
        <div
          className="flex-1 flex justify-end py-2"
          onClick={() => setShowChildren(!showChildren)}
        >
          {hasChildren &&
            (showChildren ? (
              <ChevronUpIcon className="h-5 w-5" />
            ) : (
              <ChevronDownIcon className="h-5 w-5" />
            ))}
        </div>
      </div>
      <div>
        {showChildren &&
          category.children.map((category) => (
            <Category
              key={category.id}
              category={category}
              selectedCategoryIdPath={selectedCategoryIdPath}
              onClick={onClick}
            />
          ))}
      </div>
    </div>
  );
});
