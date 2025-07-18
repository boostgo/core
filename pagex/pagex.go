package pagex

const defaultPageSize = 20

type Pagination struct {
	Page int   `default:"1"`
	Size int64 `default:"20"`
}

func MaxPages(size, count int64) int64 {
	if size <= 0 {
		size = defaultPageSize
	}

	if count <= 0 {
		return 1
	}

	maxPages := count / size
	if maxPages == 0 {
		maxPages++
	}

	return maxPages
}
