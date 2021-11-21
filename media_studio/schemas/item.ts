import ItemPreview from "../components/ItemPreview"

export default {
  name: 'item',
  title: 'Item',
  type: 'object',
  fields: [
    {
      name: 'id',
      title: 'ID',
      type: 'string',
    },
  ],
  preview: {
    select: {
      id: 'id'
    },
    component: ItemPreview
  }
}

