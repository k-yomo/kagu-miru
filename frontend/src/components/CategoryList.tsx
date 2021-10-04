import React, { memo, useEffect, useState } from 'react';
import itemCategories from '@src/static/itemCategories.json';
import { ChevronUpIcon, ChevronDownIcon } from '@heroicons/react/solid';

interface Props {
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

export default memo(function CategoryList({
  selectedCategoryId,
  onClickCategory,
  onClearCategory,
}: Props) {
  const selectedCategoryIdPath = findSelectedCategoryIdPath(
    itemCategories,
    selectedCategoryId
  );
  const [displayedCategories, setDisplayedCategories] =
    useState(itemCategories);
  // const [showMore, setShowMore] = useState(false)

  // useEffect(() => {
  //   setDisplayedCategories(showMore ? itemCategories : itemCategories.slice(0, 10))
  // }, [showMore])

  return (
    <>
      <span
        className="my-1 cursor-pointer border-text-primary hover:border-b-[1px] text-sm "
        onClick={onClearCategory}
      >
        ALL
      </span>
      <div className="space-y-1.5">
        {displayedCategories.map((category) => (
          <Category
            key={category.id}
            category={category}
            selectedCategoryIdPath={selectedCategoryIdPath}
            onClick={onClickCategory}
          />
        ))}
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
  const isSelected =
    selectedCategoryIdPath.length > 0 &&
    selectedCategoryIdPath[selectedCategoryIdPath.length - 1] === category.id;
  const [showChildren, setShowChildren] = useState(false);

  useEffect(() => {
    // don't close if already open
    if (!showChildren) {
      setShowChildren(selectedCategoryIdPath.includes(category.id));
    }
  }, [selectedCategoryIdPath]);

  return (
    <div className="ml-6">
      <div className="flex items-center cursor-pointer ">
        <span
          className={`hover:underline text-sm ${isSelected ? 'font-bold' : ''}`}
          onClick={() => onClick(category.id)}
        >
          {category.name}
        </span>
        {category.children.length > 0 &&
          (showChildren ? (
            <ChevronUpIcon
              className="h-5 w-5"
              onClick={() => setShowChildren(false)}
            />
          ) : (
            <ChevronDownIcon
              className="h-5 w-5"
              onClick={() => setShowChildren(true)}
            />
          ))}
      </div>
      <div className="mt-1.5 space-y-1.5">
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
