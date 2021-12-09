// ==UserScript==
// @name         bump
// @version      0.0.1
// @include      https://www.mcbbs.net/home.php?mod=space&uid=*
// @author       xmdhs
// @description  bump。
// @namespace https://greasyfork.org/users/166541
// ==/UserScript==

(async () => {
    let u = new URL(location.href)
    const uid = u.searchParams.get('uid')
    let d: data = {
        data: [],
        msg: "",
        code: 0
    }
    let f = await fetch(`https://auto.xmdhs.top/getforuid?uid=` + uid)
    d = await f.json()
    if (d.code != 0) {
        console.warn(d.msg)
        alert(d.msg)
        return
    }
    if (d.data == null) {
        console.log("没有数据")
        return
    }

    let p = document.querySelector("#ct > div > div.bm.bw0 > div > div.bm_c.u_profile")
    let div = document.createElement("div")
    p ? p.appendChild(div) : console.warn("没有找到")
    let text = document.createElement("h2")
    text.className = "mbn"
    text.innerText = "提升卡记录"
    div.appendChild(text)
    div.appendChild(makeTable(d.data))
    let psts = document.querySelector("#psts")
    psts && (psts.className = "pbm mbm bbda cl")
    interface rmsg {
        msg: string,
        code: number
    }

    interface data extends rmsg {
        data: {
            uid: number,
            name: string,
            tid: number,
            time: number,
            operation: string,
        }[]
    }

    function countData(data: data["data"]): { [tid: number]: { count: number, lastime: number } } {
        let m: { [tid: number]: { count: number, lastime: number } } = {}
        for (const v of data) {
            if (v.operation.indexOf("提升卡") == -1) {
                continue
            }
            if (m[v.tid] == undefined) {
                m[v.tid] = { count: 0, lastime: v.time }
            } else {
                m[v.tid].count++
                v.time > m[v.tid].lastime && (m[v.tid].lastime = v.time)
            }
        }
        return m
    }

    function makeTable(data: data["data"]): HTMLTableElement {
        const c = countData(data)
        let table = document.createElement("table")
        table.className = "bm dt"
        let tbody = document.createElement("tbody")
        table.appendChild(tbody)
        let tr = document.createElement("tr")
        tbody.appendChild(tr)
        tr.innerHTML = `<th class="xw1">tid</th><th class="xw1">数量</th><th class="xw1">上一次顶贴时间</th>`

        for (const v in c) {
            let trr = document.createElement("tr")
            tbody.appendChild(trr)
            addTr(trr, `<a href="https://www.mcbbs.net/thread-${v}-1-1.html" target="_blank">${v}</a>`, true)
            addTr(trr, String(c[v].count))
            addTr(trr, transformTime(c[v].lastime))
        }
        return table
    }

    function addTr(item: Element, v: string, h?: boolean) {
        let t = document.createElement("td")
        if (h === true) {
            t.innerHTML = v
        } else {
            t.innerText = v
        }
        item.appendChild(t)
    }

    function transformTime(timestamp: number): string {
        var time = new Date(timestamp * 1000);
        var y = time.getFullYear();
        var M = time.getMonth() + 1;
        var d = time.getDate();
        var h = time.getHours();
        var m = time.getMinutes();
        return y + '-' + addZero(M) + '-' + addZero(d) + ' ' + addZero(h) + ':' + addZero(m)
    }
    function addZero(m: number): string {
        return m < 10 ? '0' + String(m) : String(m);
    }
})()