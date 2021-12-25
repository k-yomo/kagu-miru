# Objects

### About objects

[Objects](https://graphql.github.io/graphql-spec/June2018/#sec-Objects) in GraphQL represent the resources you can access. An object can contain a list of fields, which are specifically typed.

### Facet

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>facetType</strong> (<a href="enums.md#facettype">FacetType!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>title</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>totalCount</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>values</strong> (<a href="objects.md#facetvalue">[FacetValue!]!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### FacetValue

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>count</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>id</strong> (<a href="scalars.md#id">ID!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>name</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### GetSimilarItemsResponse

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>itemConnection</strong> (<a href="objects.md#itemconnection">ItemConnection!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>searchId</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### Item

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>affiliateUrl</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>averageRating</strong> (<a href="scalars.md#float">Float!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>categoryId</strong> (<a href="scalars.md#id">ID!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>colors</strong> (<a href="enums.md#itemcolor">[ItemColor!]!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>description</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>id</strong> (<a href="scalars.md#id">ID!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>imageUrls</strong> (<a href="scalars.md#string">[String!]!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>name</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>platform</strong> (<a href="enums.md#itemsellingplatform">ItemSellingPlatform!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>price</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>reviewCount</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>status</strong> (<a href="enums.md#itemstatus">ItemStatus!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>url</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### ItemCategory

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>Parent</strong> (<a href="objects.md#itemcategory">ItemCategory</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>children</strong> (<a href="objects.md#itemcategory">[ItemCategory!]!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>id</strong> (<a href="scalars.md#id">ID!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>level</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>name</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>parentId</strong> (<a href="scalars.md#id">ID</a>)</td> 
    <td></td>
  </tr>
</table>

---

### ItemConnection

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>nodes</strong> (<a href="objects.md#item">[Item!]!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>pageInfo</strong> (<a href="objects.md#pageinfo">PageInfo!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### PageInfo

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>page</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>totalCount</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>totalPage</strong> (<a href="scalars.md#int">Int!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### QuerySuggestionsResponse

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>query</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>suggestedQueries</strong> (<a href="scalars.md#string">[String!]!</a>)</td> 
    <td></td>
  </tr>
</table>

---

### SearchResponse

  

#### Fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>facets</strong> (<a href="objects.md#facet">[Facet!]!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>itemConnection</strong> (<a href="objects.md#itemconnection">ItemConnection!</a>)</td> 
    <td></td>
  </tr>
  <tr>
    <td><strong>searchId</strong> (<a href="scalars.md#string">String!</a>)</td> 
    <td></td>
  </tr>
</table>

---