import CustomHtmlPreview from '../components/CustomHtmlPreview';

export default {
  name: 'customHtml',
  title: 'カスタムHTML',
  type: 'object',
  fields: [
    {
      name: 'html',
      title: 'HTML',
      type: 'text',
    },
  ],
  preview: {
    select: {
      html: 'html',
    },
    component: CustomHtmlPreview,
  },
};
