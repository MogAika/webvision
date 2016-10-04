function wsProgressPercent(s, v) {
	s.find("progress").css('width', v+'%').val(v);
}

function wsProgressState(s, state, text) {
	s.find("progress").css('width', 100+'%').attr('value', 100).removeClass("progress-info").addClass("progress-" + state);
	s.find(".file_status").removeClass("text-primary").removeClass("text-info").addClass("text-" + state).text(text);
}

function wsProgressError(s, err) {
	wsProgressState(s, "danger", "Error: " + err);
}

function wsProgressUploaded(s) {
	wsProgressState(s, "success", "Uploaded");
}

function wsGetUploadElement(fname) {
	return $('<div class="upload">\
			<div class="row">\
				<div class="col-xs-6 text-xs-right file_name">' + fname + '</div>\
				<div class="col-xs-6 text-xs-left text-primary file_status">Uploading</div>\
			</div>\
			<div class="row">\
				<div class="col-xs-12">\
					<progress class="progress-striped progress" value="0" max="100"></progress>\
				</div>\
			</div>\
		</div>');
}

function wsUploadFile(file) {
	var fd = new FormData();
	
	fd.append('fl', file);
	
	var row = wsGetUploadElement(file.name);
	
	$("#uploads").append(row);
	$.ajax({
		xhr: function() {
			var xhr = new window.XMLHttpRequest();
			xhr.upload.addEventListener("progress", function(evt) {
				if (evt.lengthComputable) {
					var percentComplete = evt.loaded / evt.total;
					percentComplete = percentComplete * 100;
					wsProgressPercent(row, percentComplete);
				}
			}, false);
			return xhr;
		},
		url: "/api/upload",
		dataType: "json",
		method: "post",
		data: fd,
		processData: false,
		contentType: false, 
		mimeType: 'multipart/form-data', 
		fail: function(result) {
			wsProgressError(row, result);
		},
		success: function(data) {
			if (data.error) {
				wsProgressError(row, data.error);
			} else {
				wsProgressUploaded(row);
			}
		},
	});
}

function wsUploadOnChange() {
	if (this.files.length) {
		$("#ws-upload-select").text("Selected " + this.files.length + " files");
	} else {
		$("#ws-upload-select").text("Select files");
	}
}

function wsUploadSelectFiles() {
	$("#ws-upload-file").trigger('click');
}

function wsUpload() {
	var files = $("#ws-upload-file")[0].files;
	
	if (!files.length) {
		return;
	}
	
	for (var i = 0; i < files.length; i++) {
		wsUploadFile(files[i]);
	}
}

function isValidURL(str) {
   var a  = document.createElement('a');
   a.href = str;
   return (a.host && a.host != window.location.host);
}

function wsUploadByUrlWarn(text) {
	$warn = $("#ws-upload-url-warn");
	if (text) {
		$warn.text(text);
		$warn.show();
		$warn.parent().addClass("has-danger");
	} else {
		$warn.hide();
		$warn.parent().removeClass("has-danger");
	}
}

function wsUploadOnUrlSubmit() {
	var url = $("#ws-upload-url").val();
	if (!isValidURL(url)) {
		wsUploadByUrlWarn('Not a valid url');
		return;
	}
	wsUploadByUrlWarn();
	
	wsUploadByUrl(url);
}

function wsUploadByUrl(url) {
	var row = wsGetUploadElement(url);
	$("#uploads").append(row);
	wsProgressState(row, "info", "Server downloading file");
	
	$.ajax({
		url: "/api/upload",
		dataType: "json",
		data: {'url': url},
		method: "get",
		fail: function(result) {
			wsProgressError(row, result);
		},
		success: function(data) {
			if (data.error) {
				wsProgressError(row, data.error);
			} else {
				wsProgressUploaded(row);
			}
		},
	});
}

$(function() {
	$("#ws-upload-file").change(wsUploadOnChange);
	$("#ws-upload-select").click(wsUploadSelectFiles);
	$("#ws-upload-url").keyup(function(event) {
	    if(event.keyCode == 13) { wsUploadOnUrlSubmit(); }
	});
	$("#ws-upload-url-btn").click(wsUploadOnUrlSubmit);
});