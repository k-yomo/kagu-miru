import React, { memo } from 'react';

interface Props {
  name: string;
}

export default memo(function CategoryTag({ name }: Props) {
  return (
    <span className="inline-flex items-center px-2.5 py-1.5 rounded bg-indigo-100 dark:bg-indigo-600 text-indigo-800 dark:text-indigo-100 text-xs focus:outline-none">
      {name}
    </span>
  );
});
