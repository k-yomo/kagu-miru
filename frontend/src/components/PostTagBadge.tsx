import React, { memo } from 'react';
import Link from 'next/link';
import { routes } from '@src/routes/routes';

interface Props {
  name: string;
}

export default memo(function PostTagBadge({ name }: Props) {
  return (
    <Link href={routes.mediaTag(name)}>
      <a className="inline-flex items-center px-2.5 py-1.5 rounded shadow-md dark:border-2 dark:border-gray-800 text-xs focus:outline-none">
        #{name}
      </a>
    </Link>
  );
});
