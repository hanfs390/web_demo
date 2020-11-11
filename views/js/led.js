var ledOp = ""; /* judge the operation creat or modify */
var allGroupNameAndId = new Array;
var ledModifyId = 0;

$("#ledAdd").click(function() {
	ledOp = "Add";
	$("#ledList").css("display", "none");
	$("#ledConf").css("display", "block");
	$("#ledConf").empty(); /* empty the html */
	$("#ledConf").append(
		`<header class="panel-heading" style="overflow: auto">
			Create a New Led Timer Conf
			<button type="button" id="ledBack" class="btn btn-default" style="float: right;">Back</button>
		</header>`
	); /* add the led header */
	$("#ledConf").append(
		`<div class="panel-body">
			<form class="form-horizontal tasi-form" Name="ledConfForm">
				<div class="form-group">
					<label class="col-sm-2 control-label">GroupName</label>
					<div class="col-sm-10">
						<select class="form-control m-bot15" Name="ledGroupNameSelect">
						</select>
					</div>
				</div>
				<div class="form-group" Name="clock">
					<label class="col-sm-2 control-label">Start</label>
					<div class="col-sm-2">
						<div class="input-group clockpicker">
						    <input type="text" class="form-control" Name="start" value="09:30">
						    <span class="input-group-addon">
						        <span class="glyphicon glyphicon-time"></span>
						    </span>
						</div>
					</div>
					<div class="col-sm-2">
					</div>
					<label class="col-sm-2 control-label">Stop</label>
					<div class="col-sm-2">
						<div class="input-group clockpicker">
							<input type="text" class="form-control" Name="stop" value="09:30">
							<span class="input-group-addon">
								<span class="glyphicon glyphicon-time"></span>
							</span>
						</div>
					</div>
					<button data-dismiss="alert" class="close close-sm" type="button">
						<i class="icon-remove"></i>
					</button>
				</div>
				<button type="button" id="clockAdd" class="btn btn-success" style="float: right;">Add a Led Timer</button>
			</form>
		</div>
		<button type="button" class="btn btn-success btn-lg btn-block" id="ledCommitButton">Finish and Save</button>`
	)
	$('.clockpicker').clockpicker();
	//alert($("form[Name=ledConfForm] > div").size());
	$("#clockAdd").click(function() {
		$('#clockAdd').before(`<div class="form-group" Name="clock">
					<label class="col-sm-2 control-label">Start</label>
					<div class="col-sm-2">
						<div class="input-group clockpicker">
						    <input type="text" class="form-control" Name="start" value="09:30">
						    <span class="input-group-addon">
						        <span class="glyphicon glyphicon-time"></span>
						    </span>
						</div>
					</div>
					<div class="col-sm-2">
					</div>
					<label class="col-sm-2 control-label">Stop</label>
					<div class="col-sm-2">
						<div class="input-group clockpicker">
							<input type="text" class="form-control" Name="stop" value="09:30">
							<span class="input-group-addon">
								<span class="glyphicon glyphicon-time"></span>
							</span>
						</div>
					</div>
					<button data-dismiss="alert" class="close close-sm" type="button">
						<i class="icon-remove"></i>
					</button>
				</div>`)
	});
	$.ajax({
		type: "GET",
		url: "/group",
		dataType: "json",
		success: function(data) {
			allGroupNameAndId.length = 0;
			$.each(data, function(i, item) {
				var group = Object.create(null);
				group.Id = item.Id;
				group.Name = item.Name;
				$("select[name=ledGroupNameSelect]").append(`<option>` + group.Name + `</option>`);
				allGroupNameAndId.push(group);
			});
		},
		error: function(jqXHR, textStatus, errorThrown) {
			alert(jqXHR.responseText);
		},
	});
	$("#ledBack").click(function() {
		$("#ledList").css("display", "block");
		$("#ledConf").css("display", "none");
		$("#ledConf").empty();
	});
	$("#ledCommitButton").click(function() {
		var groupName;
		var TimerArr = new Array;
		$("form[Name=ledConfForm]").find("div[Name=clock]").each(function() {
			var start=$(this).find("input[Name=start]").val(),;
			temp = {
				"Start": $(this).find("input[Name=start]").val(),
				"Stop": $(this).find("input[Name=stop]").val(),
			};
			
		});
	});
});
/* del the led conf */
$("#ledDel").click(function() {
	$("#led_table").find(".checkboxes").each(function() {
		if ($(this).is((":checked"))) {
			wlanId = $(this).val();
			$.ajax({
				type: "DELETE",
				url: "/led?id=" + wlanId,
				dataType: "json",
				success: function(data) {
					window.location.reload();
				},
				error: function() {
					alert("请求失败");
				},
			});
		}
	});
});
