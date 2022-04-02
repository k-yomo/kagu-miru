export const routes = {
  home: () => '/media',
  search: () => '/search', // deprecated

  media: () => '/media',
  mediaPost: (slug: string) => `/media/posts/${slug}`,
  mediaCategory: (categoryId: string) => `/media/categories/${categoryId}`,
  mediaTag: (tag: string) => `/media/tags/${tag}`,

  contact: () => '/contact',
  privacyPolicy: () => '/privacy-policy',
};
