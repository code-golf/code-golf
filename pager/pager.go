package pager

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

const PerPage = 50

type Pager struct {
	lastPage, First, Last, Offset, Page, Total int
	base, Prev, Next                           *url.URL
}

func New(r *http.Request) *Pager {
	page, _ := strconv.Atoi(r.FormValue("page"))
	page = max(page, 1)

	return &Pager{base: r.URL, Offset: (page - 1) * PerPage, Page: page}
}

func (p *Pager) Calculate() bool {
	// Falling off the end is a 404, but no results isn't.
	if p.Total == 0 && p.Page > 1 {
		return true
	}

	p.lastPage = int(math.Ceil(float64(p.Total) / float64(PerPage)))
	if p.lastPage == 0 {
		p.lastPage = 1
	}

	if p.Total > 0 {
		p.First = (p.Page-1)*PerPage + 1
	}

	if p.Page == p.lastPage {
		p.Last = p.Total
	} else {
		p.Last = p.Page * PerPage
	}

	if p.Page > 1 {
		p.Prev = changePage(p.base, p.Page-1)
	}

	if p.Page < p.lastPage {
		p.Next = changePage(p.base, p.Page+1)
	}

	return false
}

func changePage(u *url.URL, page int) *url.URL {
	q := u.Query()
	if page == 1 {
		q.Del("page")
	} else {
		q.Set("page", strconv.Itoa(page))
	}

	new := *u
	new.RawQuery = q.Encode()

	return &new
}
