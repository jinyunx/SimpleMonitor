function detailPage(attr_container,graph_container) {
    view = getQueryString("view");
    attr = getQueryString("attr");

    queryInstances(view, function (instances) {
        console.log(instances);
        instances.forEach((item)=>{
            fillAttrDiv(attr_container, item, item);
            queryDetailReportData(graph_container, view, item, attr);
        });
    });
}

function queryInstances(view, ondata) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/view_info?' + 'view=' + view, true);
    httpRequest.send();
    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            console.log(httpRequest.responseText);
            ondata(JSON.parse(httpRequest.responseText)["instances"]);
        }
    }
}

function queryDetailReportData(container, view, instance, attr) {
    var httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/view_detail?' + 'view=' + view + '&attr=' + attr + '&instance=' + instance, true);
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

            containerName = "container" + instance;
            elem = document.createElement("div");
            elem.ondblclick = function(){
                window.open("/?instance="+instance)
            };

            elem.id = containerName;
            elem.style = "max-width:800px;height:400px;border: 1px solid black;margin: 2px";
            container.insertBefore(elem, container.firstChild);
            draw(containerName, instance, attr, lastWeek, yestoday, today);
        }
    };
}
