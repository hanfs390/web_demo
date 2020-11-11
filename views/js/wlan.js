var allGroupNameAndId = new Array

var wlanop = "";
var wlanModifyId = 0; //the ID of modify wlan
/* create the wlan */
$("#wlanCreateAndModifyButton").click(function() {
	var groupname = $.trim($("select[name=wlanGroupNameSelect]").val());
	var groupid;
	for (var i = 0; i < allGroupNameAndId.length; i++) {
		if (groupname == allGroupNameAndId[i].Name) {
			groupid = allGroupNameAndId[i].Id;
		}
	}
	var ssid = $.trim($("input[name=wlanSsidInput]").val());
	var disabled_str = $.trim($("button[name=wlanDisabledButton]").text());
	var disabled = 0;
	if (disabled_str == "YES") {
		disabled = 1;
	}
	var hidden = $.trim($("button[name=wlanHiddenButton]").text());

	var radio = "";
	if ($.trim($("button[name=wlan2GButton]").text()) == "2G enabled") {
		radio += "2G;"
	}
	if ($.trim($("button[name=wlan5GButton]").text()) == "5G enabled") {
		radio += "5G;"
	}
	if ($.trim($("button[name=wlan5G2Button]").text()) == "5G2 enabled") {
		radio += "5G2;"
	}
	var encryption = $.trim($("select[name=wlanEncryptionSelect]").val());
	var key = "";
	var authserver = "";
	var authsecret = "";
	var authport = "";
	if (encryption == "NONE") {

	} else if (encryption == "WPA2-PSK+TKIP") {
		key = $.trim($("input[name=wlanKeyInput]").val());
	} else if (encryption == "WPA2-PSK+CCMP") {
		key = $.trim($("input[name=wlanKeyInput]").val());
	} else if (encryption == "WPA2-EAP+CCMP") {
		authserver = $.trim($("input[name=wlanServerIPInput]").val());
		authport = $.trim($("input[name=wlanServerPortInput]").val());
		authsecret = $.trim($("input[name=wlanServerKeyInput]").val());
	}
	var vlanswitch_str = $.trim($("button[name=wlanVlanButton]").text());
	var vlanswitch;
	var vlanid = 0;
	if (vlanswitch_str == "enabled") {
		vlanid = parseInt($("input[name=wlanVlanIdInput]").val());
		vlanswitch = 1;
	} else {
		vlanswitch = 0;
	}

	data = {
		"GroupName": groupname,
		"GroupId": groupid,
		"Ssid": ssid,
		"Disabled": disabled,
		"Hidden": hidden,
		"Radio": radio,
		"Encryption": encryption,
		"Key": key,
		"AuthSecret": authsecret,
		"AuthServer": authserver,
		"AuthPort": authport,
		"VlanSwitch": vlanswitch,
		"VlanId": vlanid,
	};

	var method = "";
	if (wlanop == "Add") {
		method = "POST";
	} else if (wlanop == "Update") {
		method = "PUT";
		data.Id = wlanModifyId;
	}
	dataStr = JSON.stringify(data);

	$.ajax({
		url: "/wlan",
		type: method,
		data: dataStr,
		success: function(data) {
			window.location.reload();
		},
		error: function(jqXHR, textStatus, errorThrown) {
			alert(jqXHR.responseText);
		},
		dataType: "json",
		contentType: "application/json"
	});
});
/* show the create wlan web */
$("#wlanAdd").click(function() {
	wlanop = "Add";
	$("#wlanList").css("display", "none");
	$("#wlanNew").css("display", "block");
	$("#wlanHeader").remove();
	$("#wlanNew").prepend(
		`<header class="panel-heading" style="overflow: auto" id="wlanHeader">
									Create a new WLAN
									<button type="button" id="wlanBack" class="btn btn-default" style="float: right;">Back</button>
								</header>`
	);
	$("select[name=wlanGroupNameSelect]").empty();

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
				$("select[name=wlanGroupNameSelect]").append(`<option>` + group.Name + `</option>`);
				allGroupNameAndId.push(group);
			});
		},
		error: function(jqXHR, textStatus, errorThrown) {
			alert(jqXHR.responseText);
		},
	});
	$("#wlanBack").click(function() {
		$("#wlanList").css("display", "block");
		$("#wlanNew").css("display", "none");
		window.location.reload();
	});
});
$("#wlanDel").click(function() {
	$("#wlan_table").find(".checkboxes").each(function() {
		if ($(this).is((":checked"))) {
			wlanId = $(this).val();
			$.ajax({
				type: "DELETE",
				url: "/wlan?id=" + wlanId,
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
/* show the modify wlan web */
$("#wlanUpdate").click(function() {
	$("#wlan_table").find(".checkboxes").each(function() {
		if ($(this).is((":checked"))) {
			wlanModifyId = parseInt($(this).val());
			wlanop = "Update";
			$("select[name=wlanGroupNameSelect]").empty();
			$("#wlanList").css("display", "none");
			$("#wlanNew").css("display", "block");
			$("#wlanHeader").remove();
			$("#wlanNew").prepend(
				`<header class="panel-heading" style="overflow: auto" id="wlanHeader">
											Update a WLAN
											<button type="button" id="wlanBack" class="btn btn-default" style="float: right;">Back</button>
										</header>`
			);
			$("#wlanBack").click(function() {
				$("#wlanList").css("display", "block");
				$("#wlanNew").css("display", "none");
				window.location.reload();
			});
			wlanId = $(this).val();
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
						$("select[name=wlanGroupNameSelect]").append(`<option>` + group.Name + `</option>`);
						allGroupNameAndId.push(group);
					});

					$.ajax({
						type: "GET",
						url: "/wlan/id?id=" + wlanId,
						dataType: "json",
						success: function(data) {
							/* show the Group web */
							$("select[name=wlanGroupNameSelect]").val(data.GroupName);
							$("input[name=wlanSsidInput]").val(data.Ssid);
							var disabled;
							if (data.Disabled == "1") {
								disabled = "YES";
							} else {
								disabled = "NO";
							}
							$("button[name=wlanDisabledButton]").text(disabled);
							var hidden;
							if (data.Hidden == "enable") {
								hidden = "disable";
							} else {
								hidden = "enable";
							}
							$("button[name=wlanHiddenButton]").text(hidden);
							if (data.Radio.search("2G") != -1) {
								$("button[name=wlan2GButton]").text("2G enabled");
							} else if (data.Radio.search("5G") != -1) {
								$("button[name=wlan5GButton]").text("5G enabled");
							} else if (data.Radio.search("5G2") != -1) {
								$("button[name=wlan5G2Button]").text("5G2 enabled");
							}
							$("select[name=wlanEncryptionSelect]").val(data.Encryption);

							if (data.Encryption == "NONE") {} else if ((data.Encryption == "WPA2-PSK+CCMP") || (data.Encryption ==
									"WPA2-PSK+TKIP")) {
								$("#wlanEncryptionCfg").empty();
								$("#wlanEncryptionCfg").css("display", "block");
								$("#wlanEncryptionCfg").append(
									`<label class="col-sm-2 control-label">Password</label>
																<div class="col-sm-10">
																	<input type="Password" Name="wlanKeyInput" class="form-control">
																</div>`
								);
								$("input[name=wlanKeyInput]").val(data.Key);
							} else if (data.Encryption == "WPA2-EAP+CCMP") {
								$("#wlanEncryptionCfg").empty();
								$("#wlanEncryptionCfg").css("display", "block");
								$("#wlanEncryptionCfg").append(
									`<label class="col-sm-2 control-label">Server IP</label>
																<div class="col-sm-10">
																	<input type="text" Name="wlanServerIPInput" class="form-control">
																</div>`
								);
								$("#wlanEncryptionCfg").append(
									`<label class="col-sm-2 control-label">Server Port</label>
																<div class="col-sm-10">
																	<input type="text" Name="wlanServerPortInput" class="form-control">
																</div>`
								);
								$("#wlanEncryptionCfg").append(
									`<label class="col-sm-2 control-label">Password</label>
																<div class="col-sm-10">
																	<input type="Password" Name="wlanServerKeyInput" class="form-control">
																</div>`
								);
								$("input[name=wlanServerIPInput]").val(data.AuthServer);
								$("input[name=wlanServerPortInput]").val(data.AuthPort);
								$("input[name=wlanServerKeyInput]").val(data.AuthSecret);
							}
							var vlanswitch;
							if (data.VlanSwitch == 0) {
								vlanswitch = "disabled";
							} else {
								vlanswitch = "enabled";
								$("#wlanVlanCfg").empty();
								$("#wlanVlanCfg").css("display", "block");
								$("#wlanVlanCfg").append(
									`<label class="col-sm-2 control-label">VLAN ID</label>
																<div class="col-sm-10">
																	<input type="Number" class="form-control" name="wlanVlanIdInput" onblur ="javascript:if(this.value<0){
																	this.value=0;};if(this.value>4094){this.value=4094;}">
																</div>`
								);
								$("input[name=wlanVlanIdInput]").val(data.VlanId);
							}
							$("button[name=wlanVlanButton]").text(vlanswitch);
						},
						error: function(jqXHR, textStatus, errorThrown) {
							alert(jqXHR.responseText);
						},
					});
				},
				error: function(jqXHR, textStatus, errorThrown) {
					alert(jqXHR.responseText);
				},
			});
			return false /* return each */
		}
	});
});

function initWlan() {
	$("button[name=wlanDisabledButton]").click(function() {
		if ($(this).text() == "NO") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("YES");
		} else if ($(this).text() == "YES") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("NO");
		}
	});
	$("button[name=wlanHiddenButton]").click(function() {
		if ($(this).text() == "disable") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("enable");
		} else if ($(this).text() == "enable") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("disable");
		}
	});
	$("button[name=wlan2GButton]").click(function() {
		if ($(this).text() == "2G enabled") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("2G disabled");
		} else if ($(this).text() == "2G disabled") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("2G enabled");
		}
	});
	$("button[name=wlan5GButton]").click(function() {
		if ($(this).text() == "5G enabled") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("5G disabled");
		} else if ($(this).text() == "5G disabled") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("5G enabled");
		}
	});
	$("button[name=wlan5G2Button]").click(function() {
		if ($(this).text() == "5G2 enabled") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("5G2 disabled");
		} else if ($(this).text() == "5G2 disabled") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("5G2 enabled");
		}
	});
	$("select[name=wlanEncryptionSelect]").change(function() {
		var entrypt = $(this).val();
		if (entrypt == "NONE") {
			$("#wlanEncryptionCfg").empty();
			$("#wlanEncryptionCfg").css("display", "none");
		} else if ((entrypt == "WPA2-PSK+TKIP") || (entrypt == "WPA2-PSK+CCMP")) {
			$("#wlanEncryptionCfg").empty();
			$("#wlanEncryptionCfg").css("display", "block");
			$("#wlanEncryptionCfg").append(
				`<label class="col-sm-2 control-label">Password</label>
											<div class="col-sm-10">
												<input type="Password" Name="wlanKeyInput" class="form-control">
											</div>`
			);
		} else if (entrypt == "WPA2-EAP+CCMP") {
			$("#wlanEncryptionCfg").empty();
			$("#wlanEncryptionCfg").css("display", "block");
			$("#wlanEncryptionCfg").append(
				`<label class="col-sm-2 control-label">Server IP</label>
											<div class="col-sm-10">
												<input type="text" Name="wlanServerIPInput" class="form-control">
											</div>`
			);
			$("#wlanEncryptionCfg").append(
				`<label class="col-sm-2 control-label">Server Port</label>
											<div class="col-sm-10">
												<input type="text" Name="wlanServerPortInput" class="form-control">
											</div>`
			);
			$("#wlanEncryptionCfg").append(
				`<label class="col-sm-2 control-label">Password</label>
											<div class="col-sm-10">
												<input type="Password" Name="wlanServerKeyInput" class="form-control">
											</div>`
			);
		}
	});
	$("button[name=wlanVlanButton]").click(function() {
		if ($(this).text() == "enabled") {
			$(this).removeClass("btn-white");
			$(this).addClass("btn-black");
			$(this).text("disabled");
			$("#wlanVlanCfg").empty();
			$("#wlanVlanCfg").css("display", "none");
		} else if ($(this).text() == "disabled") {
			$(this).removeClass("btn-black");
			$(this).addClass("btn-white");
			$(this).text("enabled");
			$("#wlanVlanCfg").empty();
			$("#wlanVlanCfg").css("display", "block");
			$("#wlanVlanCfg").append(
				`<label class="col-sm-2 control-label">VLAN ID</label>
											<div class="col-sm-10">
												<input type="Number" class="form-control" name="wlanVlanIdInput" onblur ="javascript:if(this.value<0){
												this.value=0;};if(this.value>4094){this.value=4094;}">
											</div>`
			);
		}
	});
}
