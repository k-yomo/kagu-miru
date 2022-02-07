import { Rule } from '@sanity/types';
import BoxPreview from '../components/BoxPreview';

export default {
  title: 'ボックス',
  name: 'box',
  type: 'object',
  fields: [
    {
      title: 'タイトル',
      name: 'title',
      type: 'string',
    },
    {
      title: '内容',
      name: 'content',
      type: 'text',
      validation: (Rule: Rule) => Rule.required(),
    },
  ],
  preview: {
    select: {
      title: 'title',
      content: 'content',
    },
    component: BoxPreview,
  },
};
