import itemCategories from '@src/static/itemCategories.json';

export const itemCategoryIdNameMap = setItemCategoryIdNameMap(
  itemCategories,
  {}
);

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
