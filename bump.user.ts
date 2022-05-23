// ==UserScript==
// @name         提升卡记录
// @version      0.0.5
// @include      https://www.mcbbs.net/home.php?mod=space*
// @include      https://www.mcbbs.net/?*
// @author       xmdhs
// @license MIT
// @description  查看用户的提升卡使用记录
// @namespace https://greasyfork.org/users/166541
// ==/UserScript==

(async () => {
    const uid = getuid();
    if (uid === null) {
        alert("无法获取 uid")
        return;
    }
    let d: data = {
        data: [],
        msg: "",
        code: 0
    }
    let f = await fetch(`https://auto.xmdhs.com/getforuid?uid=` + uid)
    d = await f.json()
    if (d.code != 0) {
        console.warn(d.msg)
        alert(d.msg)
        return
    }
    let p = document.querySelector("#ct > div > div.bm.bw0 > div > div.bm_c.u_profile")
    let div = document.createElement("div")
    p ? p.appendChild(div) : console.warn("没有找到")
    let text = document.createElement("h2")
    text.className = "mbn"
    text.innerText = "提升卡记录"
    div.appendChild(text)
    if (d.data == null) {
        console.log("没有数据")
        return
    }
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
                m[v.tid] = { count: 1, lastime: v.time }
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

    function getuid(): string | null {
        let u = new URL(location.href)
        let uid = u.searchParams.get('uid')
        if (uid && uid.length > 0) {
            return uid
        }
        let dom = document.querySelector("#uhd > div > div > a")
        if ((dom as HTMLAnchorElement).href.length > 0) {
            u = new URL((dom as HTMLAnchorElement).href)
            return u.searchParams.get('uid')
        }
        return null
    }
})()