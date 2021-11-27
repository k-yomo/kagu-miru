import CustomHtmlPreview from '../components/CustomHtmlPreview';

export default {
  name: 'customHtml',
  title: 'Custom HTML',
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
