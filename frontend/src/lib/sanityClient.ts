import client from '@sanity/client';
import imageUrlBuilder from '@sanity/image-url';
import { SanityImageSource } from '@sanity/image-url/lib/types/types';
import { SANITY_DATASET_ENV } from '@src/config/env';

const projectId = 'iwkc43by';

export const sanityClient = client({
  projectId: projectId,
  dataset: SANITY_DATASET_ENV,
  useCdn: true,
});

export const sanityPreviewClient = client({
  projectId: projectId,
  dataset: SANITY_DATASET_ENV,
  useCdn: false,
  withCredentials: true,
});

export function buildSanityImageSrc(source: SanityImageSource) {
  return imageUrlBuilder(sanityClient)
    .image(source)
    .auto('format')
    .maxWidth(1000)
    .quality(50);
}
