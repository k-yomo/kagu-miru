// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AppliedMetadata struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type Event struct {
	ID        EventID                `json:"id"`
	Action    Action                 `json:"action"`
	CreatedAt time.Time              `json:"createdAt"`
	Params    map[string]interface{} `json:"params"`
}

type Facet struct {
	Title      string        `json:"title"`
	FacetType  FacetType     `json:"facetType"`
	Values     []*FacetValue `json:"values"`
	TotalCount int           `json:"totalCount"`
}

type FacetValue struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type GetSimilarItemsInput struct {
	ItemID   string `json:"itemId"`
	Page     *int   `json:"page"`
	PageSize *int   `json:"pageSize"`
}

type GetSimilarItemsResponse struct {
	SearchID       string          `json:"searchId"`
	ItemConnection *ItemConnection `json:"itemConnection"`
}

type Item struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Status        ItemStatus          `json:"status"`
	URL           string              `json:"url"`
	AffiliateURL  string              `json:"affiliateUrl"`
	Price         int                 `json:"price"`
	ImageUrls     []string            `json:"imageUrls"`
	AverageRating float64             `json:"averageRating"`
	ReviewCount   int                 `json:"reviewCount"`
	CategoryID    string              `json:"categoryId"`
	Colors        []ItemColor         `json:"colors"`
	Platform      ItemSellingPlatform `json:"platform"`
}

type ItemCategory struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Level    int             `json:"level"`
	ParentID *string         `json:"parentId"`
	Parent   *ItemCategory   `json:"Parent"`
	Children []*ItemCategory `json:"children"`
}

type ItemConnection struct {
	PageInfo *PageInfo `json:"pageInfo"`
	Nodes    []*Item   `json:"nodes"`
}

type PageInfo struct {
	Page       int `json:"page"`
	TotalPage  int `json:"totalPage"`
	TotalCount int `json:"totalCount"`
}

type QuerySuggestionsDisplayActionParams struct {
	Query            string   `json:"query"`
	SuggestedQueries []string `json:"suggestedQueries"`
}

type QuerySuggestionsResponse struct {
	Query            string   `json:"query"`
	SuggestedQueries []string `json:"suggestedQueries"`
}

type SearchClickItemActionParams struct {
	SearchID string `json:"searchId"`
	ItemID   string `json:"itemId"`
}

type SearchDisplayItemsActionParams struct {
	SearchID    string       `json:"searchId"`
	SearchFrom  SearchFrom   `json:"searchFrom"`
	SearchInput *SearchInput `json:"searchInput"`
	ItemIds     []string     `json:"itemIds"`
}

type SearchFilter struct {
	CategoryIds []string              `json:"categoryIds"`
	Platforms   []ItemSellingPlatform `json:"platforms"`
	BrandNames  []string              `json:"brandNames"`
	Colors      []ItemColor           `json:"colors"`
	MinPrice    *int                  `json:"minPrice"`
	MaxPrice    *int                  `json:"maxPrice"`
	MinRating   *int                  `json:"minRating"`
	Metadata    []*AppliedMetadata    `json:"metadata"`
}

type SearchInput struct {
	Query    string         `json:"query"`
	Filter   *SearchFilter  `json:"filter"`
	SortType SearchSortType `json:"sortType"`
	Page     *int           `json:"page"`
	PageSize *int           `json:"pageSize"`
}

type SearchResponse struct {
	SearchID       string          `json:"searchId"`
	ItemConnection *ItemConnection `json:"itemConnection"`
	Facets         []*Facet        `json:"facets"`
}

type SimilarItemsDisplayItemsActionParams struct {
	SearchID             string                `json:"searchId"`
	GetSimilarItemsInput *GetSimilarItemsInput `json:"getSimilarItemsInput"`
	ItemIds              []string              `json:"itemIds"`
}

type Action string

const (
	ActionDisplay   Action = "DISPLAY"
	ActionClickItem Action = "CLICK_ITEM"
)

var AllAction = []Action{
	ActionDisplay,
	ActionClickItem,
}

func (e Action) IsValid() bool {
	switch e {
	case ActionDisplay, ActionClickItem:
		return true
	}
	return false
}

func (e Action) String() string {
	return string(e)
}

