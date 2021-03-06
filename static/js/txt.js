function showHead() {
    $.ajax({
        type: "GET",
        url: "/txt/head",
        dataType: "json",
        success: function (data) {
            var navHtml="";
            $.each(data, function(i, item) {
                navHtml = navHtml + `<li><a name="dir" onclick="showTxtByDir(&quot;`+item.Dir+`&quot;)">`+item.Dir+`</a></li>`;
            });
            $("#navHead").append(navHtml);
        }
    });
}
function showAllTxt() {
    $.ajax({
        type: "GET",
        url: "/txt/allTxt",
        dataType: "json",
        success: function (data) {
            var txtHtml="";
            $.each(data, function(i, item) {
                var temp = item.FileName
                txtHtml = txtHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href=/txt/nameTxt?name=`+ temp +` target="view_window">`+temp+`</a>
                        </div>
                    </div>`;
            });
            $("#showTxtList").empty();
            $("#showTxtList").append(txtHtml);
        }
    });
}
function showTxtByDir(dir) {
    $.ajax({
        type: "GET",
        url: "/txt/dirTxt?dir="+dir,
        dataType: "json",
        success: function (data) {
            var txtHtml="";
            $.each(data, function(i, item) {
                var temp = item.FileName
                txtHtml = txtHtml +
                    `<div class="col-sm-6 col-md-4">
                        <div class="thumbnail">
                            <a href=/txt/nameTxt?name=`+ temp +` target="view_window">`+temp+`</a>
                        </div>
                    </div>`;
            });
            $("#showTxtList").empty();
            $("#showTxtList").append(txtHtml);
        }
    });
}

