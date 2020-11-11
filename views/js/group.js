
/* show the create group web */
$("#groupAdd").click(function() {
	$("#groupList").css("display", "none");
	$("#groupNew").css("display", "block");
	
	$("#groupBack").click(function() {
		$("#groupList").css("display", "block");
		$("#groupNew").css("display", "none");
	});
});
/* del the group */
$("#groupDel").click(function() {
    $("#group_table").find(".checkboxes").each(function () {
       if ($(this).is((":checked"))) {
           groupId = $(this).val();
           $.ajax({
               type: "DELETE",
               url: "/group?id="+groupId,
               dataType: "json",
               success:function(data){
                   window.location.reload();
               },
               error:function(){
                   alert("请求失败");
               },
           });
       }
    });
});
/* update the group */
$("#groupUpdate").click(function() {
    $("#group_table").find(".checkboxes").each(function () {
        if ($(this).is((":checked"))) {
            groupId = $(this).val();
            $.ajax({
                type: "GET",
                url: "/group/id?id="+groupId,
                dataType: "json",
                success:function(data){
                    /* show the Group web */
                    $("#groupList").css("display", "none");
                    $("#groupModify").css("display", "block");
                    $("#groupModify").append(`                            <div class="panel-body">
                                <form class="form-horizontal tasi-form">
                                    <div class="form-group">
                                        <label class="col-sm-2 control-label">Group Name</label>
                                        <div class="col-sm-10">
                                            <input type="text" class="form-control" value="`+ data.Name +`" readonly="readonly">
                                        </div>
                                    </div>
                                    <div class="form-group">
                                        <label class="col-sm-2 control-label">Tag</label>
                                        <div class="col-sm-10">
                                            <input type="text" id="modifyTag" class="form-control" value="`+ data.Tag +`">
                                        </div>
                                    </div>
                                    <div class="form-group">
                                        <label class="col-sm-2 control-label">Detail</label>
                                        <div class="col-sm-10">
                                            <input type="text" id="modifyDetail" class="form-control" value="`+ data.Detail +`">
                                        </div>
                                    </div>
                                </form>
                            </div>
                            <button type="button" class="btn btn-success btn-lg btn-block" onclick="updateGroup(`+groupId+`)">Finish and Save</button>`);
                    $("#groupModifyBack").click(function() {
                        $("#groupList").css("display", "block");
                        $("#groupModify").css("display", "none");
                        window.location.reload();
                    });
                },
                error:function(){
                    alert("请求失败");
                },
            });
            return false /* return each */
        }
    });
});
/* create the group */
function createGroup() {
    var name = $.trim($("#groupName").val());
    if (name == "") {
        alert("group name must be not empty");
    }
    var tag = $.trim($("#tag").val());
    var detail = $.trim($("#detail").val());

    data = {"Name":name, "Tag":tag, "Detail":detail};
    dataStr = JSON.stringify(data);

    $.ajax({
        url:"/group",
        type:"POST",
        data:dataStr,
        success:function(data){
            window.location.reload()
        },
        error:function(){
            alert("请求失败");
        },
        dataType:"json",
        contentType : "application/json"
    });
}
function updateGroup(id) {
    var tag = $.trim($("#modifyTag").val());
    var detail = $.trim($("#modifyDetail").val());
    data = {"Id":id, "Tag":tag, "Detail":detail};
    dataStr = JSON.stringify(data);

    $.ajax({
        url:"/group",
        type:"PUT",
        data:dataStr,
        success:function(data){
            window.location.reload()
        },
        error:function(){
            alert("请求失败");
        },
        dataType:"json",
        contentType : "application/json"
    });
}
/* show the detail of group */
function showGroupInfo(value) {
    $.ajax({
        type: "GET",
        url: "/group/id?id="+value,
        dataType: "json",
        success:function(data){
            /* show the Group web */
            $("#groupList").css("display", "none");
            $("#groupUpdate").css("display", "block");

            $("#groupBack").click(function() {
                $("#groupList").css("display", "block");
                $("#groupUpdate").css("display", "none");
            });
        },
        error:function(){
            alert("请求失败");
        },
    });
}