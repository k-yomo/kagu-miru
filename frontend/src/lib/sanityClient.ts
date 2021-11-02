import client from '@sanity/client';
import imageUrlBuilder from '@sanity/image-url';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';

export const sanityClient = client({
  projectId: 'iwkc43by',
  dataset: 'production',
  useCdn: true,
});

export function buildSanityImageSrc(source: SanityImageSource) {
  return imageUrlBuilder(sanityClient).image(source);
}
