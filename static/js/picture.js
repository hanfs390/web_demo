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
var sizePage = 20;
function showAllPicture() {
    $.ajax({
        type: "GET",
        url: "/picture/allPicture",
        dataType: "json",
        success: function (data) {
            var pictureHtml="";
            var allPictureArr = new Array(2000);
            var allPictureNumber = 0;
            var allCurrentPage = 0;
            var allPageNumber = 0;
            $.each(data, function(i, item) {
                allPictureNumber++;
                allPictureArr[i] = item;
            });
            allPageNumber = (allPictureNumber / sizePage) + 1;
            for (var i = 0; i < sizePage; i++) {
                var temp = allPictureArr[i].Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
            }

            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
            var pageHtml = "";
            for (var i = 0; i < allPageNumber; i++) {
                pageHtml = pageHtml + ` <li><a name="page" aria-label="`+ i.toString() +`">`+ (i+1).toString() +`</a></li>`;
            }
            $(".pagination").empty();
            $(".pagination").append(`<li>
                                <a name="page" aria-label="Previous">
                                    <span aria-hidden="true">&laquo;</span>
                                </a>
                            </li>
                            ` + pageHtml +`
                            <li>
                                <a name="page" aria-label="Next">
                                    <span aria-hidden="true">&raquo;</span>
                                </a>
                            </li>`);

            $("a[name=page]").click(function () {
                var flag = $(this).attr("aria-label");
                pictureHtml = "";
                if (flag == "Previous") {
                    if (allCurrentPage >= 1) {
                        allCurrentPage--;
                    } else {
                        allCurrentPage = 0;
                    }
                } else if (flag == "Next") {
                    if (allCurrentPage < allPageNumber) {
                        allCurrentPage++;
                    } else {
                        allCurrentPage = allPageNumber;
                    }
                } else {
                    allCurrentPage = parseInt(flag);
                }
                var offset = (allCurrentPage * sizePage) - 1;
                for (var i = offset; i < (offset + sizePage); i++) {
                    var temp = allPictureArr[i].Show
                    pictureHtml = pictureHtml +
                        `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
                }
                $("#showPictureList").empty();
                $("#showPictureList").append(pictureHtml);
            });
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
            var allPictureArr = new Array(2000);
            var allPictureNumber = 0;
            var allCurrentPage = 0;
            var allPageNumber = 0;
            $.each(data, function(i, item) {
                allPictureNumber++;
                allPictureArr[i] = item;
            });
            allPageNumber = (allPictureNumber / sizePage) + 1;
            for (var i = 0; i < sizePage; i++) {
                var temp = allPictureArr[i].Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
            }

            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
            var pageHtml = "";
            for (var i = 0; i < allPageNumber; i++) {
                pageHtml = pageHtml + ` <li><a name="page" aria-label="`+ i.toString() +`">`+ (i+1).toString() +`</a></li>`;
            }
            $(".pagination").empty();
            $(".pagination").append(`<li>
                                <a name="page" aria-label="Previous">
                                    <span aria-hidden="true">&laquo;</span>
                                </a>
                            </li>
                            ` + pageHtml +`
                            <li>
                                <a name="page" aria-label="Next">
                                    <span aria-hidden="true">&raquo;</span>
                                </a>
                            </li>`);

            $("a[name=page]").click(function () {
                var flag = $(this).attr("aria-label");
                pictureHtml = "";
                if (flag == "Previous") {
                    if (allCurrentPage >= 1) {
                        allCurrentPage--;
                    } else {
                        allCurrentPage = 0;
                    }
                } else if (flag == "Next") {
                    if (allCurrentPage < allPageNumber) {
                        allCurrentPage++;
                    } else {
                        allCurrentPage = allPageNumber;
                    }
                } else {
                    allCurrentPage = parseInt(flag);
                }
                var offset = (allCurrentPage * sizePage) - 1;
                for (var i = offset; i < (offset + sizePage); i++) {
                    var temp = allPictureArr[i].Show
                    pictureHtml = pictureHtml +
                        `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
                }
                $("#showPictureList").empty();
                $("#showPictureList").append(pictureHtml);
            });
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
            var allPictureArr = new Array(2000);
            var allPictureNumber = 0;
            var allCurrentPage = 0;
            var allPageNumber = 0;
            $.each(data, function(i, item) {
                allPictureNumber++;
                allPictureArr[i] = item;
            });
            allPageNumber = (allPictureNumber / sizePage) + 1;
            for (var i = 0; i < sizePage; i++) {
                var temp = allPictureArr[i].Show
                pictureHtml = pictureHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
            }

            $("#showPictureList").empty();
            $("#showPictureList").append(pictureHtml);
            var pageHtml = "";
            for (var i = 0; i < allPageNumber; i++) {
                pageHtml = pageHtml + ` <li><a name="page" aria-label="`+ i.toString() +`">`+ (i+1).toString() +`</a></li>`;
            }
            $(".pagination").empty();
            $(".pagination").append(`<li>
                                <a name="page" aria-label="Previous">
                                    <span aria-hidden="true">&laquo;</span>
                                </a>
                            </li>
                            ` + pageHtml +`
                            <li>
                                <a name="page" aria-label="Next">
                                    <span aria-hidden="true">&raquo;</span>
                                </a>
                            </li>`);

            $("a[name=page]").click(function () {
                var flag = $(this).attr("aria-label");
                pictureHtml = "";
                if (flag == "Previous") {
                    if (allCurrentPage >= 1) {
                        allCurrentPage--;
                    } else {
                        allCurrentPage = 0;
                    }
                } else if (flag == "Next") {
                    if (allCurrentPage < allPageNumber) {
                        allCurrentPage++;
                    } else {
                        allCurrentPage = allPageNumber;
                    }
                } else {
                    allCurrentPage = parseInt(flag);
                }
                var offset = (allCurrentPage * sizePage) - 1;
                for (var i = offset; i < (offset + sizePage); i++) {
                    var temp = allPictureArr[i].Show
                    pictureHtml = pictureHtml +
                        `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href="/picture/namePicture?name=`+allPictureArr[i].FileName+`" target="view_window" ><img src=`+allPictureArr[i].Url+`/1-front.jpg`+`></a>
                            <div class="caption">
                                <p>`+ temp+`</p>
                                <p>Label:`+allPictureArr[i].Label+`</p>
                            </div>
                        </div>
                    </div>`;
                }
                $("#showPictureList").empty();
                $("#showPictureList").append(pictureHtml);
            });
        }
    });
}
