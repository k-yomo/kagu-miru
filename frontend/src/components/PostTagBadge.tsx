import React, { memo } from 'react';

interface Props {
  name: string;
}

export default memo(function PostTagBadge({ name }: Props) {
  return (
    <span className="inline-flex items-center px-2.5 py-1.5 rounded shadow-md text-xs focus:outline-none">
      #{name}
    </span>
  );
});
