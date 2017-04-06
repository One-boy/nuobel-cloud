

window.onload = function () {
	
	//输入框置空
	login_form.username.value = "";
	login_form.pass.value = "";
	register_form.username.value="";
	register_form.pass.value="";
	register_form.repass.value="";
	register_form.email.value="";

}
//登录
function login() {

	//登录表单提交
	var username = login_form.username.value;
	var pass = login_form.pass.value;
	var autologin = login_form.autologin.checked;
	if (username == "") {
		showTips("login_tips", "用户名不能为空");
		hideTips("login_tips");
		return
	}
	if (pass == "") {
		showTips("login_tips", "密码不能为空");
		hideTips("login_tips");
		return
	}
	//提交
	var data = "username=" + username + "&pass=" + pass+"&autologin="+autologin;
	http.httpPost("/nuobel/login", data, function (resp) {
		if (resp === 1) {
			alert("网络错误");
			return
		}
		//解析json
		try {
			var jsonText = JSON.parse(resp)
			if (jsonText.code != 0) {
				alert(jsonText.mess)
				return
			}
			window.location.href="/nuobel/index";
		} catch (error) {
			alert("解析json错误")
			console.log(error)
		}

	})

}
//去注册
function toRegister() {
	//点击我要注册
	removeClass(document.getElementById('login_box'), "out")
	removeClass(document.getElementById('register_box'), "in")
	addClass(document.getElementById('login_box'), "in")
	addClass(document.getElementById('register_box'), "out");
}
//
function toLogin() {
	//点击返回登录
	removeClass(document.getElementById('login_box'), "in")
	removeClass(document.getElementById('register_box'), "out")
	addClass(document.getElementById('register_box'), "in")
	addClass(document.getElementById('login_box'), "out");
}
//注册
function register() {
	//注册表单提交
	var username = register_form.username.value;
	var pass = register_form.pass.value;
	var repass = register_form.repass.value;
	var email = register_form.email.value;
	if (username == "") {
		showTips("register_tips", "用户名不能为空");
		hideTips("register_tips");
		return false;
	}
	if (pass == "") {
		showTips("register_tips", "密码不能为空");
		hideTips("register_tips");
		return false;
	}
	if (repass == "") {
		showTips("register_tips", "重复密码不能为空");
		hideTips("register_tips");
		return false;
	}
	if (email == "") {
		showTips("register_tips", "邮箱不能为空");
		hideTips("register_tips");
		return false;
	}
	if (pass != repass) {
		showTips("register_tips", "两次密码不相同");
		hideTips("register_tips");
		return false;
	}
	//
	var data = "username=" + username + "&pass=" + pass + "&repass=" + repass + "&email=" + email;

	http.httpPost("/nuobel/register", data, function (resp) {
		if (resp === 1) {
			alert("网络错误")
			return
		}
		try {
			var jsonText = JSON.parse(resp)
			console.log(jsonText)
			if (jsonText.code != 0) {
				alert(jsonText.mess)
				console.log(jsonText.mess)
				return
			}
			alert("注册成功，马上去登录!")
			toLogin();
		} catch (error) {
			alert("解析json出错")
			console.error(error)
			return
		}
	})
}

//提示超时取消
function hideTips(type) {
	var tips = document.getElementById(type);
	if (!window.timeout1) {
		window.timeout1 = setTimeout(function () {
			tips.innerHTML = "";
			tips.style.display = "none";
			window.timeout1 = null;
		}, 2000);
	}

}
//显示提示
function showTips(type, text) {
	var tips = document.getElementById(type);
	tips.innerHTML = text;
	tips.style.display = "block";
}



//是否有这个class
function hasClass(obj, name) {
	return obj.className.match(new RegExp('(\\s|^)' + name + '(\\s|$)'));
}
//增加class
function addClass(obj, name) {
	if (!hasClass(obj, name)) {
		obj.className += " " + name;
	}
}
//removeclasss
function removeClass(obj, name) {
	if (hasClass(obj, name)) {
		var reg = new RegExp('(\\s|^)' + name + '(\\s|$)');
		obj.className = obj.className.replace(reg, ' ');
	}
}