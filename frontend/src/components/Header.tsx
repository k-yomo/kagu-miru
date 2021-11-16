import React, { memo, useEffect, useState } from 'react';
import Link from 'next/link';
import { useTheme } from 'next-themes';
import { MoonIcon, SunIcon } from '@heroicons/react/outline';
import { routes } from '@src/routes/routes';

export default memo(function Header() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  return (
    <header>
      <div className="relative bg-white dark:bg-black mx-auto px-4 sm:px-6 border-b-[1px] border-gray-200 dark:border-gray-800">
        <div className="flex justify-between items-center h-20 md:justify-start md:space-x-10">
          <div className="flex justify-start lg:w-0 lg:flex-1">
            <Link href="/">
              <a className="tracking-wider hover:underline text-black dark:text-white font-bold">
                カグミル
              </a>
            </Link>
          </div>

          <div className="flex items-center justify-end flex-1 divide-x-2 divide-black dark:divide-white">
            <div className="pr-3">
              <Link href={routes.media()}>
                <a>
                  <span className="font-bold hover:underline">メディア</span>
                </a>
              </Link>
            </div>
            <div className="pl-3">
              {mounted &&
                (theme === 'light' ? (
                  <SunIcon
                    className="w-8 h-8 cursor-pointer dark:text-white hover:text-gray-500 dark:hover:text-gray-300"
                    onClick={() => setTheme('dark')}
                  />
                ) : (
                  <MoonIcon
                    className="w-8 h-8 cursor-pointer dark:text-white hover:text-gray-500 dark:hover:text-gray-300"
                    onClick={() => setTheme('light')}
                  />
                ))}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
});
