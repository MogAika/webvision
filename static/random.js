function wsBlockShow($bl) {
	$bl.show();
	var video = $bl.find("video").first();
	if (video.length) {
		video[0].play();
	}
}

function wsNextVideoClick() {
	$(".ws-random-card").first().remove();
	var nextVideo = $(".ws-random-card").first();
	if (nextVideo.length) {
		wsBlockShow(nextVideo);
	}
	wsRequestRandomBlock();
}

function wsCreateRandomBlock(o) {
	var ptype = o.Type.split('/')[0];
	switch (ptype) {
		case "audio":
			var card = '<audio controls preload="none" loop><source src="' + o.Url + '"></audio>';
			break;
		case "video":
			var card = '<video controls loop><source src="' + o.Url + '"></source></div>';
			break;
		case "image":
			var card = '<a style="background-image:url(\'' + o.Url + '\'" href="' + o.Url + '" onclick="return false;">';
			break;
	}
	var block = $('<div class="ws-random-card"><div class="ws-data ws-data-' + ptype + '">' + card + '</div></div>').hide();
	
	$(".container").append(block);
	if ($(".ws-random-card").length === 1) {
		wsBlockShow(block);
	}

	return block;
}

function wsRequestRandomBlock() {
	$.ajax({
		url: "/api/random",
		method: "get",
		dataType: "json",
		fail: function(result) {
			console.log("fail", result);
		},
		success: function(data) {
			data.Url = "/data/" + data.Url;
			wsCreateRandomBlock(data);
		}
	});
}

window.onload = function() {
	wsRequestRandomBlock();
	wsRequestRandomBlock();
	$("#ws-next-video-btn").click(wsNextVideoClick);
};