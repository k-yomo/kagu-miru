import React, { memo } from 'react';
import Link from 'next/link';
import { routes } from '@src/routes/routes';

interface Props {
  id: string;
  name: string;
  enableLink?: boolean;
}

export default memo(function PostCategoryBadge({
  id,
  name,
  enableLink,
}: Props) {
  const Badge = () => {
    return (
      <span className="inline-flex items-center px-2.5 py-1.5 rounded bg-gradient-to-r from-pink-500 dark:from-pink-600 to-rose-500 dark:to-rose-600 text-white text-xs focus:outline-none">
        {name}
      </span>
    );
  };
  if (enableLink) {
    return (
      <Link href={routes.mediaCategory(id)}>
        <a>
          <Badge />
        </a>
      </Link>
    );
  }
  return <Badge />;
});
