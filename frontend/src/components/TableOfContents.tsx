import React, { memo, useEffect } from 'react';
import { ChevronDownIcon } from '@heroicons/react/outline';
// @ts-ignore
import smoothScroll from 'smoothscroll-polyfill';

interface Props {
  headings: string[];
}

export default memo(function TableOfContents({ headings }: Props) {
  useEffect(() => {
    smoothScroll.polyfill();
  }, []);

  const onClickHeading = (i: number) => {
    const section = document.getElementsByTagName('h2');
    section[i]?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };
  const listItems = headings.map((heading, i) => (
    <li key={heading} className="list-none ml-0">
      <a
        href={`#${i}`}
        className="flex items-center mt-2 cursor-pointer"
        onClick={() => onClickHeading(i)}
      >
        <span className="ml-1">{heading}</span>
      </a>
    </li>
  ));
  return (
    <div className="px-2 sm:px-4 py-4 bg-gray-50 dark:bg-gray-900">
      <div className="flex items-center space-x-1 font-bold">
        <ChevronDownIcon className="inline w-3 h-3" />
        <div>目次</div>
      </div>
      <ul className="divide-y-[1px] divide-gray-200 dark:divide-gray-800">
        {listItems}
      </ul>
    </div>
  );
});
