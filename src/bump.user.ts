// ==UserScript==
// @name         提升卡记录
// @version      0.1.2
// @match        https://www.mcbbs.net/home.php?mod=space*
// @match        https://www.mcbbs.net/?*
// @match        https://www.mcbbs.net/forum.php?mod=viewthread&tid=*
// @match        https://www.mcbbs.net/thread-*.html
// @author       xmdhs
// @license      MIT
// @description  查看用户的提升卡使用记录
// @namespace    https://greasyfork.org/users/166541
// ==/UserScript==

(async () => {
    if (location.href.startsWith("https://www.mcbbs.net/home.php?mod=space") || location.href.startsWith("https://www.mcbbs.net/?")) {
        await userPage()
        return
    } else {
        let doms: NodeListOf<Element>;
        try {
            doms = document.querySelectorAll("div.pi > div > a.xw1")
        } catch (e) {
            return
        }
        let i = 0
        for (const dom of Array.from(doms)) {
            i++
            if (!(dom instanceof HTMLAnchorElement)) {
                continue
            }
            let u = new URL(dom.href);
            const uid = u.searchParams.get("uid");
            if (!uid) continue;
            dosome(uid, dom, 0)
            if (i > 5) {
                await new Promise((r) => setTimeout(r, 1000))
                i = 0
            }
        }
    }


    async function userPage() {
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
        const profile = document.querySelector(".bm_c.u_profile")
        if (profile && profile.lastElementChild) {
            profile.lastElementChild.className = "pbm mbm bbda cl"
        }
        const p = document.querySelector("#ct > div > div.bm.bw0 > div > div.bm_c.u_profile")
        const div = document.createElement("div")
        p ? p.appendChild(div) : console.warn("没有找到")
        const text = document.createElement("h2")
        text.className = "mbn"
        text.innerText = "提升卡记录"
        div.appendChild(text)
        if (d.data == null) {
            console.log("没有数据")
            return
        }
        div.appendChild(makeTable(d.data))
    }

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

    async function dosome(uid: string, dom: HTMLElement, i: number) {
        i++
        if (i > 3) {
            console.warn(`${uid} 失败超过3次，跳过`)
            return
        }
        try {
            let d: data = {
                data: [],
                msg: "",
                code: 0
            }
            const data = await getData(uid);
            if (data.length == 0) {
                return
            }
            const dd = document.createElement("dd")
            let c = 0
            data.forEach(v => {
                if (v.operation.indexOf("提升卡") != -1) {
                    c++
                }
            })
            dd.textContent = `${c} 张`
            const dt = document.createElement("dt")
            const timg = document.createElement("img")
            timg.src = "https://www.mcbbs.net/source/plugin/mcbbs_mcserver_plus/magic/magic_serverBump.small.gif"
            timg.style.verticalAlign = "middle"
            dt.textContent = ` 提升`
            dt.style.color = "red"
            dt.style.fontWeight = "bold"
            dt.insertBefore(timg, dt.firstChild)
            const dl = dom?.parentNode?.parentNode?.parentNode?.querySelector("dl.pil")
            dl?.appendChild(dt)
            dl?.appendChild(dd)
        } catch (e) {
            console.warn(e)
            await new Promise((r) => setTimeout(r, 2000))
            await dosome(uid, dom, i)
        }
    }

    async function getData(uid: string): Promise<data["data"]> {
        let f = await fetch(`https://auto.xmdhs.com/getforuid?uid=` + uid);
        let d: data = {
            data: [],
            msg: "",
            code: 0
        }
        d = await f.json();
        if (d.code != 0) {
            throw new Error(d.msg);
        }
        return d.data
    }
})()