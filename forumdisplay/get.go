package forumdisplay

import (
	"fmt"
	"regexp"
	"strconv"
)

func getThreadData(fid, cookie string, page int) ([]threadData, error) {
	link := `https://www.mcbbs.net/forum.php?mod=forumdisplay&fid=` + fid + `&filter=author&orderby=dateline&forumdefstyle=yes&page=` + strconv.Itoa(page)
	s, err := httpGet(link, cookie)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	m := make([]threadData, 0, 30)
	list := tidlistreg.FindAll(s, -1)
	for _, v := range list {
		tidtemp := tidreg.FindSubmatch(v)
		if len(tidtemp) != 2 {
			return nil, Regerr{Msg: string(v)}
		}
		tid := string(tidtemp[1])
		titiletemp := titilereg.FindSubmatch(v)
		if len(titiletemp) != 2 {
			return nil, Regerr{Msg: string(v)}
		}
		titile := string(titiletemp[1])
		m = append(m, threadData{tid: tid, titile: titile})
	}
	return m, nil
}

type Regerr struct {
	Msg string
}

func (r Regerr) Error() string {
	return "len != 2: " + r.Msg
}

var (
	tidlistreg = regexp.MustCompile(`<a href=".{10,200}" onclick="atarget\(this\)" class="s xst">.{1,100}</a>`)
	tidreg     = regexp.MustCompile(`tid=([0-9]{1,10})&amp;`)
	titilereg  = regexp.MustCompile(`class="s xst">(.{1,100})</a>`)
)

type threadData struct {
	tid    string
	titile string
	uid    string
	name   string
}
