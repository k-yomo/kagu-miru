import React from 'react';
import { GetServerSideProps } from 'next';
import groq from 'groq';
import { sanityClient } from '@src/lib/sanityClient';
import { SITE_ROOT_URL } from '@src/config/constant';
import { routes } from '@src/routes/routes';

export default function SiteMap() {
  return <></>;
}

export const getServerSideProps: GetServerSideProps = async ({ res }) => {
  const sitemapUrls = [
    routes.media(),
    routes.contact(),
    routes.privacyPolicy(),
  ].map(
    (path) => `
  <url>
    <loc>${SITE_ROOT_URL}${path}</loc>
    <changefreq>daily</changefreq>
    <priority>1</priority>
  </url>
`
  );

  const query = groq`{
      "posts": *[_type == 'post']{slug,_updatedAt},
    }`;
  const { posts } = await sanityClient.fetch(query);
  sitemapUrls.push(
    ...posts.map(
      ({
        _updatedAt,
        slug,
      }: {
        _updatedAt: string;
        slug: { current: string };
      }) => `
  <url>
    <loc>${SITE_ROOT_URL}${routes.mediaPost(slug.current)}</loc>
    <lastmod>${_updatedAt}</lastmod>
    <changefreq>daily</changefreq>
    <priority>0.5</priority>
  </url>
`
    )
  );

  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">${sitemapUrls.join(
    ''
  )}</urlset>
`;
  res.setHeader('Content-Type', 'text/xml');
  res.write(sitemap);
  res.end();
  return {
    props: {},
  };
};
