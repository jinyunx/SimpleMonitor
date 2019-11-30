function start(attr_container,graph_container) {
    instance = getQueryString("instance");
    queryAttrAndName(instance, function (attrArr) {
        console.log(attrArr);
        attrArr.forEach((item)=>{
            fillAttrDiv(attr_container, item["name"]);
            queryReportData(graph_container, instance, item["attr"], item["name"]);
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

function queryAttrAndName(instance, ondata) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/attr?instance=' + instance, true);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log(httpRequest.responseText);
            ondata(JSON.parse(httpRequest.responseText)["attr_name"]);
        }
    }
}

function queryReportData(container, instance, attr, name) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/r?instance=' + instance+ '&attr=' + attr, true);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log(httpRequest.responseText);
            d = JSON.parse(httpRequest.responseText);
            data = [];
            d["counters"].forEach((item)=>{
                data.push([item["t"] * 1000, item["counter"]])
            });
            console.log(d);

            containerName = "container" + attr;
            elem = document.createElement("div");
            elem.id = containerName;
            elem.style = "max-width:800px;height:400px;border: 1px solid black;margin: 2px";
            container.insertBefore(elem, container.firstChild);
            draw(containerName,instance, name, [],[], data);
        }
    };
}
