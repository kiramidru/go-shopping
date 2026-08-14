package catalog

type VolumeRequest struct {
	Query      *string `json:"q" form:"q" binding:"omitempty"`
	MaxResults int     `json:"maxResults" form:"maxResults" binding:"omitempty,min=1,max=40"`
	StartIndex int     `json:"startIndex" form:"startIndex" binding:"omitempty,min=0"`
}

type VolumeResponse struct {
	ID         string     `json:"id"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
	SaleInfo   SaleInfo   `json:"saleInfo"`
}

type VolumeInfo struct {
	Title         string     `json:"title"`
	Subtitle      string     `json:"subtitle"`
	Authors       []string   `json:"authors"`
	Publisher     string     `json:"publisher"`
	PublishedDate string     `json:"publishedDate"`
	PageCount     int        `json:"pageCount"`
	Categories    []string   `json:"categories"`
	Language      string     `json:"language"`
	ImageLinks    ImageLinks `json:"imageLinks"`
}

type ImageLinks struct {
	Thumbnail string `json:"thumbnail"`
}

type SaleInfo struct {
	ListPrice ListPrice `json:"listPrice"`
}

type ListPrice struct {
	Amount float64 `json:"amount"`
}

type CatalogResponse struct {
	TotalItems int               `json:"totalItems"`
	Items      []*VolumeResponse `json:"items"`
}