func (e *Action) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Action(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Action", str)
	}
	return nil
}

func (e Action) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ErrorCode string

const (
	ErrorCodeNotFound ErrorCode = "NOT_FOUND"
	ErrorCodeInternal ErrorCode = "INTERNAL"
)

var AllErrorCode = []ErrorCode{
	ErrorCodeNotFound,
	ErrorCodeInternal,
}

func (e ErrorCode) IsValid() bool {
	switch e {
	case ErrorCodeNotFound, ErrorCodeInternal:
		return true
	}
	return false
}

func (e ErrorCode) String() string {
	return string(e)
}

func (e *ErrorCode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ErrorCode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ErrorCode", str)
	}
	return nil
}

func (e ErrorCode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type EventID string

const (
	EventIDSearch           EventID = "SEARCH"
	EventIDQuerySuggestions EventID = "QUERY_SUGGESTIONS"
	EventIDSimilarItems     EventID = "SIMILAR_ITEMS"
)

var AllEventID = []EventID{
	EventIDSearch,
	EventIDQuerySuggestions,
	EventIDSimilarItems,
}

func (e EventID) IsValid() bool {
	switch e {
	case EventIDSearch, EventIDQuerySuggestions, EventIDSimilarItems:
		return true
	}
	return false
}

func (e EventID) String() string {
	return string(e)
}

func (e *EventID) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventID(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventID", str)
	}
	return nil
}

func (e EventID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FacetType string

const (
	FacetTypeCategoryIDS FacetType = "CATEGORY_IDS"
	FacetTypeBrandNames  FacetType = "BRAND_NAMES"
	FacetTypeColors      FacetType = "COLORS"
	FacetTypeMetadata    FacetType = "METADATA"
)

var AllFacetType = []FacetType{
	FacetTypeCategoryIDS,
	FacetTypeBrandNames,
	FacetTypeColors,
	FacetTypeMetadata,
}

func (e FacetType) IsValid() bool {
	switch e {
	case FacetTypeCategoryIDS, FacetTypeBrandNames, FacetTypeColors, FacetTypeMetadata:
		return true
	}
	return false
}

func (e FacetType) String() string {
	return string(e)
}

func (e *FacetType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FacetType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FacetType", str)
	}
	return nil
}

func (e FacetType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ItemColor string

const (
	ItemColorWhite       ItemColor = "WHITE"
	ItemColorYellow      ItemColor = "YELLOW"
	ItemColorOrange      ItemColor = "ORANGE"
	ItemColorPink        ItemColor = "PINK"
	ItemColorRed         ItemColor = "RED"
	ItemColorBeige       ItemColor = "BEIGE"
	ItemColorSilver      ItemColor = "SILVER"
	ItemColorGold        ItemColor = "GOLD"
	ItemColorGray        ItemColor = "GRAY"
	ItemColorPurple      ItemColor = "PURPLE"
	ItemColorBrown       ItemColor = "BROWN"
	ItemColorGreen       ItemColor = "GREEN"
	ItemColorBlue        ItemColor = "BLUE"
	ItemColorBlack       ItemColor = "BLACK"
	ItemColorNavy        ItemColor = "NAVY"
	ItemColorKhaki       ItemColor = "KHAKI"
	ItemColorWineRed     ItemColor = "WINE_RED"
	ItemColorTransparent ItemColor = "TRANSPARENT"
)

var AllItemColor = []ItemColor{
	ItemColorWhite,
	ItemColorYellow,
	ItemColorOrange,
	ItemColorPink,
	ItemColorRed,
	ItemColorBeige,
	ItemColorSilver,
	ItemColorGold,
	ItemColorGray,
	ItemColorPurple,
	ItemColorBrown,
	ItemColorGreen,
	ItemColorBlue,
	ItemColorBlack,
	ItemColorNavy,
	ItemColorKhaki,
	ItemColorWineRed,
	ItemColorTransparent,
}

func (e ItemColor) IsValid() bool {
	switch e {
	case ItemColorWhite, ItemColorYellow, ItemColorOrange, ItemColorPink, ItemColorRed, ItemColorBeige, ItemColorSilver, ItemColorGold, ItemColorGray, ItemColorPurple, ItemColorBrown, ItemColorGreen, ItemColorBlue, ItemColorBlack, ItemColorNavy, ItemColorKhaki, ItemColorWineRed, ItemColorTransparent:
		return true
	}
	return false
}

func (e ItemColor) String() string {
	return string(e)
}

func (e *ItemColor) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ItemColor(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ItemColor", str)
	}
	return nil
}

