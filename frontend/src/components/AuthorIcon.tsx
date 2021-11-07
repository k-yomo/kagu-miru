import React, { memo } from 'react';

interface Props {
  name: string;
  imgSrc: string;
}

export default memo(function AuthorIcon({ name, imgSrc }: Props) {
  return (
    <div className="flex items-center">
      <img
        src={imgSrc}
        alt={`${name}のプロフィール画像`}
        className="w-8 h-8 object-cover object-center rounded-full"
      />
      <span className="ml-1">{name}</span>
    </div>
  );
});
