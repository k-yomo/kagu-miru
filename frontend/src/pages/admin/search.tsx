import React from 'react';
import { SearchProvider } from '@src/contexts/search';
import { TopPageInner } from '@src/pages';

export default function AdminSearchPage() {
  return (
    <SearchProvider isAdmin>
      <TopPageInner isAdmin />
    </SearchProvider>
  );
}
