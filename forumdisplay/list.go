package forumdisplay

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type Thread struct {
	Subject string `json:"subject"`
	Tid     string `json:"tid"`
}

var c = http.Client{
	Timeout: 10 * time.Second,
}

func GetForumList(fid int, page int, cookie string) ([]Thread, error) {
	d, err := getThreadData(strconv.Itoa(fid), cookie, page)
	if err != nil {
		return nil, fmt.Errorf("GetForumList: %w", err)
	}
	t := make([]Thread, 0, len(d))
	for _, v := range d {
		t = append(t, Thread{
			Subject: v.titile,
			Tid:     v.tid,
		})
	}
	return t, nil
}

var pageReg = regexp.MustCompile(`<span title="共 (\d{1,5}?) 页"> / \d{1,7} 页</span>`)

var ErrNotFound = fmt.Errorf("NotFound")

func GetForumPage(fid int, cookie string) (int, error) {
	b, err := httpGet(`https://www.mcbbs.net/forum.php?mod=forumdisplay&fid=`+strconv.Itoa(fid), cookie)
	if err != nil {
		return 0, fmt.Errorf("GetForumPage: %w", err)
	}
	l := pageReg.FindSubmatch(b)
	if len(l) != 2 {
		return 1, nil
	}

	i, err := strconv.Atoi(string(l[1]))
	if err != nil {
		return 0, fmt.Errorf("GetForumPage: %w", err)
	}
	return i, nil
}
