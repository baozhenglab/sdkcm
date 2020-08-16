package sdkcm

import "strings"

type OrderBy struct {
	Key    string
	IsDesc bool
}

type Paging struct {
	Cursor     *UID      `json:"cursor" form:"-"`
	NextCursor string    `json:"next_cursor" form:"-"`
	CursorStr  string    `json:"-" form:"cursor"`
	Limit      int       `json:"limit" form:"limit"`
	Total      int       `json:"total" form:"-"`
	Page       int       `json:"page" form:"page"`
	HasNext    bool      `json:"has_next" form:"-"`
	OrderBy    string    `json:"order_by" form:"order_by"`
	OB         []OrderBy `json:"-" form:"-"`
}

func (p *Paging) FullFill() {
	if p.Cursor != nil && p.Cursor.localID == 0 {
		p.Cursor = nil
	}

	if p.CursorStr != "" {
		uid, err := FromBase58(p.CursorStr)
		if err == nil {
			p.Cursor = &uid
		}
	}

	if p.Limit <= 0 {
		p.Limit = 50
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	if strings.TrimSpace(p.OrderBy) == "" {
		p.OrderBy = "id desc"
		p.OB = []OrderBy{{Key: "id", IsDesc: true}}
	} else {
		p.OB = getOrderBy(p.OrderBy)
	}
}

func getOrderBy(ord string) []OrderBy {
	comps := strings.Split(ord, ",")
	result := make([]OrderBy, len(comps))

	for i := range comps {
		kvs := strings.Split(strings.TrimSpace(comps[i]), " ")
		result[i] = OrderBy{Key: kvs[0], IsDesc: len(kvs) == 1 || kvs[1] == "-1"}
	}

	return result
}
