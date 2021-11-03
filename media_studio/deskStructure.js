import S from '@sanity/desk-tool/structure-builder'
import client from 'part:@sanity/base/client'
import EditIcon from 'part:@sanity/base/edit-icon'



const getCategoryMenuItems = id => {
    const customEditButton = S.menuItem()
        .icon(EditIcon)
        .title(`Edit Category`)
        .showAsAction({ whenCollapsed: true })
        .intent({
            type: `edit`,
            params: { id, type: `category` }
        })

    const defaultItems = S.documentTypeList(`category`).getMenuItems()
    return [...defaultItems, customEditButton]
}



const subCategoryList = async (categoryId) => {
    const category = await client.getDocument(categoryId)

    return S.documentTypeList(`category`)
        .title(category.name)
        .filter('parent._ref == $categoryId')
        .params({ categoryId })
        .menuItems(getCategoryMenuItems(categoryId))
        .canHandleIntent(() => false)
        .initialValueTemplates([
            S.initialValueTemplateItem(
                `subCategory`,
                { parentCategoryId: categoryId }
            )
        ])
        .child(subCategoryList)
}



export default () => S.list()
    .title('Content')
    .items([
        S.listItem()
            .title(`Categories`)
            .child(
                S.documentTypeList(`category`)
                    .title(`Categories`)
                    .filter('_type == "category" && !defined(parent)')
                    .canHandleIntent(() => false)
                    .child(subCategoryList)
            ),

        ...S.documentTypeListItems().filter(item => {
            const id = item.getId()
            return ![`category`].includes(id)
        })
    ])
