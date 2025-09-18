package score

import "context"

type GetScoreboardInput struct {
	Page     int `params:"page"`
	PageSize int `params:"page_size"`
}

func (gsi *GetScoreboardInput) Validate(ctx context.Context) error {
	if gsi.Page < 1 {
		gsi.Page = 1
	}
	if gsi.PageSize > maxPageSize {
		gsi.PageSize = maxPageSize
	}
	if gsi.PageSize == 0 {
		gsi.PageSize = defaultPageSize
	}
	return nil
}
