package forumdisplay

import (
	"fmt"
	"io"
	"net/http"
)

func httpGet(url, cookie string) ([]byte, error) {
	reqs, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("httpGet: %w", err)
	}
	reqs.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	reqs.Header.Set("cookie", cookie)
	rep, err := c.Do(reqs)
	if rep != nil {
		defer rep.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("httpGet: %w", err)
	}
	b, err := io.ReadAll(rep.Body)
	if err != nil {
		return nil, fmt.Errorf("httpGet: %w", err)
	}
	return b, nil
}
