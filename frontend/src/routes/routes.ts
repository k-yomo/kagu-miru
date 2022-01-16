export const routes = {
  home: () => '/',
  search: () => '/search',

  itemDetail: (itemId: string) => `/items/${itemId}`,

  media: () => '/media',
  mediaPost: (slug: string) => `/media/posts/${slug}`,
  mediaCategory: (categoryId: string) => `/media/categories/${categoryId}`,
  mediaTag: (tag: string) => `/media/tags/${tag}`,

  contact: () => '/contact',
  privacyPolicy: () => '/privacy-policy',
};
