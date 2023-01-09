package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/xxmdhs/bumpcard/forumdisplay"
	"github.com/xxmdhs/bumpcard/server"
	"github.com/xxmdhs/bumpcard/sql"
)

var fid = []int{179, 296, 647, 436}

func main() {
	s, err := sql.NewSql("data.db")
	if err != nil {
		panic(err)
	}
	defer s.Close()
	if update {
		for _, v := range fid {
			get(s, v)
		}
	} else {
		server.Server(serverport, s)
	}
}

func get(s *sql.DB, fid int) {
	w := sync.WaitGroup{}
	t := 0
	var maxpage int
	err := retry.Do(func() (err error) {
		maxpage, err = forumdisplay.GetForumPage(fid, cookie)
		return err
	}, retryOpts...)
	if err != nil {
		panic(err)
	}
	log.Printf("fid %v,总页数: %d", fid, maxpage)
	for i := 1; i <= maxpage; i++ {
		var l []forumdisplay.Thread
		err := retry.Do(func() (err error) {
			l, err = forumdisplay.GetForumList(fid, i, cookie)
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
					l, err = forumdisplay.GetActionData(tid, cookie)
					return err
				}, retryOpts...)
				if err != nil {
					panic(err)
				}
				err = retry.Do(func() (err error) {
					tx, err := s.NewTx()
					defer tx.Rollback()
					if err != nil {
						return err
					}
					err = s.Del(tx, tid)
					if err != nil {
						return err
					}
					for _, v := range l {
						d := sql.ActionData{
							Operation: v.Operation,
							Time:      v.Time,
							UID:       v.UID,
							Name:      v.Name,
							TID:       v.TID,
						}
						err = s.Save(tx, d)
						if err != nil {
							return err
						}
					}
					return tx.Commit()
				}, retryOpts...)
				if err != nil {
					panic(err)
				}

			}(v)
			if t >= threads {
				w.Wait()
				time.Sleep(1 * time.Second)
				t = 0
			}
		}
		log.Printf("fid: %v, page: %v", fid, i)
	}
}

var (
	threads    int
	serverport int
	update     bool
	cookie     string
)

var retryOpts = []retry.Option{
	retry.Attempts(7),
	retry.Delay(time.Second * 2),
	retry.MaxDelay(5 * time.Second),
	retry.LastErrorOnly(true),
	retry.OnRetry(func(n uint, err error) {
		log.Printf("retry %d: %v", n, err)
	}),
}

func init() {
	flag.IntVar(&threads, "threads", 6, "threads")
	flag.IntVar(&serverport, "serverport", 2517, "server port")
	flag.BoolVar(&update, "update", false, "update")
	flag.Parse()
	var err error
	var c []byte
	c, err = os.ReadFile("cookie.txt")
	if err != nil {
		panic(err)
	}
	cookie = string(c)
}
