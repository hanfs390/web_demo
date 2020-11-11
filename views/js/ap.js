var allGroupNameAndId = new Array
/* create the ap */
function createAP() {
	var mac = $.trim($("#mac").val());
	if (mac == "") {
		alert("ap mac must be not empty");
	}
	var tag = $.trim($("#tag").val());
	var detail = $.trim($("#detail").val());
	var model = $.trim($("#model").val());
	data = {
		"Mac": mac,
		"Model": model,
		"Tag": tag,
		"Detail": detail
	};
	dataStr = JSON.stringify(data);

	$.ajax({
		url: "/ap",
		type: "POST",
		data: dataStr,
		success: function(data) {
			window.location.reload()
		},
		error: function() {
			alert("请求失败");
		},
		dataType: "json",
		contentType: "application/json"
	});
}
/* show the create ap web */
$("#apAdd").click(function() {
	$("#apList").css("display", "none");
	$("#apNew").css("display", "block");

	$("#apBack").click(function() {
		$("#apList").css("display", "block");
		$("#apNew").css("display", "none");
	});
});
/* del the ap */
$("#apDel").click(function() {
	$("#ap_table").find(".checkboxes").each(function() {
		if ($(this).is((":checked"))) {
			apId = $(this).val();
			$.ajax({
				type: "DELETE",
				url: "/ap?id=" + apId,
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


var modifyAPDevCfg = Object.create(null); /* the ap dev cfg */
var currentBand = ""

/* callback of save modify */
function updateAp(id) {
	var mac = $.trim($("input[name=Mac]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var tag = $.trim($("input[name=Tag]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var groupName = $.trim($("select[name=apGroupNameSelect]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var alias = $.trim($("input[name=Alias]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var detail = $.trim($("input[name=Detail]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var blocked = $.trim($("input[name=Blocked]").val().replace(/(^\s*)|(\s*$)/g, ""));
	var groupid;
	for (var i = 0; i < allGroupNameAndId.length; i++) {
		if (groupName == allGroupNameAndId[i].Name) {
			groupid = allGroupNameAndId[i].Id;
		}
	}
	if (currentBand == "BAND_2G") {
		modifyAPDevCfg.Channel2G = $("#channel").val()
		modifyAPDevCfg.TxPower2G = $("#txpower").val()
		modifyAPDevCfg.HtMode2G = $("#htmode").val()
		modifyAPDevCfg.HwMode2G = $("#mode").val()
	} else if (currentBand == "BAND_5G") {
		modifyAPDevCfg.Channel5G1 = $("#channel").val()
		modifyAPDevCfg.TxPower5G1 = $("#txpower").val()
		modifyAPDevCfg.HtMode5G1 = $("#htmode").val()
		modifyAPDevCfg.HwMode5G1 = $("#mode").val()
	} else if (currentBand == "BAND_5G2") {
		modifyAPDevCfg.Channel5G2 = $("#channel").val()
		modifyAPDevCfg.TxPower5G2 = $("#txpower").val()
		modifyAPDevCfg.HtMode5G2 = $("#htmode").val()
		modifyAPDevCfg.HwMode5G2 = $("#mode").val()
	}
	var wired = $("#wired").children();
	for (var i=0; i < wired.length; i++) {
		 if (i==0) {
			modifyAPDevCfg.LAN1=$("#port1").val();
		 }
		 if (i==1) {
		 	modifyAPDevCfg.LAN2=$("#port2").val();
		 }
		 if (i==2) {
		 	modifyAPDevCfg.LAN3=$("#port3").val();
		 }
		 if (i==3) {
			modifyAPDevCfg.LAN4=$("#port4").val();
		 }
	}
	data = {
		"Id": id,
		"Mac": mac,
		"GroupId": groupid,
		"GroupName": groupName,
		"Alias": alias,
		"Tag": tag,
		"Detail": detail,
		"Blocked": parseInt(blocked),
		"DevCfg": JSON.stringify(modifyAPDevCfg)
	};
	
	dataStr = JSON.stringify(data);
	alert(dataStr);
	$.ajax({
		url: "/ap",
		type: "PUT",
		data: dataStr,
		success: function(data) {
			window.location.reload()
		},
		error: function() {
			alert("请求失败");
		},
		dataType: "json",
		contentType: "application/json"
	});
}
/* show the radio and wired web*/
function generateRadioAndWired(model) {
	$.ajax({
		type: "GET",
		url: "/ap/model?model=" + model,
		dataType: "json",
		success: function(data) {
			/* band head */
			var bandarr = new Array
			if (data.Radio & 0x1) {
				bandarr.push("BAND_2G");
			}
			if (data.Radio & (0x1 << 1)) {
				bandarr.push("BAND_5G");
			}
			if (data.Radio & (0x1 << 2)) {
				bandarr.push("BAND_5G2");
			}
			var bandButtonsHTML;
			var temp = "";
			for (var i=0; i < bandarr.length; i++) {
				temp += `<button class="btn btn-white" type="button" id="`+bandarr[i]+`">`+bandarr[i]+`</button>`;
			}
			bandButtonsHTML = `<div class="btn-row">
							<div class="btn-group">`+ temp +`
							</div>
							</div>`;
			$("#radio").append(bandButtonsHTML)
			$("#radio").append(`<div class="form-group">
                                      <label class="control-label col-lg-2" for="inputSuccess">Channel</label>
                                      <div class="col-lg-10">
                                          <select class="form-control m-bot15" id="channel">
                                          </select>
                                      </div>
                                  </div>
								  <div class="form-group">
                                      <label class="control-label col-lg-2" for="inputSuccess">HTmode</label>
                                      <div class="col-lg-10">
                                          <select class="form-control m-bot15" id="htmode">
                                              <option>open</option>
                                          </select>
                                      </div>
                                  </div>
								  <div class="form-group">
                                      <label class="control-label col-lg-2" for="inputSuccess">TxPower</label>
                                      <div class="col-lg-10">
                                          <select class="form-control m-bot15" id="txpower">
                                              <option>open</option>
                                          </select>
                                      </div>
                                  </div><div class="form-group">
                                      <label class="control-label col-lg-2" for="inputSuccess">Mode</label>
                                      <div class="col-lg-10">
                                          <select class="form-control m-bot15" id="mode">
                                              <option>open</option>
                                          </select>
                                      </div>
                                  </div>`)
			/* radio */
			for (var i=0; i < bandarr.length; i++) {
				$("#"+bandarr[i]).click(bandarr[i], function(band) {
					var hopeBand = band.data
					if (currentBand == hopeBand) {
						return;
					}
					if (currentBand == "BAND_2G") {
						$("#"+currentBand).removeClass("btn-black");
						$("#"+currentBand).addClass("btn-white");
						
						modifyAPDevCfg.Channel2G = $("#channel").val()
						modifyAPDevCfg.TxPower2G = $("#txpower").val()
						modifyAPDevCfg.HtMode2G = $("#htmode").val()
						modifyAPDevCfg.HwMode2G = $("#mode").val()
					} else if (currentBand == "BAND_5G") {
						$("#"+currentBand).removeClass("btn-black");
						$("#"+currentBand).addClass("btn-white");
						
						modifyAPDevCfg.Channel5G1 = $("#channel").val()
						modifyAPDevCfg.TxPower5G1 = $("#txpower").val()
						modifyAPDevCfg.HtMode5G1 = $("#htmode").val()
						modifyAPDevCfg.HwMode5G1 = $("#mode").val()
					} else if (currentBand == "BAND_5G2") {
						$("#"+currentBand).removeClass("btn-black");
						$("#"+currentBand).addClass("btn-white");
						
						modifyAPDevCfg.Channel5G2 = $("#channel").val()
						modifyAPDevCfg.TxPower5G2 = $("#txpower").val()
						modifyAPDevCfg.HtMode5G2 = $("#htmode").val()
						modifyAPDevCfg.HwMode5G2 = $("#mode").val()
					}
					if (hopeBand == "BAND_2G") {
						/* init channel */
						$("#channel").empty()
						$("#channel").append(`<option>auto</option>`);
						for (var i=0; i < data.Chan2G.length; i++) {
							$("#channel").append(`<option>`+data.Chan2G[i]+`</option>`);
						}
						$("#channel").val(modifyAPDevCfg.Channel2G);
						/* init txpower */
						$("#txpower").empty()
						$("#txpower").append(`<option>auto</option>`);
						for (var i=0; i < 25; i++) {
							$("#txpower").append(`<option>`+i+`</option>`);
						}
						$("#txpower").val(modifyAPDevCfg.TxPower2G);
						
						/* init htmode */
						$("#htmode").empty()
						$("#htmode").append(`<option>auto</option>`);
						$("#htmode").append(`<option>ht20</option>`);
						$("#htmode").append(`<option>ht40</option>`);
						
						$("#htmode").val(modifyAPDevCfg.HtMode2G);
						
						/* init mode */
						$("#mode").empty()
						$("#mode").append(`<option>auto</option>`);
						$("#mode").append(`<option>11ng</option>`);
						$("#mode").append(`<option>11n</option>`);
						$("#mode").append(`<option>11g</option>`);
						
						$("#mode").val(modifyAPDevCfg.HwMode2G);						
					} else if (hopeBand == "BAND_5G") {
						/* init channel */
						$("#channel").empty()
						$("#channel").append(`<option>auto</option>`);
						for (var i=0; i < data.Chan5G.length; i++) {
							$("#channel").append(`<option>`+data.Chan5G[i]+`</option>`);
						}
						$("#channel").val(modifyAPDevCfg.Channel5G1);
						/* init txpower */
						$("#txpower").empty()
						$("#txpower").append(`<option>auto</option>`);
						for (var i=0; i < 25; i++) {
							$("#txpower").append(`<option>`+i+`</option>`);
						}
						$("#txpower").val(modifyAPDevCfg.TxPower5G1);
						
						/* init htmode */
						$("#htmode").empty()
						$("#htmode").append(`<option>auto</option>`);
						$("#htmode").append(`<option>ht20</option>`);
						$("#htmode").append(`<option>ht40</option>`);
						$("#htmode").append(`<option>ht80</option>`);
						
						$("#htmode").val(modifyAPDevCfg.HtMode5G1);
						
						/* init mode */
						$("#mode").empty()
						$("#mode").append(`<option>auto</option>`);
						$("#mode").append(`<option>11ac</option>`);
						$("#mode").append(`<option>11an</option>`);
						
						$("#mode").val(modifyAPDevCfg.HwMode5G1);						
					} else if (hopeBand == "BAND_5G2") {
						/* init channel */
						$("#channel").empty()
						$("#channel").append(`<option>auto</option>`);
						for (var i=0; i < data.Chan5G2.length; i++) {
							$("#channel").append(`<option>`+data.Chan5G2[i]+`</option>`);
						}
						$("#channel").val(modifyAPDevCfg.Channel5G2);
						/* init txpower */
						$("#txpower").empty()
						$("#txpower").append(`<option>auto</option>`);
						for (var i=0; i < 25; i++) {
							$("#txpower").append(`<option>`+i+`</option>`);
						}
						$("#txpower").val(modifyAPDevCfg.TxPower5G2);
						
						/* init htmode */
						$("#htmode").empty()
						$("#htmode").append(`<option>auto</option>`);
						$("#htmode").append(`<option>ht20</option>`);
						$("#htmode").append(`<option>ht40</option>`);
						$("#htmode").append(`<option>ht80</option>`);
						
						$("#htmode").val(modifyAPDevCfg.HtMode5G2);
						
						/* init mode */
						$("#mode").empty()
						$("#mode").append(`<option>auto</option>`);
						$("#mode").append(`<option>11ac</option>`);
						$("#mode").append(`<option>11an</option>`);
						
						$("#mode").val(modifyAPDevCfg.HwMode5G2);						
					}
					currentBand = hopeBand;
					$("#"+currentBand).addClass("btn-black");
					$("#"+currentBand).removeClass("btn-white");
				});
			}
			$("#"+bandarr[0]).click(); /* do a click */
			/* wired */
			for (var i=0; i < data.PortNumber; i++) {
				var port = i+1;
				$("#wired").append(`<div class="form-group">
									    <label for="inputEmail1" class="col-lg-2 control-label">Port_`+port+`</label>
									    <div class="col-lg-10">
									        <input type="Number" class="form-control" id="port`+port+`" onblur ="javascript:if(this.value<0){
this.value=0;};if(this.value>4094){this.value=4094;}">
									        <p class="help-block">VLANID (1-4094) 0 或者 空 代表不设置VLAN</p>
									    </div>
									</div>`)
				if (i==0) {
					$("#port1").val(modifyAPDevCfg.LAN1);
				}
				if (i==1) {
					$("#port2").val(modifyAPDevCfg.LAN2);
				}
				if (i==2) {
					$("#port3").val(modifyAPDevCfg.LAN3);
				}
				if (i==3) {
					$("#port4").val(modifyAPDevCfg.LAN4);
				}
			}
		},
		error: function() {
			alert("请求失败");
		},
	});
}
/* open modify web */
$("#apUpdate").click(function() {
	$("#ap_table").find(".checkboxes").each(function() {
		if ($(this).is((":checked"))) {
			apId = $(this).val();
			$.ajax({
				type: "GET",
				url: "/ap/id?id=" + apId,
				dataType: "json",
				success: function(data) {
					/* show the Group web */
					/* init global variables */
					modifyAP=data;
					modifyAPDevCfg=Object.create(null);
					currentBand="";
					if (data.DevCfg != "") {
						modifyAPDevCfg=JSON.parse(data.DevCfg);
					}
					
					$("#apList").css("display", "none");
					$("#apModify").css("display", "block");
					$("#apModify").append(
						`<header class="panel-heading" style="overflow: auto">
                                Update a New AP
                                <button type="button" id="apModifyBack" class="btn btn-default" style="float: right;">Back</button>
                        </header>
						<div class="panel-body">
                            <form class="form-horizontal tasi-form" id="modifyAPBaiseInfo">
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">Mac</label>
                                    <div class="col-sm-10">
                                        <input type="text" class="form-control" name="Mac" value="` + data.Mac + `" readonly="readonly">
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">GroupName</label>
                                    <div class="col-sm-10">
                                        <select class="form-control m-bot15" Name="apGroupNameSelect">
											<option>default</option>
                                        </select>
                                    </div>
                                </div>
                                <div class="form-group">
									<label class="col-sm-2 control-label">Alias</label>
									<div class="col-sm-10">
										<input type="text" class="form-control" name="Alias" value="` + data.Alias + `">
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">Tag</label>
                                    <div class="col-sm-10">
                                        <input type="text" class="form-control" name="Tag" value="` +data.Tag +`">
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">Detail</label>
                                    <div class="col-sm-10">
                                        <input type="text" class="form-control" name="Detail" value="` +data.Detail +`" >
                                    </div>
                                </div>
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">Blocked</label>
                                    <div class="col-sm-10">
                                        <input type="text" class="form-control" name="Blocked" value="` +data.Blocked + `">
                                    </div>
                                </div>
                            </form>
                        </div>
                        <div class="row">
							<div class="col-lg-6">
								<section class="panel">
									<header class="panel-heading">Radio</header>
									<div class="panel-body" id="radio"></div>
								</section>
							</div>
							<div class="col-lg-6">
								<section class="panel">
									<header class="panel-heading">Wired</header>
									<div class="panel-body" id="wired"></div>
								</section>
							</div>
						</div>
						<button type="button" class="btn btn-success btn-lg btn-block" onclick="updateAp(` + apId + `)">Finish and Save</button>`);

					/* get radio and wired info from board */
					generateRadioAndWired(data.Model)
					
					$("#apModifyBack").click(function() { /* close the modify web */
						$("#apList").css("display", "block");
						$("#apModify").css("display", "none");
						$("#apModify").empty()
					});
					$.ajax({
						type: "GET",
						url: "/group",
						dataType: "json",
						success: function(g) {
							allGroupNameAndId.length = 0;
							$.each(g, function(i, item) {
								var group = Object.create(null);
								group.Id = item.Id;
								group.Name = item.Name;
								$("select[name=apGroupNameSelect]").append(`<option>` + group.Name + `</option>`);
								allGroupNameAndId.push(group);
								$("select[name=apGroupNameSelect]").val(data.GroupName);
							});
						},
						error: function(jqXHR, textStatus, errorThrown) {
							alert(jqXHR.responseText);
						},
					});
				},
				error: function() {
					alert("请求失败");
				},
			});
			return false /* return each */
		}
	});
});

/* show the detail of ap */
function showApInfo(value) {
	$.ajax({
		type: "GET",
		url: "/ap/id?id=" + value,
		dataType: "json",
		success: function(data) {
			/* show the Group web */
			$("#apList").css("display", "none");
			$("#apUpdate").css("display", "block");

			$("#apBack").click(function() {
				$("#apList").css("display", "block");
				$("#apUpdate").css("display", "none");
			});
		},
		error: function() {
			alert("请求失败");
		},
	});
}