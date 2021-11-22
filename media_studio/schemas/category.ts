import { Rule } from "@sanity/types"

export default {
  name: 'category',
  title: 'Category',
  type: 'document',
  fields: [
    {
      name: 'id',
      title: 'ID',
      type: 'string',
      validation: (Rule: Rule) => Rule.required()
    },
    {
      name: `parent`,
      title: `Parent Category`,
      type: `reference`,
      to: [{ type: `category` }]
    },
    {
      name: 'name',
      title: 'Name',
      type: 'string',
      validation: (Rule: Rule) => Rule.required()
    },
    {
      name: 'description',
      title: 'Description',
      type: 'text',
    },
  ],
}
