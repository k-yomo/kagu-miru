import React, { memo } from 'react';

interface Props {
  name: string;
}

export default memo(function PostTag({ name }: Props) {
  return (
    <span className="inline-flex items-center px-2.5 py-1.5 rounded bg-gradient-to-r from-primary-500 dark:from-primary-600 to-rose-500 dark:to-rose-600 text-white text-xs focus:outline-none">
      {name}
    </span>
  );
});
