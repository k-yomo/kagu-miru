import React, { memo, useEffect, useState } from 'react';
import itemCategories from '@src/static/itemCategories.json';
import { ChevronDownIcon, ChevronUpIcon } from '@heroicons/react/solid';

interface Props {
  // used to sort categories
  displayedItemTopLevelCategoryIds: string[];
  // used to highlight selected category
  selectedCategoryId?: string;
  onClickCategory: (categoryId: string) => void;
  onClearCategory: () => void;
}

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

export default memo(function CategoryList({
  displayedItemTopLevelCategoryIds,
  selectedCategoryId,
  onClickCategory,
  onClearCategory,
}: Props) {
  const [categories, setCategories] = useState(itemCategories);
  const selectedCategoryIdPath = findSelectedCategoryIdPath(
    itemCategories,
    selectedCategoryId
  );
  const [showMore, setShowMore] = useState(false);

  useEffect(() => {
    const categoryIdCountMap = buildCategoryIdCountMap(
      displayedItemTopLevelCategoryIds
    );
    const sortedCategories = itemCategories.sort(function (a, b) {
      const aCount = categoryIdCountMap[a.id] ? categoryIdCountMap[a.id] : 0;
      const bCount = categoryIdCountMap[b.id] ? categoryIdCountMap[b.id] : 0;
      if (aCount > bCount) return -1;
      if (aCount < bCount) return 1;
      return 0;
    });
    setCategories([...sortedCategories]);
  }, [displayedItemTopLevelCategoryIds]);

  return (
    <>
      <span
        className="my-1 cursor-pointer border-text-primary hover:border-b-[1px] text-sm"
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

function Category({
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
    // don't close if already open
    if (!showChildren) {
      setShowChildren(selectedCategoryIdPath.includes(category.id));
    }
  }, [selectedCategoryIdPath, category.id]);

  return (
    <div className="ml-4">
      <div
        className={`flex items-center justify-between px-2 cursor-pointer ${
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
}
