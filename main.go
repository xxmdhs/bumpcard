package main

import (
	"flag"
	"log"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/xmdhs/bumpcard/forumdisplay"
	"github.com/xmdhs/bumpcard/sql"
)

const fid = 179

func main() {
	w := sync.WaitGroup{}
	t := 0
	for i := 1; i <= maxpage; i++ {
		list, err := retry(forumdisplay.GetForumList, []interface{}{int(fid), i}, 5, func(e error) {
			log.Println(e)
		})
		if err != nil {
			log.Fatal(err)
		}
		l := list[0].([]forumdisplay.Thread)
		for _, v := range l {
			t++
			w.Add(1)
			go func(v forumdisplay.Thread) {
				defer w.Done()
				tid, _ := strconv.Atoi(v.Tid)
				list, err := retry(forumdisplay.GetActionData, []interface{}{tid}, 5, func(e error) {
					log.Println(e)
				})
				if err != nil {
					log.Fatal(err)
				}
				l := list[0].([]forumdisplay.ActionData)
				for _, v := range l {
					d := sql.ActionData{
						Operation: v.Operation,
						Time:      v.Time,
						UID:       v.UID,
						Name:      v.Name,
						TID:       v.TID,
					}
					err := sql.Save(d)
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
}

var (
	maxpage int
	threads int
)

func init() {
	flag.IntVar(&maxpage, "maxpage", 10, "max page")
	flag.IntVar(&threads, "threads", 6, "threads")
	flag.Parse()
}

func retry(f interface{}, args []interface{}, retry int, elog func(error)) ([]interface{}, error) {
	r := reflect.ValueOf(f)
	if r.Kind() != reflect.Func {
		panic("need function")
	}
	rargs := make([]reflect.Value, len(args))
	for i, arg := range args {
		rargs[i] = reflect.ValueOf(arg)
	}
	var err error
	for i := 0; i < retry; i++ {
		rv := r.Call(rargs)
		e := rv[len(rv)-1]
		if !e.IsNil() {
			ok := false
			err, ok = e.Interface().(error)
			if ok {
				elog(err)
				continue
			}
		}
		ilist := make([]interface{}, len(rv))
		for i, v := range rv {
			ilist[i] = v.Interface()
		}
		return ilist, nil
	}
	return nil, err
}
