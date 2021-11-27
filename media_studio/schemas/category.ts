import CategoryInput from '../components/CategoryInput';

export default {
  name: 'category',
  title: 'Category',
  type: 'object',
  fields: [
    {
      name: 'id',
      title: 'ID',
      type: 'string',
    },
    {
      name: 'names',
      title: 'Names',
      type: 'array',
      of: [{ type: 'string' }],
    },
  ],
  inputComponent: CategoryInput,
  preview: {
    select: {
      names: 'names',
    },
    prepare({ names }: { names: string[] }) {
      return {
        title: names.join(' > '),
      };
    },
  },
};
