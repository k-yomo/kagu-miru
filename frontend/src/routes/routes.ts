export const routes = {
  top: () => '/',
  media: () => '/media',
  post: (slug: string) => `/media/posts/${slug}`,
};