func (e ItemColor) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ItemSellingPlatform string

const (
	ItemSellingPlatformRakuten       ItemSellingPlatform = "RAKUTEN"
	ItemSellingPlatformYahooShopping ItemSellingPlatform = "YAHOO_SHOPPING"
	ItemSellingPlatformPaypayMall    ItemSellingPlatform = "PAYPAY_MALL"
)

var AllItemSellingPlatform = []ItemSellingPlatform{
	ItemSellingPlatformRakuten,
	ItemSellingPlatformYahooShopping,
	ItemSellingPlatformPaypayMall,
}

func (e ItemSellingPlatform) IsValid() bool {
	switch e {
	case ItemSellingPlatformRakuten, ItemSellingPlatformYahooShopping, ItemSellingPlatformPaypayMall:
		return true
	}
	return false
}

func (e ItemSellingPlatform) String() string {
	return string(e)
}

func (e *ItemSellingPlatform) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ItemSellingPlatform(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ItemSellingPlatform", str)
	}
	return nil
}

func (e ItemSellingPlatform) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ItemStatus string

const (
	ItemStatusActive   ItemStatus = "ACTIVE"
	ItemStatusInactive ItemStatus = "INACTIVE"
)

var AllItemStatus = []ItemStatus{
	ItemStatusActive,
	ItemStatusInactive,
}

func (e ItemStatus) IsValid() bool {
	switch e {
	case ItemStatusActive, ItemStatusInactive:
		return true
	}
	return false
}

func (e ItemStatus) String() string {
	return string(e)
}

func (e *ItemStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ItemStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ItemStatus", str)
	}
	return nil
}

func (e ItemStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SearchFrom string

const (
	SearchFromURL             SearchFrom = "URL"
	SearchFromOpenSearch      SearchFrom = "OPEN_SEARCH"
	SearchFromSearch          SearchFrom = "SEARCH"
	SearchFromQuerySuggestion SearchFrom = "QUERY_SUGGESTION"
	SearchFromFilter          SearchFrom = "FILTER"
	SearchFromMedia           SearchFrom = "MEDIA"
)

var AllSearchFrom = []SearchFrom{
	SearchFromURL,
	SearchFromOpenSearch,
	SearchFromSearch,
	SearchFromQuerySuggestion,
	SearchFromFilter,
	SearchFromMedia,
}

func (e SearchFrom) IsValid() bool {
	switch e {
	case SearchFromURL, SearchFromOpenSearch, SearchFromSearch, SearchFromQuerySuggestion, SearchFromFilter, SearchFromMedia:
		return true
	}
	return false
}

func (e SearchFrom) String() string {
	return string(e)
}

func (e *SearchFrom) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SearchFrom(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SearchFrom", str)
	}
	return nil
}

func (e SearchFrom) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SearchSortType string

const (
	SearchSortTypeBestMatch   SearchSortType = "BEST_MATCH"
	SearchSortTypePriceAsc    SearchSortType = "PRICE_ASC"
	SearchSortTypePriceDesc   SearchSortType = "PRICE_DESC"
	SearchSortTypeReviewCount SearchSortType = "REVIEW_COUNT"
	SearchSortTypeRating      SearchSortType = "RATING"
)

var AllSearchSortType = []SearchSortType{
	SearchSortTypeBestMatch,
	SearchSortTypePriceAsc,
	SearchSortTypePriceDesc,
	SearchSortTypeReviewCount,
	SearchSortTypeRating,
}

func (e SearchSortType) IsValid() bool {
	switch e {
	case SearchSortTypeBestMatch, SearchSortTypePriceAsc, SearchSortTypePriceDesc, SearchSortTypeReviewCount, SearchSortTypeRating:
		return true
	}
	return false
}

func (e SearchSortType) String() string {
	return string(e)
}

func (e *SearchSortType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SearchSortType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SearchSortType", str)
	}
	return nil
}

func (e SearchSortType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
