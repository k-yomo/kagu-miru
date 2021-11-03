import client from '@sanity/client';
import imageUrlBuilder from '@sanity/image-url';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';
import { SANITY_DATASET_ENV } from '@src/config/env';

export const sanityClient = client({
  projectId: 'iwkc43by',
  dataset: SANITY_DATASET_ENV,
  useCdn: true,
});

export function buildSanityImageSrc(source: SanityImageSource) {
  return imageUrlBuilder(sanityClient).image(source);
}
