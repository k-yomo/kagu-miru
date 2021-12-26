import itemCategories from '@src/static/itemCategories.json';

export const findCategoryNameById = (id: string): string => {
  return itemCategoryIdNameMap[id];
};

// findCategoryIdsById finds hierarchical id list from L0 to the given category
export const findCategoryIdsById = (id: string): string[] => {
  return itemCategoryIdsMap[id];
};

const itemCategoryIdNameMap = setItemCategoryIdNameMap(itemCategories, {});

function setItemCategoryIdNameMap(
  categories: typeof itemCategories,
  map: { [key: string]: string }
): { [key: string]: string } {
  if (categories.length === 0) {
    return map;
  }
  categories.forEach((category) => {
    map[category.id] = category.name;
    setItemCategoryIdNameMap(category.children, map);
  });

  return map;
}

const itemCategoryIdsMap = setItemCategoryIdsMap(itemCategories, [], {});

function setItemCategoryIdsMap(
  categories: typeof itemCategories,
  parentIds: string[],
  map: { [key: string]: string[] }
): { [key: string]: string[] } {
  if (categories.length === 0) {
    return map;
  }
  categories.forEach((category) => {
    const categoryIds = [...parentIds, category.id];
    map[category.id] = categoryIds;
    setItemCategoryIdsMap(category.children, categoryIds, map);
  });

  return map;
}
