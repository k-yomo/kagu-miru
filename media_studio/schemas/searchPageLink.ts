import { Rule } from '@sanity/types';
import SearchPageLinkPreview from '../components/SearchPageLinkPreview';

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
      title: 'URL',
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
