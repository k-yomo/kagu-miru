package yahoo_shopping

type Item struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HeadLine    string `json:"headLine"`
	Url         string `json:"url"`
	InStock     bool   `json:"inStock"`
	Code        string `json:"code"`
	Condition   string `json:"condition"`
	ImageId     string `json:"imageId"`
	Image       struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
	} `json:"image"`
	Review struct {
		Rate  float64 `json:"rate"`
		Count int     `json:"count"`
		Url   string  `json:"url"`
	} `json:"review"`
	AffiliateRate       float64 `json:"affiliateRate"`
	Price               int     `json:"price"`
	PremiumPrice        int     `json:"premiumPrice"`
	PremiumPriceStatus  bool    `json:"premiumPriceStatus"`
	PremiumDiscountType *string `json:"premiumDiscountType,omitempty"`
	PremiumDiscountRate *int    `json:"premiumDiscountRate,omitempty"`
	PriceLabel          struct {
		Taxable         bool `json:"taxable"`
		DefaultPrice    int  `json:"defaultPrice"`
		DiscountedPrice *int `json:"discountedPrice,omitempty"`
		FixedPrice      *int `json:"fixedPrice,omitempty"`
		PremiumPrice    *int `json:"premiumPrice,omitempty"`
		PeriodStart     *int `json:"periodStart,omitempty"`
		PeriodEnd       *int `json:"periodEnd,omitempty"`
	} `json:"priceLabel"`
	Point struct {
		Amount        int `json:"amount"`
		Times         int `json:"times"`
		PremiumAmount int `json:"premiumAmount"`
		PremiumTimes  int `json:"premiumTimes"`
	} `json:"point"`
	Shipping struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"shipping"`
	GenreCategory struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Depth int    `json:"depth"`
	} `json:"genreCategory"`
	ParentGenreCategories []struct {
		Depth int    `json:"depth"`
		Id    int    `json:"id"`
		Name  string `json:"name"`
	} `json:"parentGenreCategories"`
	Brand struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"brand"`
	ParentBrands []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"parentBrands"`
	JanCode     string `json:"janCode"`
	ReleaseDate *int   `json:"releaseDate,omitempty"`
	Seller      struct {
		SellerId      string `json:"sellerId"`
		Name          string `json:"name"`
		Url           string `json:"url"`
		IsBestSeller  bool   `json:"isBestSeller"`
		IsPMallSeller bool   `json:"isPMallSeller"`
		Review        struct {
			Rate  float64 `json:"rate"`
			Count int     `json:"count"`
		} `json:"review"`
		ImageId string `json:"imageId"`
	} `json:"seller"`
	Delivery struct {
		Area     string `json:"area"`
		DeadLine *int   `json:"deadLine,omitempty"`
		Day      *int   `json:"day,omitempty"`
	} `json:"delivery"`
	Payment string `json:"payment"`
}
