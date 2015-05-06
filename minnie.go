package minnie

import (
	"net/url"
	"testing"

	"github.com/sclevine/agouti"
)

type Finder interface {
	Find(string) *agouti.Selection
	All(string) *agouti.MultiSelection
}

type Page struct {
	Page Finder
	T    *testing.T
}

func New(t *testing.T, p *agouti.Page) *Page {
	return &Page{
		Page: p,
		T:    t,
	}
}

func Find(str string, f Finder) (*agouti.Selection, error) {
	return f.Find(str), nil
}

func All(str string, f Finder) (*agouti.MultiSelection, error) {
	return f.All(str), nil
}

func (p *Page) check(err error) {
	if err != nil {
		p.Error(err)
	}
}

func (p *Page) Visit(url string) {
	page, ok := p.Page.(*agouti.Page)
	if !ok {
		p.Fatal("cannot call Visit(string) on %T: invalid type", p.Page)
	}

	p.check(page.Navigate(url))
}

func (p *Page) find(s string, f Finder) *agouti.Selection {
	if f == nil {
		f = p.Page
	}

	sel, err := Find(s, f)
	if err != nil {
		p.Fatal(err)
	}

	return sel
}

func (p *Page) Find(s string) *agouti.Selection {
	return p.find(s, p.Page)
}

// TODO FindContain both string and selector

func (p *Page) all(s string, sel Finder) (*agouti.MultiSelection, error) {
	if sel == nil {
		sel = p.Page
	}

	return All(s, sel)
}

func (p *Page) All(s string) (*agouti.MultiSelection, error) {
	return p.all(s, nil)
}

func (p *Page) Click(s string) {
	p.check(p.find(s, nil).Click())
}

func (p *Page) Hover(s string) {
	p.check(p.Find(s).MouseToElement())
}

func (p *Page) Fill(s string, v string) {
	p.check(p.find(s, nil).Fill(v))
}

func (p *Page) Confirm() {
	page, ok := p.Page.(*agouti.Page)
	if !ok {
		p.Fatalf("cannot call call ConfirmPopup() on %T: invalid type", p.Page)
	}

	p.check(page.ConfirmPopup())
}

// TODO unconfirm or cancel, not sure what the write name for this is quite yet

// Count counts selector and returns the MultiSelection and count
func (p *Page) Count(s string) (*agouti.MultiSelection, int) {
	sels, _ := p.All(s)

	n, err := sels.Count()
	if err != nil {
		p.Fatal(err)
	}

	return sels, n
}

// URL returns the page url parsed through net/url
func (p *Page) URL() *url.URL {
	page, ok := p.Page.(*agouti.Page)
	if !ok {
		p.Fatal("cannot call URL() on %T: invalid type", p.Page)
	}

	urlStr, err := page.URL()
	if err != nil {
		p.Fatal(err)
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		p.Fatal(err)
	}

	return u
}

func (p *Page) Text(s string) string {
	str, err := p.Find(s).Text()
	if err != nil {
		p.Fatal(err)
	}

	return str
}

func (p *Page) Visible(s string) bool {
	ok, err := p.Find(s).Visible()
	if err != nil {
		p.Fatal(err)
	}

	return ok
}

func (p *Page) Within(s string, fn func(*Page)) {
	in := p
	if s != "" {
		in = &Page{
			Page: p.Find(s),
			T:    p.T,
		}
	}

	fn(in)
}

type AssertFunc func(*Page)

func (p *Page) Assert(fn AssertFunc) {
	fn(p)
}

/*
Testing shortcuts

*/

func (p *Page) Fatal(v ...interface{}) {
	p.T.Fatal(v...)
}

func (p *Page) Fatalf(s string, v ...interface{}) {
	p.T.Fatalf(s, v...)
}

func (p *Page) Error(v ...interface{}) {
	p.T.Error(v...)
}

func (p *Page) Errorf(s string, v ...interface{}) {
	p.T.Errorf(s, v...)
}
