function generateContentViews() {
    var svg = d3.select("#content-views"),
        margin = {top: 20, right: 20, bottom: 30, left: 38},
        width = svg.attr("width") - margin.left - margin.right,
        height = svg.attr("height") - margin.top - margin.bottom;

    var x = d3.scaleTime().range([0, width]),
        y = d3.scaleLinear().range([height, 0]),
        z = d3.scaleOrdinal().range(["#393b79", "#5254a3" , "#6b6ecf"]);

    var stack = d3.stack();

    var area = d3.area()
        .x(function(d, i) { return x(d.data.date); })
        .y0(function(d) { return y(d[0]); })
        .y1(function(d) { return y(d[1]); });

    var g = svg.append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    d3.csv("../data?q=content_view", type, function(error, data) {
        if (error) throw error;

        var keys = data.columns.slice(1);

        x.domain(d3.extent(data, function(d) { return d.date; }));
        z.domain(keys);
        stack.keys(keys);

        var layer = g.selectAll(".layer")
            .data(stack(data))
            .enter().append("g")
            .attr("class", "layer");

        layer.append("path")
            .attr("class", "area")
            .style("fill", function(d) { return z(d.key); })
            .attr("d", area);

        layer.filter(function(d) { return d[d.length - 1][1] - d[d.length - 1][0] > 0.01; })
            .append("text")
            .attr("class", "label")
            .attr("x", width - 6)
            .attr("y", function(d) { return y((d[d.length - 1][0] + d[d.length - 1][1]) / 2); })
            .attr("dy", ".35em")
            .style("font-size", "10px")
            .style("text-anchor", "end")
            .text(function(d) { return d.key; });

        g.append("g")
            .attr("class", "axis axis--x")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x).ticks(5));

        g.append("g")
            .attr("class", "axis axis--y")
            .call(d3.axisLeft(y).ticks(5, "%"));
    });
}

function type(d, i, columns) {
    var parseDate = d3.timeParse("%Y %b %d");
    d.date = parseDate(d.date);
    for (var i = 1, n = columns.length; i < n; ++i) d[columns[i]] = d[columns[i]] / 100;
    return d;
}

function generateMAUs() {
    var svg = d3.select("#mau"),
        margin = {top: 20, right: 20, bottom: 30, left: 38},
        width = svg.attr("width") - margin.left - margin.right,
        height = svg.attr("height") - margin.top - margin.bottom,
        g = svg.append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    var parseTime = d3.timeParse("%Y %b %d");

    var x = d3.scaleTime()
        .rangeRound([0, width]);

    var y = d3.scaleLinear()
        .rangeRound([height, 0]);

    var area = d3.area()
        .x(function(d) { return x(d.date); })
        .y1(function(d) { return y(d.count); });

    d3.csv("../data?q=maus", function(d) {
        d.date = parseTime(d.date);
        d.count = +d.count;
        return d;
    }, function(error, data) {
        if (error) throw error;

        x.domain(d3.extent(data, function(d) { return d.date; }));
        y.domain([0, d3.max(data, function(d) { return d.count; })]);
        area.y0(y(0));

        g.append("path")
            .datum(data)
            .attr("fill", "#5254a3")
            .attr("d", area);

        g.append("g")
            .attr("class", "axis axis--x")
            .attr("transform", "translate(0," + height + ")")
            .call(d3.axisBottom(x).ticks(5));

        g.append("g")
            .attr("class", "axis axis--y")
            .call(d3.axisLeft(y).ticks(5));
    });
}

generateContentViews();
generateMAUs();