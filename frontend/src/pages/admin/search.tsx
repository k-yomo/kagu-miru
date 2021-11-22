import React from 'react';
import { SearchProvider } from '@src/contexts/search';
import { TopPageInner } from '@src/pages';
import Head from 'next/head';

export default function AdminSearchPage() {
  return (
    <>
      <Head>
        <meta name="robots" content="noindex,nofollow,noarchive" />
      </Head>
      <SearchProvider isAdmin>
        <TopPageInner isAdmin />
      </SearchProvider>
    </>
  );
}
