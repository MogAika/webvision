wsBrowseRequested = false;
wsBrowseEnd = false;
wsBrowseLoaded = -1;
wsLastPlayedVideo = null;

wsGetVideoBlock = function(url, adds) {
	return '<video controls loop maximized ' + adds + '><source src="' + url + '"></video>';
}

wsBrowseInsert = function(o) {
	var ptype = o.Type.split('/')[0];
	switch (ptype) {
		case "audio":
			var card = '<audio controls preload="none" loop><source src="' + o.Url + '"></audio>';
			break;
		case "video":
			if (o.Thumb != null) {
				var card = '<div class="ws-data-lazyvideo"><a href="' + o.Url + '"><img src="' + o.Thumb + '"/></a></div>';
			} else {
				var card = wsGetVideoBlock(o.Url, 'preload="meta"');
			}
			break;
		case "image":
			var card = '<img src="' + o.Url + '">';
			break;
	}
	card = $('<div class="ws-card"><div class="ws-data ws-data-' + ptype + '">' + card + '</div></div>');
	card.insertBefore($("#ws-request-trigger"));
	card.find(".ws-data-lazyvideo").click(wsLazyVideoOnClick);
}

wsRequestMedia = function() {
	wsBrowseRequested = true;
	var data = wsBrowseLoaded != -1 ? ("s=" + wsBrowseLoaded) : "";
	$.ajax({
		url: "/",
		method: "get",
		dataType: "json",
		data: data,
		fail: function(result) {
			console.log("fail", result);
			wsBrowseRequested = false;
		},
		success: function(data) {
			for (var i in data) {
				var obj = data[i];
				if (obj.Id < wsBrowseLoaded || wsBrowseLoaded < 0) {
					wsBrowseLoaded = obj.Id;
				}
				wsBrowseInsert(obj);
			}
			wsBrowseRequested = false;
		},
	});
}

wsLazyVideoOnClick = function(ev) {
	var current = $(ev.target).parent().parent().parent();
	var src = current.find("a").attr("href");
	
	if (wsLastPlayedVideo != null) {
		var player = wsLastPlayedVideo.find("video");
		player.appendTo(current);
		player.find("source").attr('src', src);
		player[0].load();
		wsLastPlayedVideo.find(".ws-data-lazyvideo").show();
	} else {
		current.append(wsGetVideoBlock(src, 'autoplay'));
	}
	
	wsLastPlayedVideo = current;
	wsLastPlayedVideo.find(".ws-data-lazyvideo").hide();
	return false;
};

$(document).ready(function() {
	wsRequestMedia();
	$(document).scroll(function() {
		if (!wsBrowseEnd && !wsBrowseRequested) {
			if (($(window).scrollTop() + $(window).height()) >= $("#ws-request-trigger").position().top - 512) {
				wsRequestMedia();
			}
		}
	});
});