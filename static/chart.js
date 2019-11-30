Highcharts.setOptions({ global: { useUTC: false } });
function draw(container, day, attrName, lastWeek, yestday, today) {
    Highcharts.chart(container, {
        title: {
            text: day + ' 属性: ' + attrName
        },
        yAxis: {
            title: {
                text: attrName
            }
        },
        xAxis: {
            title: {
                text: '时间'
            },
            type: 'datetime',
            labels: {
                format: '{value:%H:%M}'
            },
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'middle'
        },
        plotOptions: {
            series: {
                label: {
                    connectorAllowed: false
                },
            }
        },
        series: [{
            name: '上周',
            data: lastWeek
        }, {
            name: '昨天',
            data: yestday
        }, {
            name: '今天',
            data: today
        }],
        responsive: {
            rules: [{
                condition: {
                    maxWidth: 500
                },
                chartOptions: {
                    legend: {
                        layout: 'horizontal',
                        align: 'center',
                        verticalAlign: 'bottom'
                    }
                }
            }]
        }
    });
}
