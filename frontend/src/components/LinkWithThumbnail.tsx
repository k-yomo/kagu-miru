import React from 'react';
import Link from 'next/link';
import { truncate } from '@src/lib/string';

interface Props {
  url: string;
  title: string;
  subTitle: string;
  imgSrc: string;
}

export default function LinkWithThumbnail({
  url,
  title,
  subTitle,
  imgSrc,
}: Props) {
  return (
    <Link href={url}>
      <a>
        <div className="flex items-center h-[100px] sm:h-[150px] shadow-md rounded-md">
          <div className="w-[30%] h-full overflow-hidden">
            <img
              src={imgSrc}
              alt={title}
              className="w-full h-full object-cover object-center"
            />
          </div>
          <div className="ml-2 sm:ml-3 w-[70%]">
            <div className="mb-2 text-xl font-bold underline">{title}</div>
            <div className="hidden sm:block text-text-secondary dark:text-text-secondary-dark">
              {truncate(subTitle, 100)}
            </div>
          </div>
        </div>
      </a>
    </Link>
  );
}
