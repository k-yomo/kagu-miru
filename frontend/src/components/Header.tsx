import React, { memo } from 'react';
import Link from 'next/link';
import { routes } from '@src/routes/routes';
import SearchBarModal from '@src/components/SearchBarModal';

export default memo(function Header() {
  return (
    <header>
      <div className="relative bg-white dark:bg-black mx-auto px-4 sm:px-6 border-b-[1px] border-gray-200 dark:border-gray-800">
        <div className="flex justify-between items-center h-20 md:justify-start md:space-x-10">
          <div className="flex justify-start lg:w-0 lg:flex-1">
            <Link href={routes.home()}>
              <a className="tracking-wider hover:underline text-black dark:text-white font-bold">
                カグミル
              </a>
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
});
