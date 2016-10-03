wsUpoladProgress = function(s, v) {
	s.find("progress").css('width', v+'%').attr('value', v);
}
wsEndProgress = function(s) {
	s.find("progress").css('width', 100+'%').attr('value', 100).addClass("progress-success");
	s.find(".file_status").removeClass("text-primary").addClass("text-success").text("Uploaded");
}
wsErrorProgress = function(s, v) {
	s.find("progress").css('width', 100+'%').attr('value', 100).addClass("progress-danger");
	s.find(".file_status").removeClass("text-primary").addClass("text-danger").text("Error:" + v);
}

wsGetUploadElement = function(fname) {
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

wsUploadFile = function(file) {
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
					wsUpoladProgress(row, percentComplete);
				}
			}, false);
			return xhr;
		},
		url: "/upload",
		method: "post",
		data: fd,
		processData: false,
		contentType: false, 
		mimeType: 'multipart/form-data', 
		fail: function(result) {
			wsErrorProgress(row, result);
		},
		success: function(data) {
			if (data != "") {
				wsErrorProgress(row, data);
			} else {
				wsEndProgress(row);
			}
		},
	});
}

wsUploadOnChange = function(obj) {
	if (obj.files.length) {
		$("#ws-upload-select").text("Selected " + obj.files.length + " files");
	} else {
		$("#ws-upload-select").text("Select files");
	}
}

wsUploadSelectFiles = function(obj) {
	$("#ws-upload-file").trigger('click');
}

wsUpload = function(obj) {
	var files = $("#ws-upload-file")[0].files;
	
	if (!files.length) {
		return;
	}
	
	for (var i = 0; i < files.length; i++) {
		wsUploadFile(files[i]);
	}
}