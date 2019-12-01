function start(attr_container,graph_container) {
    view = getQueryString("view");
    instance = getQueryString("instance");
    v = view;
    vName = "view";
    if (view === null || view.length === 0) {
        v = instance;
        vName = "instance";
    }
    queryAttrAndName(v, vName, function (attrArr) {
        console.log(attrArr);
        attrArr.forEach((item)=>{
            fillAttrDiv(attr_container, item["name"]);
            queryReportData(graph_container, v, vName, item["attr"], item["name"]);
        });
    });
}

function getQueryString(name) {
    var reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
    var r = window.location.search.substr(1).match(reg);
    if (r != null) {
        return unescape(r[2]);
    }
    return null;
}

function fillAttrDiv(container, attr) {
    a = document.createElement("a");
    a.href = "javascript:void(0);";
    a.title = attr;
    a.onclick = "";
    a.innerText = attr;

    div = document.createElement("div");
    div.style = "padding: 10px 20px 0 0;float:left";
    div.append(a);

    container.appendChild(div)
}

function queryAttrAndName(v, vName, ondata) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/attr?' + vName + '=' + v, true);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log(httpRequest.responseText);
            ondata(JSON.parse(httpRequest.responseText)["attr_name"]);
        }
    }
}

function queryReportData(container, v, vName, attr, name) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/r?' + vName + '=' + v + '&attr=' + attr, true);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log(httpRequest.responseText);
            d = JSON.parse(httpRequest.responseText);
            console.log(d);

            today = [];
            d["today"]["counters"].forEach((item)=>{
                today.push([item["t"] * 1000, item["counter"]])
            });

            yestoday = [];
            if (d["yestoday"]["counters"] != null) {
                d["yestoday"]["counters"].forEach((item)=>{
                    yestoday.push([item["t"] * 1000 + 24 * 3600 * 1000, item["counter"]])
                });
            }


            lastWeek = [];
            if (d["last_week"]["counters"] != null) {
                d["last_week"]["counters"].forEach((item) => {
                    lastWeek.push([item["t"] * 1000 + 7 * 24 * 3600 * 1000, item["counter"]])
                });
            }

            containerName = "container" + attr;
            elem = document.createElement("div");
            elem.id = containerName;
            elem.style = "max-width:800px;height:400px;border: 1px solid black;margin: 2px";
            container.insertBefore(elem, container.firstChild);
            draw(containerName, v, name, lastWeek,yestoday, today);
        }
    };
}
