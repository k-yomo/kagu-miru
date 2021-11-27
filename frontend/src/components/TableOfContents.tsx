import React, { memo } from 'react';
import Link from 'next/link';
import { ChevronRightIcon } from '@heroicons/react/solid';

interface Props {
  headings: string[];
}

export default memo(function TableOfContents({ headings }: Props) {
  const onClickHeading = (i: number) => {
    const section = document.getElementsByTagName('h2');
    section[i]?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };
  const listItems = headings.map((heading, i) => (
    <li key={heading} className="list-none ml-0">
      <Link href={`#${i}`}>
        <a className="flex items-center" onClick={() => onClickHeading(i)}>
          <ChevronRightIcon className="h-5 w-5 min-w-5" />
          <span className="ml-1">{heading}</span>
        </a>
      </Link>
    </li>
  ));
  return (
    <div className="px-2 sm:px-4 py-4 bg-gray-50 dark:bg-gray-900">
      <div className="font-bold">目次</div>
      <ul className="space-y-3 font-bold">{listItems}</ul>
    </div>
  );
});
