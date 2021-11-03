import T from '@sanity/base/initial-value-template-builder'


const subCategoryTemplate = T.template({
    id: `subCategory`,
    title: `Sub-category`,
    schemaType: `category`,
    parameters: [
        {
            name: `parentCategoryId`,
            title: `Parent Category ID`,
            type: `string`
        }
    ],
    value: parameters => ({
        parent: {
            _type: `reference`,
            _ref: parameters.parentCategoryId
        }
    })
})


export default [...T.defaults(), subCategoryTemplate]
