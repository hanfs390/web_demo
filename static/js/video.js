function showHead() {
    $.ajax({
        type: "GET",
        url: "/video/head",
        dataType: "json",
        success: function (data) {
            var navHtml="";
            $.each(data, function(i, item) {
                var labelCount = 0;
                var labelHtml = "";
                $.each(item.Label, function(j, label) {
                    labelCount++;
                    labelHtml = labelHtml + `<li><a name="label" onclick="showVideoByLabel(&quot;`+label+`&quot;)">`+label+`</a></li>`
                });
                if (labelCount > 0) {
                    navHtml = navHtml +
                        `<li class="dropdown">
                            <a name="dir" onclick="showVideoByDir(&quot;`+item.Dir+`&quot;)" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">`+ item.Dir +` <span class="caret"></span></a>
                            <ul class="dropdown-menu">
                                `+ labelHtml +`
                            </ul>
                        </li>`;
                } else {
                    navHtml = navHtml + `<li><a name="dir" onclick="showVideoByDir(&quot;`+item.Dir+`&quot;)">`+item.Dir+`</a></li>`;
                }
            });
            $("#navHead").append(navHtml);
        }
    });
}
function showAllVideo() {
    $.ajax({
        type: "GET",
        url: "/video/allVideo",
        dataType: "json",
        success: function (data) {
            var videoHtml="";
            $.each(data, function(i, item) {
                var temp = item.FileName
                videoHtml = videoHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href=`+item.Url+`><img src=`+item.Picture+`></a>
                            <div class="caption">
                                <p>Name:`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showVideoList").empty();
            $("#showVideoList").append(videoHtml);
        }
    });
}
function showVideoByDir(dir) {
    $.ajax({
        type: "GET",
        url: "/video/dirVideo?dir="+dir,
        dataType: "json",
        success: function (data) {
            var videoHtml="";
            $.each(data, function(i, item) {
                var temp = item.FileName
                videoHtml = videoHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href=`+item.Url+`><img src=`+item.Picture+`></a>
                            <div class="caption">
                                <p>Name:`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showVideoList").empty();
            $("#showVideoList").append(videoHtml);
        }
    });
}
function showVideoByLabel(label) {
    $.ajax({
        type: "GET",
        url: "/video/labelVideo?label="+label,
        dataType: "json",
        success: function (data) {
            var videoHtml="";
            $.each(data, function(i, item) {
                var temp = item.FileName
                videoHtml = videoHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href=`+item.Url+`><img src=`+item.Picture+`></a>
                            <div class="caption">
                                <p>Name:`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showVideoList").empty();
            $("#showVideoList").append(videoHtml);
        }
    });
}
