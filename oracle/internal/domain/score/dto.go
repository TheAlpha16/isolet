package score

import "context"

type GetScoreboardInput struct {
	Page     int
	PageSize int
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
	if gsi.PageSize < minPageSize {
		gsi.PageSize = minPageSize
	}
	return nil
}
