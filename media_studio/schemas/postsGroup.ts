import { Rule } from '@sanity/types';

export default {
  name: 'postsGroup',
  title: '記事グループ',
  type: 'document',
  fields: [
    {
      name: 'id',
      title: 'id',
      type: 'string',
      validation: (Rule: Rule) => Rule.required(),
    },
    {
      name: 'title',
      title: 'タイトル',
      type: 'string',
      validation: (Rule: Rule) => Rule.required(),
    },
    {
      title: '対象記事',
      name: 'posts',
      type: 'array',
      of: [{ type: 'reference', to: { type: 'post' } }],
      validation: (Rule: Rule) => Rule.required(),
    },
  ],
};
