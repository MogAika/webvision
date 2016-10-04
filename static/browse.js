wsBrowseRequested = false;
wsBrowseLoaded = -1;
wsBrowseEnd = false;
wsLastPlayedVideo = null;

function wsGetVideoBlock(url, adds) {
	return '<video controls loop maximized ' + adds + '><source src="' + url + '"></video>';
}

function wsBrowseInsert(o) {
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
	if (ptype == "video") {
		card.find(".ws-data-lazyvideo").click(wsLazyVideoOnClick);
	}
}

function wsRequestMedia() {
	if (wsBrowseRequested) { return; }
	wsBrowseRequested = true;
	var data = {'count': 25};
	if (wsBrowseLoaded !== -1) {
		data['start'] = wsBrowseLoaded;
	}
	$.ajax({
		url: "/api/query",
		method: "get",
		dataType: "json",
		data: data,
		fail: function(result) {
			console.log("fail", result);
			wsBrowseRequested = false;
		},
		success: function(data) {
			if (data.length == 0) {
				wsBrowseEnd = true;
			} else {
				for (var i in data) {
					var obj = data[i];
					if (obj.Id < wsBrowseLoaded || wsBrowseLoaded < 0) {
						wsBrowseLoaded = obj.Id;
						$("#ws-request-trigger").show();
					}
					data[i].Url = "/data/" + data[i].Url;
					if (data[i].Thumb != null) {
						data[i].Thumb = "/data/" + data[i].Thumb;
					}
					wsBrowseInsert(obj);
				}
			}
			wsBrowseRequested = false;
		},
	});
}

function wsLazyVideoOnClick(ev) {
	var current = $(ev.target).parent();
	while (!current.hasClass("ws-data-video")) {
		current = current.parent();
	}
	
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

$(function() {
	wsRequestMedia();
	$(document).scroll(function() {
		if (!wsBrowseRequested && !wsBrowseEnd) {
			if (($(window).scrollTop() + $(window).height()) >= $("#ws-request-trigger").position().top - 512) {
				wsRequestMedia();
			}
		}
	});
});
