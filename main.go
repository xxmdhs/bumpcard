package main

import (
	"flag"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/xxmdhs/bumpcard/forumdisplay"
	"github.com/xxmdhs/bumpcard/server"
	"github.com/xxmdhs/bumpcard/sql"
)

const fid = 179

func main() {
	s, err := sql.NewSql("data.db")
	if err != nil {
		panic(err)
	}
	if update {
		w := sync.WaitGroup{}
		t := 0
		for i := 1; i <= maxpage; i++ {
			var l []forumdisplay.Thread
			err := retry.Do(func() (err error) {
				l, err = forumdisplay.GetForumList(fid, i)
				return err
			}, retryOpts...)
			if err != nil {
				panic(err)
			}
			for _, v := range l {
				t++
				w.Add(1)
				go func(v forumdisplay.Thread) {
					defer w.Done()
					tid, _ := strconv.Atoi(v.Tid)
					var l []forumdisplay.ActionData
					err := retry.Do(func() (err error) {
						l, err = forumdisplay.GetActionData(tid)
						return err
					}, retryOpts...)
					if err != nil {
						panic(err)
					}
					err = retry.Do(func() (err error) {
						return s.Del(tid)
					}, retryOpts...)
					if err != nil {
						panic(err)
					}
					for _, v := range l {
						d := sql.ActionData{
							Operation: v.Operation,
							Time:      v.Time,
							UID:       v.UID,
							Name:      v.Name,
							TID:       v.TID,
						}

						err := retry.Do(func() (err error) {
							return s.Save(d)
						}, retryOpts...)
						if err != nil {
							panic(err)
						}
					}
				}(v)
				if t >= threads {
					w.Wait()
					time.Sleep(1 * time.Second)
					t = 0
				}
			}
			log.Println(i, "完成")
		}
	} else {
		server.Server(serverport, s)
	}
}

var (
	maxpage    int
	threads    int
	serverport int
	update     bool
)

var retryOpts = []retry.Option{
	retry.Attempts(3),
	retry.Delay(time.Second * 2),
	retry.OnRetry(func(n uint, err error) {
		log.Printf("retry %d: %v", n, err)
	}),
}

func init() {
	flag.IntVar(&maxpage, "maxpage", 10, "max page")
	flag.IntVar(&threads, "threads", 6, "threads")
	flag.IntVar(&serverport, "serverport", 2517, "server port")
	flag.BoolVar(&update, "update", false, "update")
	flag.Parse()
}
