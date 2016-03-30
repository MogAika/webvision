wsUpoladProgress = function(v) {
	$("#ws-upload-progress").css('width', v+'%').attr('value', v).text(parseInt(v) + " %")
		.removeClass("progress-success").removeClass("progress-danger").addClass("progress-info");
}
wsEndProgress = function() {
	$("#ws-upload-progress").css('width', 100+'%').attr('value', 100).text("Uploaded")
		.removeClass("progress-danger").removeClass("progress-info").addClass("progress-success");
}
wsErrorProgress = function(v) {
	$("#ws-upload-progress").css('width', 100+'%').attr('value', 100).text(v)
		.removeClass("progress-success").removeClass("progress-info").addClass("progress-danger");
}

wsUpload = function(obj) {
	var fd = new FormData();
	var files = $("#ws-upload-file")[0].files;
	if (!files.length) {
		return;
	}
	fd.append('heh', files[0]);

	wsUpoladProgress(0);
	$.ajax({
		xhr: function() {
			var xhr = new window.XMLHttpRequest();
			xhr.upload.addEventListener("progress", function(evt) {
				if (evt.lengthComputable) {
					var percentComplete = evt.loaded / evt.total;
					percentComplete = percentComplete * 100;
					wsUpoladProgress(percentComplete);
				}
			}, false);
			return xhr;
		},
		url: "/upload",
		type: "POST",
		data: fd,
		processData: false,
		contentType: false, 
		mimeType: 'multipart/form-data', 
		fail: function(result) {
			console.log(result);
			wsErrorProgress(result);
		},
		success: function(data) {
			console.log(data);
			if (data != "") {
				wsErrorProgress(data);
			} else {
				wsEndProgress();
			}
		},
	});
}