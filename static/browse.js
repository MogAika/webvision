wsBrowseRequested = false;
wsBrowseEnd = false;
wsBrowseLoaded = -1;

wsBrowseInsert = function(o) {
	var ptype = o.Type.split('/')[0];
	switch (ptype) {
		case "audio":
			var card = '<audio controls preload="none" loop><source src="' + o.Url + '"></audio>';
			break;
		case "video":
			var card = '<video controls preload="none" loop maximized><source src="' + o.Url + '"></video>';
			break;
		case "image":
			var card = '<img src="' + o.Url + '">';
			break;
	}
	console.log(card);
	card = $('<div class="ws-card"><div class="ws-data ws-data-' + ptype + '">' + card + '</div></div>');
	card.insertBefore($("#ws-request-trigger"));
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

$(function() {
	wsRequestMedia();
	$(document).scroll(function() {
		if (!wsBrowseEnd && !wsBrowseRequested) {
			if (($(window).scrollTop() + $(window).height()) >= $("#ws-request-trigger").position().top - 512) {
				wsRequestMedia();
			}
		}
	});
});