package forumdisplay

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Thread struct {
	Subject  string `json:"subject"`
	Tid      string `json:"tid"`
	Authorid string `json:"authorid"`
}

var c = http.Client{
	Timeout: 10 * time.Second,
}

func GetForumList(fid int, page int) ([]Thread, error) {
	d, err := getforumData(fid, page)
	if err != nil {
		return nil, fmt.Errorf("GetForumList: %w", err)
	}
	t := make([]Thread, 0, len(d.Variables.ForumThreadlist))
	for _, v := range d.Variables.ForumThreadlist {
		t = append(t, Thread{
			Subject:  v.Subject,
			Tid:      v.Tid,
			Authorid: v.Authorid,
		})
	}
	return t, nil
}

func getforumData(fid int, page int) (*thread, error) {
	//version=4&module=forumdisplay&fid=179&page=1&orderby=dateline
	q := url.Values{}
	q.Set("version", "4")
	q.Set("module", "forumdisplay")
	q.Set("fid", strconv.Itoa(fid))
	q.Set("page", strconv.Itoa(page))
	q.Set("orderby", "dateline")
	resp, err := c.Get("https://www.mcbbs.net/api/mobile/index.php?" + q.Encode())
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("getforumData: %w", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getforumData: %w", err)
	}
	var d thread
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, fmt.Errorf("getforumData: %w", err)
	}
	return &d, nil
}

func GetForumPage(fid int) (int, error) {
	t, err := getforumData(fid, 1)
	if err != nil {
		return 0, fmt.Errorf("GetForumPage: %w", err)
	}
	i, err := strconv.Atoi(t.Variables.Forum.Threads)
	if err != nil {
		return 0, fmt.Errorf("GetForumPage: %w", err)
	}
	return i, nil
}
