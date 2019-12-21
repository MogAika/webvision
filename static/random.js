function wsRandomPage(data) {
	this.cachedBlocks = [];
	this.currentBlock = null;
	this.urlPrefix = '/data/';
	this.savedVolume = 0.5;
	this.cacheSize = 2;
}

wsRandomPage.prototype.preload = function() {
	for (let i = this.cachedBlocks.length; i < this.cacheSize; i++) {
		this.requestBlock();
	}
}

wsRandomPage.prototype.onNext = function() {
	if (this.currentBlock) {
		let plel = this.currentBlock.getPlayableElement();
		if (plel) {
			this.savedVolume = plel[0].volume;
		}

		this.currentBlock.div.remove();
		this.currentBlock = null;
	}

	this.currentBlock = this.cachedBlocks.shift();	
	if (this.currentBlock) {
		let plel = this.currentBlock.getPlayableElement();
		if (plel) {
			plel[0].volume = this.savedVolume;
			plel[0].play();
		}

		this.currentBlock.div.show();
	}
	
	this.preload();
}

wsRandomPage.prototype.onBlockReceived = function(data) {
	data.Url = this.urlPrefix + data.Url;
	let block = new wsRandomBlock(data);
	this.cachedBlocks.push(block);
	
	let elblock = block.createElement();
	$(".container").append(elblock);
	
	if (!this.currentBlock) {
		this.onNext();
	}
}

wsRandomPage.prototype.requestBlock = function() {
	let randomPage = this;
	$.ajax({
		url: "/api/random",
		method: "get",
		dataType: "json",
		fail: function(result) {
			console.log("fail", result);
		},
		success: function(data) {
			randomPage.onBlockReceived(data);
		}
	});
}


function wsRandomBlock(data) {
    this.data = data;
    this.div = undefined;
}

wsRandomBlock.prototype.getMediaType = function() {
    return this.data.Type.split('/')[0];
}

wsRandomBlock.prototype.getElement = function() {
	return this.div;
}

wsRandomBlock.prototype.createElement = function() {
	let type = this.getMediaType();
	let card;
	let o = this.data;
	switch (type) {
		case "audio":
			card = '<audio controls preload="none" loop><source src="' + o.Url + '"></audio>';
			break;
		case "video":
			card = '<video controls loop><source src="' + o.Url + '"></source></div>';
			break;
		case "image":
			card = '<a style="background-image:url(\'' + o.Url + '\'" href="' + o.Url + '" onclick="return false;">';
			break;
		default:
			card = '<h2>Unknown card type ' + type + '</h2>';
			break;
	}
	this.div = $('<div class="ws-random-card"><div class="ws-data ws-data-' + type + '">' + card + '</div></div>').hide();
	return this.div;
}

wsRandomBlock.prototype.getPlayableElement = function() {
	switch (this.getMediaType()) {
		case "video":
			return this.div.find("video");
		case "audio":
			return this.div.find("audio");
		default:
			return null;
	}
}

window.onload = function() {
	let randomPage = new wsRandomPage();
	$("#ws-next-video-btn").click(function() {
		randomPage.onNext();
	});
	
	randomPage.preload();
};
