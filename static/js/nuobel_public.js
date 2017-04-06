//http方法
var http = {
	httpGet: function (url, callback) {
		var xhr = new XMLHttpRequest();
		var url = url;
		xhr.timeout = 5000;
		xhr.ontimeout = function(){
			callback(1);
		}
		xhr.open("GET", url, true);
		//回复
		xhr.onreadystatechange = function () {
			if (xhr.readyState == 4 && xhr.status == 200) {
				callback(xhr.responseText);
			} else if (xhr.readyState == 4 && xhr.status != 200) {
				callback(1);
			}
		}
		xhr.send();
	},
	//data:uid=xxx&sid=xxx
	httpPost: function (url, data, callback) {
		var xhr = new XMLHttpRequest();
		var url = url;
		var data = data;
		xhr.open("POST", url, true);
		//回复
		xhr.onreadystatechange = function () {
			if (xhr.readyState == 4 && xhr.status == 200) {
				callback(xhr.responseText);
			} else if (xhr.readyState == 4 && xhr.status != 200) {
				callback(1);
			}
		}

		//post要加下面这一个请求头信息,open过后才能设置
		xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
		xhr.send(data);
	},
	//file
	httpPostFile: function (obj) {
		var xhr = new XMLHttpRequest();
		var url = obj.url;
		var data = obj.data;
		//回复
		xhr.onreadystatechange = function () {
			if (xhr.readyState == 4 && xhr.status == 200) {
				obj.callback(xhr.responseText,obj.id);
			} else if (xhr.readyState == 4 && xhr.status != 200) {
				obj.callback(1,obj.id);
			}
		}
		xhr.upload.onprogress = function (evt) {
			var loaded = evt.loaded;
	
			var total = evt.total;
			var per = Math.floor(100 * loaded / total);
			obj.progress(per,obj.id);
		}
		xhr.open("POST", url, true);
		xhr.send(data);
	}
}