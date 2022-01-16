import React from 'react';
import { NextPage } from 'next';
import router from 'next/router';
import { ParsedUrlQuery } from 'querystring';
import { routes } from '@src/routes/routes';

const TopPage: NextPage = () => <>abc</>;

TopPage.getInitialProps = async ({ req, res, query }) => {
  if (typeof window === 'undefined' && res) {
    if (query) {
      res.writeHead(301, {
        Location: `${routes.search()}?${toQueryString(query)}`,
      });
    } else {
      res.writeHead(301, { Location: `${routes.search()}` });
    }
    res.end();

    return {};
  }

  router.push(`${routes.search()}${router.asPath}`);

  return {};
};

const toQueryString = (query: ParsedUrlQuery) => {
  return Object.keys(query)
    .filter((key) => query[key] !== null && query[key] !== undefined)
    .map((key) => {
      let value = query[key];

      if (Array.isArray(value)) {
        value = value.join('/');
      }

      return [
        encodeURIComponent(key),
        encodeURIComponent(value as string),
      ].join('=');
    })
    .join('&');
};

export default TopPage;
