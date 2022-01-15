import { Rule } from '@sanity/types';

export default {
  title: '検索ページリンク',
  name: 'searchPageLink',
  type: 'object',
  fields: [
    {
      title: 'タイトル',
      name: 'title',
      type: 'string',
      validation: (Rule: Rule) => Rule.required(),
    },
    {
      title: '検索ページURL',
      name: 'url',
      type: 'url',
    },
  ],
  preview: {
    select: {
      title: 'title',
    },
  },
};
