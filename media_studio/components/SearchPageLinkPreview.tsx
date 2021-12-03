import React, { memo } from 'react';
import { Stack, Card, Label, Text, Flex } from '@sanity/ui';
import { MdSearch } from 'react-icons/md';

interface Props {
  value: {
    title: string;
    query?: string;
    category?: {
      names: string[];
    };
  };
}

export default memo(function ({ value }: Props) {
  const { title, query, category } = value;
  return (
    <Card padding={1}>
      <Flex marginHeight={2} align="center">
        <MdSearch size={20} />
        検索ページリンク
      </Flex>
      <Stack marginY={1} padding={1} space={1}>
        <Text size={1}>タイトル: {title}</Text>
        {query && <Text size={1}>検索クエリ: {query}</Text>}
        {category && <Text size={1}>カテゴリー: {category.names}</Text>}
      </Stack>
    </Card>
  );
});
