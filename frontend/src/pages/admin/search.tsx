import React from 'react';
import { SearchProvider } from '@src/contexts/search';
import { SearchPageInner } from '@src/pages/search';
import Head from 'next/head';

export default function AdminSearchPage() {
  return (
    <>
      <Head>
        <meta name="robots" content="noindex,nofollow,noarchive" />
      </Head>
      <SearchProvider isAdmin>
        <SearchPageInner isAdmin />
      </SearchProvider>
    </>
  );
}
