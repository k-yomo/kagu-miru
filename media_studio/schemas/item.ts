import ItemPreview from '../components/ItemPreview';
import ItemIdInput from '../components/ItemIdInput';

export default {
  name: 'item',
  title: '商品',
  type: 'object',
  fields: [
    {
      name: 'id',
      title: 'ID',
      type: 'string',
      inputComponent: ItemIdInput,
    },
  ],
  preview: {
    select: {
      id: 'id',
    },
    component: ItemPreview,
  },
};
