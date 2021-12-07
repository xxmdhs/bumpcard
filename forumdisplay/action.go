package forumdisplay

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"time"
	_ "time/tzdata"
)

type ActionData struct {
	Operation string
	Time      int64
	UID       int
	Name      string
	TID       int
}

var actionReg = regexp.MustCompile(`.*<td><a href\="home.php\?mod=space&amp;uid\=(\d{1,15})" target\="_blank">(.{1,10}?)</a></td>.*\n.*<td><span title="(.{1,30}?)">.*</span></td>.*\n.*<td >(.{1,30}?)</td>`)

func ParseActionData(xmldata string, tid int) []ActionData {
	ret := []ActionData{}
	for _, v := range actionReg.FindAllStringSubmatch(xmldata, -1) {
		ret = append(ret, ActionData{
			Operation: v[4],
			Time: func() int64 {
				//2021-12-6 22:42
				t, err := time.Parse("2006-01-2 15:04", v[3])
				if err != nil {
					return 0
				}
				//Asia/Shanghai
				l, err := time.LoadLocation("Asia/Shanghai")
				if err != nil {
					return 0
				}
				return t.In(l).Unix()
			}(),
			UID: func() int {
				i, _ := strconv.Atoi(v[1])
				return i
			}(),
			Name: v[2],
			TID:  tid,
		})
	}
	return ret
}

func GetActionData(tid int) ([]ActionData, error) {
	//https://www.mcbbs.net/forum.php?mod=misc&action=viewthreadmod&tid=1276429&infloat=yes&handlekey=viewthreadmod&inajax=1&ajaxtarget=fwin_content_viewthreadmod
	v := url.Values{}
	v.Set("mod", "misc")
	v.Set("action", "viewthreadmod")
	v.Set("tid", strconv.Itoa(tid))
	v.Set("infloat", "yes")
	v.Set("handlekey", "viewthreadmod")
	v.Set("inajax", "1")
	v.Set("ajaxtarget", "fwin_content_viewthreadmod")
	resp, err := c.Get("https://www.mcbbs.net/forum.php?" + v.Encode())
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("GetActionData: %w", err)
	}
	xmldata, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetActionData: %w", err)
	}
	return ParseActionData(string(xmldata), tid), nil
}
