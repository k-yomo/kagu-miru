# Input objects

### About input objects

[Input objects](https://graphql.github.io/graphql-spec/June2018/#sec-Input-Objects) can be described as "composable objects" because they include a set of input fields that define the object.

### Event




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>action</strong> (<a href="enums.md#action">Action!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>createdAt</strong> (<a href="scalars.md#time">Time!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>id</strong> (<a href="enums.md#eventid">EventID!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>params</strong> (<a href="scalars.md#map">Map!</a>)</td>
    <td></td>
  </tr>
</table>

---

### GetSimilarItemsInput




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>itemId</strong> (<a href="scalars.md#id">ID!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>page</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>pageSize</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
</table>

---

### QuerySuggestionsDisplayActionParams




#### Input fields

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

### SearchClickItemActionParams




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>itemId</strong> (<a href="scalars.md#string">String!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>searchId</strong> (<a href="scalars.md#string">String!</a>)</td>
    <td></td>
  </tr>
</table>

---

### SearchDisplayItemsActionParams




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>itemIds</strong> (<a href="scalars.md#id">[ID!]!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>searchFrom</strong> (<a href="enums.md#searchfrom">SearchFrom!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>searchId</strong> (<a href="scalars.md#string">String!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>searchInput</strong> (<a href="input_objects.md#searchinput">SearchInput!</a>)</td>
    <td></td>
  </tr>
</table>

---

### SearchFilter




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>categoryIds</strong> (<a href="scalars.md#id">[ID!]!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>colors</strong> (<a href="enums.md#itemcolor">[ItemColor!]!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>maxPrice</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>minPrice</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>minRating</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>platforms</strong> (<a href="enums.md#itemsellingplatform">[ItemSellingPlatform!]!</a>)</td>
    <td></td>
  </tr>
</table>

---

### SearchInput




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>filter</strong> (<a href="input_objects.md#searchfilter">SearchFilter!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>page</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>pageSize</strong> (<a href="scalars.md#int">Int</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>query</strong> (<a href="scalars.md#string">String!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>sortType</strong> (<a href="enums.md#searchsorttype">SearchSortType!</a>)</td>
    <td></td>
  </tr>
</table>

---

### SimilarItemsDisplayItemsActionParams




#### Input fields

<table>
  <tr>
    <th>Name</th>
    <th>Description</th>
  </tr>
  <tr>
    <td><strong>getSimilarItemsInput</strong> (<a href="input_objects.md#getsimilaritemsinput">GetSimilarItemsInput!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>itemIds</strong> (<a href="scalars.md#id">[ID!]!</a>)</td>
    <td></td>
  </tr>
  <tr>
    <td><strong>searchId</strong> (<a href="scalars.md#string">String!</a>)</td>
    <td></td>
  </tr>
</table>

---