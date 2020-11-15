function showHead() {
    $.ajax({
        type: "GET",
        url: "/picture/head",
        dataType: "json",
        success: function (data) {
            var navHtml="";
            $.each(data, function(i, item) {
                var labelCount = 0;
                var labelHtml = "";
                $.each(item.Label, function(j, label) {
                    labelCount++;
                    labelHtml = labelHtml + `<li><a name="label" onclick="showPictureByLabel(&quot;`+label+`&quot;)">`+label+`</a></li>`
                });
                if (labelCount > 0) {
                    navHtml = navHtml +
                        `<li class="dropdown">
                            <a name="dir" onclick="showPictureByDir(&quot;`+item.Dir+`&quot;)" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">`+ item.Dir +` <span class="caret"></span></a>
                            <ul class="dropdown-menu">
                                `+ labelHtml +`
                            </ul>
                        </li>`;
                } else {
                    navHtml = navHtml + `<li><a name="dir" onclick="showPictureByDir(&quot;`+item.Dir+`&quot;)">`+item.Dir+`</a></li>`;
                }
            });
            $("#navHead").append(navHtml);
        }
    });
}
function showAllPicture() {
    $.ajax({
        type: "GET",
        url: "/picture/allPicture",
        dataType: "json",
        success: function (data) {
            var pictureHtml="";
            $.each(data, function(i, item) {
                var temp = item.Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+item.FileName+`" target="view_window" ><img src=`+item.Url+`/1.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
        }
    });
}
function showPictureByDir(dir) {
    $.ajax({
        type: "GET",
        url: "/picture/dirPicture?dir="+dir,
        dataType: "json",
        success: function (data) {
            var pictureHtml="";
            $.each(data, function(i, item) {
                var temp = item.Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+item.FileName+`" target="view_window"><img src=`+item.Url+`/1.jpg`+`></a>
			                <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
        }
    });
}
function showPictureByLabel(label) {
    $.ajax({
        type: "GET",
        url: "/picture/labelPicture?label="+label,
        dataType: "json",
        success: function (data) {
            var pictureHtml="";
            $.each(data, function(i, item) {
                var temp = item.Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+item.FileName+`" target="view_window"><img src=`+item.Url+`/1.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+item.Label+`</p>
                            </div>
                        </div>
                    </div>`;
            });
            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
        }
    });
}
