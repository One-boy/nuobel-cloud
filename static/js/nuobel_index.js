//获取用户信息

//工具箱
var tools = (function () {
    //尺寸转换
    var sizeToEye = function (size) {
        var rs = size / 1024  //KB
        var result;
        rs = Math.ceil(rs)
        if (rs > 1000) {
            rs = rs / 1024 //MB
            rs = Math.ceil(rs)
            result = rs + "MB"
        } else {
            result = rs + "KB"
        }
        return result
    }
    //顶部小提示
    var topTips = function (text) {
        document.getElementById("tips_display").style.top = "0rem";
        document.getElementById("tips_content").innerText = text;
        setTimeout(function () {
            document.getElementById("tips_display").style.top = "-1.5rem";
        }, 1000)
    }
    //动态设置内容框大小
    var resizeList = function () {
        document.getElementById("content_file").parentNode.style.height = window.innerHeight - 190 + "px";
    }
    //获取当前filepath
    var getFilepath = function () {
        //根据历史纪录
        return document.getElementById("history_path").lastElementChild.getAttribute("title");
    }
    //秒数转换成分钟
    var secondToMin = function(sec){
        if (sec <= 1){
            return "00:01"
        }
        var minutes = Math.floor(sec/60);
        if (minutes < 10){
            minutes = "0"+minutes;
        }
        var second =  Math.ceil(sec%60);
        if (second < 10){
            second = "0"+second;
        }
        return minutes+":"+second
    }
    return {
        SizeToEye: sizeToEye,
        Tips: topTips,
        ResizeList: resizeList,
        GetFilepath: getFilepath,
        SecondToMin:secondToMin
    }
})()
//用户弹出框
var userPopBox = (function () {
    //初始化，事件监听
    var init = function () {
        document.getElementById("index_user_box").addEventListener("mouseover", mouseover, false);
        document.getElementById("index_user_box").addEventListener("mouseout", mouseout, false);
        document.getElementById("exit").addEventListener("click", userHandle.ExitLogin, false);
    }
    var mouseover = function () {
        document.getElementById("user_pop").style.visibility = "visible";
        document.getElementById("user_pop").style.height = "1.5rem";
    }
    var mouseout = function () {
        document.getElementById("user_pop").style.visibility = "hidden";
        document.getElementById("user_pop").style.height = "0rem";
    }
    return {
        init: init
    }
})()
//左侧导航栏
var leftNavigation = (function () {
    var datatype = "all";//当前所在导航栏
    //初始化，事件监听
    var init = function () {
        document.getElementById("left_file_type").addEventListener("click", leftTypeHandle, false)
    }
    //左边导航栏点击事件
    var leftTypeHandle = function (evt) {
        var target = evt.target;
        if (!target.getAttribute("title")){
            target = target.parentNode;
        }
        var child = target.parentNode.parentNode.childNodes;
        for (var i = 0; i < child.length; i++) {
            child[i].className = "";
        }
        target.parentNode.className = "checked";
        datatype = target.getAttribute("data-type");
        
        historyModule.setHistory(target.getAttribute("title"))
        
        getList(target.getAttribute("data-type"), "");
    }
    //获取当前导航栏
    var getDataType = function () {
        return datatype;
    }
    return {
        init: init,
        GetDataType: getDataType
    }
})()
//工具栏点击事件分发
var toolBar = (function () {
    //文件操作对象
    var init = function () {
        //全选事件
        document.getElementById("checkall").addEventListener("click", checkall, false);
        //工具栏事件
        document.getElementById("tool_box").addEventListener("click", tool_bar_router, false);
    }

    var checkall = function (evt) {  //全选操作
        var isChecked = evt.target.checked;
        var items = document.getElementsByName("checkitem");
        if (isChecked) {
            for (var i = 0; i < items.length; i++) {
                items[i].checked = true;
            }
        } else {
            for (var i = 0; i < items.length; i++) {
                items[i].checked = false;
            }
        }
    }
    var tool_bar_router = function (evt) {   //工具栏点击路由分发
        var action = evt.target.getAttribute("action");
        //操作分发
        switch (action) {
            case "makedir":
                if (leftNavigation.GetDataType() != "all") {
                    tools.Tips("分类导航不能新建")
                    return
                }
                document.getElementsByClassName("pop_box")[0].style.display = "block";  //显示创建框
                document.getElementById("pop_btn").addEventListener("click", fileHandle.makedir, false); //监听创建框按钮
                break;
            case "download":
                download();
                break;
            case "delete":
                deletefile();
                break;
            case "refresh":
                fileHandle.refresh();
                break;
            case "rename":
                fileHandle.rename();
                break;
        }
    }
    var download = function () {  //下载文件预操作
        //遍历选中的文件
        var items = document.getElementsByName("checkitem");
        var willdown = {};
        var parent, filepath;
        for (var i = 0; i < items.length; i++) {
            if (items[i].checked) {
                parent = items[i].parentNode.parentNode;
                if (parent.getAttribute("data-type") === "dir") {
                    tools.Tips("暂不支持文件夹下载")
                    return
                }
                willdown[parent.getAttribute("filename")] = parent.getAttribute("filepath")
            }
        }
        if (!(function(){   //判断对象是否为空
            for (var i in willdown)
                return true
            return false
        })()) {
            tools.Tips("你未选择任何文件")
            return
        }
        
        for (var filename in willdown) {
            fileHandle.download(willdown[filename], filename);
        }
    }
    var deletefile = function () { //删除文件预操作
        //遍历选中的文件
        var items = document.getElementsByName("checkitem");
        var willdown = [];
        var parent, filepath;
        for (var i = 0; i < items.length; i++) {
            parent = items[i].parentNode.parentNode;
            if (items[i].checked) {
                willdown.push(parent.getAttribute("filename"))
            }
        }
        if (willdown.length === 0) {
            tools.Tips("您还未选择任何文件")
            return
        }
        filepath = parent.getAttribute("filepath");
        willdown.forEach(function (filename) {
            fileHandle.deletefile(filepath, filename);
        }, this);
    }

    return {
        init: init
    }
})()

//文件操作,包括新建文件夹，删除，重命名，下载,刷新等
var fileHandle = (function () {
    var listStage = document.getElementById("content_file");

    var filepath;
    var refresh = function () {  //刷新数据
        if (leftNavigation.GetDataType() != "all") {
            getList(leftNavigation.GetDataType(), "");
            return
        }
        filepath = tools.GetFilepath();
        getList("all", filepath);
    }
    var makedir = function (evt) { //新建文件夹
        filepath = tools.GetFilepath()
        var action = evt.target.getAttribute("action"); //获取动作
        if (action === "cancel") {  //取消，隐藏弹出框
            document.getElementsByClassName("pop_box")[0].style.display = "none";
            return
        }

        var newdirname = document.getElementById("newdirname").value;//获取新文件夹名称
        if (newdirname.length <= 0) {
            tools.Tips("你还未输入名称")
            return
        }
        var post_data = "path=" + filepath + "&newdirname=" + newdirname;
        http.httpPost("/nuobel/makedir", post_data, function (resp) {
            if (resp === 1) {
                tools.Tips("网络错误")
                document.getElementsByClassName("pop_box")[0].style.display = "none";
                return
            }
            var jsonText = JSON.parse(resp);
            if (jsonText.code != 0) {
                tools.Tips("创建失败")
                console.error(jsonText.mess)
            } else {
                tools.Tips("创建成功")
            }

            document.getElementsByClassName("pop_box")[0].style.display = "none"; //隐藏弹出框
            refresh();  //刷新数据
        })

    }
    var download = function (filepath, filename) { //下载文件
        var iframe = document.createElement("iframe")
        iframe.src = "/downloadfile?path=" + filepath + "&filename=" + encodeURIComponent(filename);
        iframe.style.display = "none";
        //文件下载方式
        document.body.appendChild(iframe);
    }
    var deletefile = function (filepath, filename) {  //删除文件
        var post_data = "path=" + filepath + "&filename=" + encodeURIComponent(filename);
        http.httpPost("/nuobel/deletefile", post_data, function (resp) {
            if (resp === 1) {
                tools.Tips("网络错误");
                return
            }
            var jsonText = JSON.parse(resp);
            if (jsonText.code != 0) {
                tools.Tips("抱歉，不能删除非空文件夹");
                console.error(jsonText.mess)
            } else {
                tools.Tips("删除成功");
                refresh();  //刷新数据
            }
        })
    }
    var rename = function () { //重命名
        //遍历选中的文件
        var items = document.getElementsByName("checkitem");
        var parent;
        for (var i = 0, j = 0; i < items.length; i++) {
            if (items[i].checked) {
                j++;
                parent = items[i].parentNode.parentNode;
            }
            if (j > 1) {
                tools.Tips("只能对单个文件重命名")
                return
            }
        }
        if (!parent) {
            tools.Tips("你未选择任何文件")
            return
        }
        //创建改名输入框
        var input = document.createElement("input")
        input.type = "text"
        input.value = parent.getAttribute("filename")
        input.className = "rename_input";
        input.name = "rename_input";

        parent.appendChild(input)
        input.focus(); //聚焦
        input.select(); //全选
        //离开时，提交更改
        var oldname = parent.getAttribute("filename");
        var filepath = tools.GetFilepath();
        var newname;
        input.onblur = function (evt) {
            newname = evt.target.value;
            if (newname.length < 1 || oldname == newname) {  //未填写或者未改变，不提交更改
                input.remove();
                return;
            }
            if (newname.match(/:/)) {
                tools.Tips("文件名不能包含:特殊字符")
                return
            }
            var data = "filepath=" + filepath + "&oldname=" + encodeURIComponent(oldname) + "&newname=" + newname;
            http.httpPost("/nuobel/rename", data, function (resp) {
                if (resp === 1) {
                    tools.Tips("网络错误")
                    return
                }
                var jsonText = JSON.parse(resp);
                if (jsonText.code == 0) {
                    tools.Tips("重命名成功")
                } else {
                    tools.Tips("重命名失败")
                }
                refresh();
            })
        }
    }
    return {
        refresh: refresh,
        makedir: makedir,
        download: download,
        deletefile: deletefile,
        rename: rename
    }
})()


//用户相关操作
var userHandle = (function () {
    //获取用户信息
    var getUserInfo = function () {
        http.httpGet("/nuobel/getuserinfo?time=" + (new Date()).getTime() / 1000, function (resp) {
            if (resp === 1) {
                tips("网络错误");
                return
            }
            var jsonText = JSON.parse(resp);
            if (jsonText.code != 0) {
                tips("抱歉，获取用户信息失败");
                console.error(jsonText.mess)
                return
            }
            document.getElementById("index_name").innerText = jsonText.data.username;
        })
    }
    //退出登录
    var exitLogin = function () {
        http.httpGet("/nuobel/exit?time=" + (new Date()).getTime() / 1000, function (resp) {
            if (resp === 1) {
                tips("网络错误")
                return
            }
            var jsonText = JSON.parse(resp);
            if (jsonText.code == 0) {
                window.location.href = "/nuobel"
            } else {
                tips("退出失败")
                console.error(jsonText.mess)
            }
        })
    }

    return {
        GetUserInfo: getUserInfo,
        ExitLogin: exitLogin
    }
})()
//上传状态框
var uploadBox = (function () {
    var init = function () {
        document.getElementById("open_statusbox").addEventListener("click", uploadBoxHandle, false);//监听状态框按钮
    }
    var uploadBoxHandle = function (evt) {
        var node = evt.target;
        var height = node.parentNode.parentNode.offsetHeight;
        var parent = node.parentNode.parentNode;
        var status = document.getElementById("open_statusbox");
        if (height > 40) {
            parent.style.height = "1.7rem"
            status.innerText = "+"
        } else {
            parent.style.height = "15rem"
            status.innerText = "-"
        }
    }
    return {
        init: init
    }
})()
//列表显示区相关操作方法集合
var listStage = function () {
    this.mainNode = document.getElementById("content_file");  //列表显示区主节点
    var p = document.createElement("p");  //创建一个提示元素
    p.style.fontSize = "1.3rem"
    p.style.textAlign = "center"
    this.p = p;
}
listStage.prototype.clearList = function () {  //清空舞台数据
    this.mainNode.innerHTML = "";
}
listStage.prototype.loading = function () { //显示加载中信息
    this.clearList()
    this.p.innerText = "正在加载数据..."
    this.mainNode.appendChild(this.p)
}
listStage.prototype.showtips = function (text) { //提示信息
    this.clearList()
    this.p.innerText = text
    this.mainNode.appendChild(this.p)
}
listStage.prototype.showdata = function (data) {  //显示列表数据

    var fileClass = "file_";
    switch (data.filetype) {
        case "dir":
            fileClass += "folder";
            break;
        case "mp3": case "wav": case "ogg": case "wma":
            fileClass += "music";
            break;
        case "jpg": case "png": case "jpeg": case "bmp":
            fileClass += "img";
            break;
        case "mp4": case "wmv":
            fileClass += "video";
            break;
        case "pdf": case "xls": case "ppt": case "doc":
            fileClass += "doc";
            break;
        default:
            fileClass += "unknown";
            break;
    }
    var temp1 = `<div class="content_list"
                 data-type="` + data.filetype + `" filepath="` + data.filepath + `" filename="` + data.filename + `">
                 <div class="check_box"><input type="checkbox" name="checkitem" class="check_box_input"/></div> 
                 <span class="file_type">`;
    var temp2 = `<i class="file_init ` + fileClass + `"></i></span>`;
    var temp3 = `<span class="file_name">` + data.filename + `</span>`;
    var temp4 = `<span class="file_size">` + tools.SizeToEye(parseInt(data.filesize)) + `</span>`;
    var temp5 = `<span class="file_time">` + data.filetime + `</span>`;
    var temp = temp1 + temp2 + temp3 + temp4 + temp5;
    this.mainNode.innerHTML += temp;  //添加一行数据 
}
listStage.prototype.addListener = function () {  //添加列表事件监听
    this.mainNode.addEventListener("click", this.listClick, false)
}
listStage.prototype.listClick = function (evt) {  //文件或文件夹点击事件
    var target = event.target;
    if (target.className !== "file_name") {  //不是点击的文件名,不做操作
        return
    }
    var parent = target.parentNode;  //父节点
    if (parent.getAttribute("data-type") === "dir") {  //如果点击的是文件夹
        //创建历史路径
        var historyPath = document.getElementById("history_path");
        var lastchild = historyPath.lastElementChild;   //找到最后一个子节点
        var lasttitle = lastchild.getAttribute("title")  //直接点标题
        var i = document.createElement("i");  //创建一个箭头
        i.innerText = ">";
        historyPath.appendChild(i);
        //创建文本节点
        var span = document.createElement("span");
        span.title = parent.getAttribute("filepath") + "/" + parent.getAttribute("filename");
        span.innerText = parent.getAttribute("filename");
        historyPath.appendChild(span)
        //创建点击节点
        var a = document.createElement("a")
        a.title = lasttitle;
        a.innerText = lastchild.innerText;
        historyPath.replaceChild(a, lastchild); //替换原来的文本节点
        //重新获取列表数据
        getList("all", span.title)
    } else {
        //获取当前导航位置,如果是图片视频或音乐，就在线打开
        var nowType = leftNavigation.GetDataType();
        if (nowType === "all"){
            fileHandle.download(parent.getAttribute("filepath"), parent.getAttribute("filename"));//是文件，则下载
        }else if (nowType === "image"){
            imgView.init(); //图片操作初始化
            document.getElementsByClassName("imgview_loding")[0].style.display = "block";
            document.getElementsByClassName("imgview_box")[0].style.display = "block";
            document.getElementById("imgview_action").firstElementChild.innerText = parent.getAttribute("filename");
            var img = document.createElement("img");
            img.src = "/downloadfile?path="+encodeURIComponent( parent.getAttribute("filepath"))+"&filename="+parent.getAttribute("filename");
            img.onload = function(){  //图片加载完成去掉提示操作
                document.getElementsByClassName("imgview_loding")[0].style.display = "none";
            }
            document.getElementById("imgview_view").innerHTML = "";
            document.getElementById("imgview_view").appendChild(img);
        }else if (nowType === "video"){
            console.log("打开视频")
            videoView.init();
            document.getElementsByClassName("videoview_box")[0].style.display = "block";
             document.getElementById("videoview_action").firstElementChild.innerText = parent.getAttribute("filename");
            var video = document.createElement("video");
             video.setAttribute("controls","");
            video.innerText = "抱歉，你的浏览器不支持html5在线播放视频"
            var source = document.createElement("source");
            source.src = "/downloadfile?path="+encodeURIComponent( parent.getAttribute("filepath"))+"&filename="+parent.getAttribute("filename");
            source.type = "video/mp4";
            video.appendChild(source); 
            document.getElementById("videoview_view").innerHTML = "";
            document.getElementById("videoview_view").appendChild(video);

            //播放控制条
            // var div = document.createElement("div");
            // div.className = "video_controller";
            // div.innerHTML = `<span class="video_ctl video_play" id="videoPlay">播放</span>
            // <span class="video_ctl video_time" id="videoTime">00:00/00:00</span>
            // `;
            // document.getElementById("videoview_view").appendChild(div);
           // videoView.initvideo();
        }else if (nowType === "music"){
            console.log("打开音乐")
            musicView.init();
            document.getElementsByClassName("musicview_box")[0].style.display = "block";
            document.getElementById("musicview_action").firstElementChild.innerText = parent.getAttribute("filename");
            var audio = document.createElement("audio");
            audio.innerText = "抱歉，你的浏览器不支持html5在线播放音乐!"
            audio.setAttribute("controls","");
            var source = document.createElement("source");
            source.src = "/downloadfile?path="+encodeURIComponent( parent.getAttribute("filepath"))+"&filename="+parent.getAttribute("filename");
            source.type = "audio/mp3";
            audio.appendChild(source); 
            document.getElementById("musicview_view").innerHTML = "";
            document.getElementById("musicview_view").appendChild(audio);
        }else if (nowType === "doc"){
            fileHandle.download(parent.getAttribute("filepath"), parent.getAttribute("filename"));//是文件，则下载
        }
        
    }
}
//图片预览模块
var imgView = (function(){
    var rotate ;
    var init = function(){  //初始化，主要是监听事件
        rotate = 90;
        document.getElementById("imgview_action").addEventListener("click",imgEventHandle,false)
    }
   
    var imgEventHandle = function(evt){
        var target = evt.target;
        var action = target.getAttribute("action");
        if (action == "close"){
            document.getElementsByClassName("imgview_box")[0].style.display = "none";
        }else if (action == "down"){
            var iframe = document.createElement("iframe")
            iframe.src = document.getElementById("imgview_view").firstElementChild.getAttribute("src");
            iframe.style.display = "none";
            //文件下载方式
            document.body.appendChild(iframe);
        }else if (action == "rotate"){
           
            document.getElementById("imgview_view").firstElementChild.style.transform="rotate("+rotate+"deg)";
            rotate += 90;
        }
    }
    return {
        init:init
    }
})()
//音乐试听模块
var musicView = (function(){
  
    var init = function(){  //初始化，主要是监听事件
        document.getElementById("musicview_action").addEventListener("click",musicEventHandle,false)
    }
    var musicEventHandle = function(evt){
        var target = evt.target;
        var action = target.getAttribute("action");
        if (action == "close"){
            document.getElementById("musicview_view").innerHTML = "";
            document.getElementsByClassName("musicview_box")[0].style.display = "none";
        }else if (action == "down"){
            var iframe = document.createElement("iframe")
            iframe.src = document.getElementById("musicview_view").firstElementChild.firstElementChild.getAttribute("src");
            iframe.style.display = "none";
            //文件下载方式
            document.body.appendChild(iframe);
        }
    }
    return {
        init:init
    }
})()
//视频模块
var videoView = (function(){
  
    var init = function(){  //初始化，主要是监听事件
        document.getElementById("videoview_action").addEventListener("click",videoEventHandle,false)
    }
     var initvideo = function(){
         var video = document.getElementsByTagName("video")[0];
         
         var videolen=0;
         //可播放
         video.addEventListener("canplay",function(){
                
                 videolen = tools.SecondToMin(video.duration)
                document.getElementById("videoTime").innerText = "00:00/"+videolen;
         },false)
         //点击播放或暂停
        document.getElementById("videoPlay").addEventListener("click",function(){
            video.webkitEnterFullscreen()
            if (video.paused){
                video.play();
                this.innerText = "暂停"
            }else {
                video.pause();
                this.innerText = "播放"
            }
        },false)
        //时间改变
        video.addEventListener("timeupdate",function(evt){
            console.log(tools.SecondToMin(evt.target.currentTime),evt.target.currentTime)
            document.getElementById("videoTime").innerText = tools.SecondToMin(evt.target.currentTime)+"/"+videolen;
        },false)
       
    }
    var videoEventHandle = function(evt){
        var target = evt.target;
        var action = target.getAttribute("action");
        if (action == "close"){
            document.getElementById("videoview_view").innerHTML = "";
            document.getElementsByClassName("videoview_box")[0].style.display = "none";
        }else if (action == "down"){
            var iframe = document.createElement("iframe")
            iframe.src = document.getElementById("videoview_view").firstElementChild.firstElementChild.getAttribute("src");
            iframe.style.display = "none";
            //文件下载方式
            document.body.appendChild(iframe);
        }
    }
    return {
        init:init,
        initvideo:initvideo
    }
})()
//历史纪录模块
var historyModule = (function () {
    var node = document.getElementById("history_path");
    var init = function () {
        //路径点击
        node.addEventListener("click", historyHandle, false)
    }
    var historyHandle = function (evt) {
        var target = evt.target;
        if (target.nodeName == "A" || target.nodeName == "a") {
            var x = target.nextSibling;
            //删除之后的所有节点
            while (x) {
                var b = x;
                x = x.nextSibling;
                b.parentNode.removeChild(b);
            }
        } else {
            return
        }
        //创建不可点击节点
        var span = document.createElement("span")
        span.title = target.getAttribute("title");
        span.innerText = target.innerText;
        //把当前结点变成不可点
        target.parentNode.replaceChild(span, target)
        //请求文件
        getList("all", target.getAttribute("title"));
    }
    var setHistory = function (text) {
        node.innerHTML = '<span title="">'+text+'</span>'
    }
    return {
        init: init,
        setHistory: setHistory
    }
})()
window.onload = function () {
    userHandle.GetUserInfo();   //用户信息获取
    addListener();  //时间监听
    tools.ResizeList();  //设置列表显示区高度
    toolBar.init();  //工具栏初始化
    leftNavigation.init();  //左侧导航栏初始化
    userPopBox.init(); //用户弹出框初始化
    historyModule.init();  //历史纪录初始化
    getList();  //获取数据

}

//获取数据列表
function getList(filtype, filepath) {
    var filetype = filtype || "all"; //获取文件类型
    var filepath = encodeURIComponent(filepath || ""); //获取的路径
    var gettime = (new Date()).getTime() / 1000;  //时间戳
    var mystage = new listStage(); //new一个列表操作对象
    mystage.loading();  //显示加载中
    http.httpGet("/nuobel/getfilelist?filetype=" + filetype + "&filepath=" + filepath + "&time=" + gettime, function (resp) {
        mystage.clearList()
        if (resp === 1) {
            tools.Tips("网络错误")
            mystage.showtips("网络出现问题啦,请重试一下~")
            return
        }
        var jsonText = JSON.parse(resp);
        if (jsonText.code != 0) {
            mystage.showtips("服务器开小差去了~")
            return
        }
        var i = 0, j = jsonText.data.length;
        if (j <= 0) {
            mystage.showtips("还没有数据哦~")
            return
        }
        for (i; i < j; i++) {  //渲染数据
            mystage.showdata({
                filename: jsonText.data[i].name,
                filepath: jsonText.data[i].path,
                filetype: jsonText.data[i].type,
                filesize: jsonText.data[i].size,
                filetime: jsonText.data[i].time
            })
        }
        mystage.addListener(); //添加文件名点击事件
    })
}


//部分事件监听
function addListener() {
    document.getElementById("upload_input").addEventListener("change", uploadFileHandle, false)
    window.onresize = tools.ResizeList;

}


//文件上传处理
function uploadFileHandle() {
    //
    uploadBox.init();  //上传框初始化
    var fileinput = document.getElementById("upload_input");
    var detailRow = document.getElementById("detail_row");
    var filepath = tools.GetFilepath();
    //打开上传状态容器
    document.getElementsByClassName("upload_box")[0].style.display = "block";
    //多个文件上传
    for (var i = 0; i < fileinput.files.length; i++) {
        //创建状态栏
        var div = document.createElement("div");
        div.className = "detail_list";
        console.log("fileinfo", fileinput.files[i])
        div.id = "detail_row_" + detailRow.childNodes.length + 1;
        div.innerHTML = `<div class="detail_loading"></div>
                    <span class="detail_name">`+ fileinput.files[i].name + `</span>
                    <span class="detail_status">0%</span>
                    <span class="detail_size">`+ tools.SizeToEye(fileinput.files[i].size) + `</span>
                    <span class="detail_dest">`+ filepath + `</span>`
        detailRow.appendChild(div);
        var file = fileinput.files[i];
        var fd = new FormData();
        fd.append("fielnames", file);
        var obj = {
            url: "uploadfile?path=" + filepath,
            data: fd,
            id: div.id,
            callback: function (resp, id) {
                console.log("resp=", resp)
                console.log("id=", id)
                var child = document.getElementById(id).childNodes;
                for (var i = 0; i < child.length; i++) {
                    if (child[i].className == "detail_status") {
                        if (resp == "success") {
                            child[i].innerText = "已完成";
                            fileHandle.refresh() //完成刷新
                        } else {
                            child[i].innerText = "失败了";
                        }
                    }
                }
            },
            progress: function (progress, id) {
                console.log("进度:", progress)
                console.log("divid=", id)
                var child = document.getElementById(id).childNodes;
                for (var i = 0; i < child.length; i++) {
                    if (child[i].className == "detail_loading") {
                        child[i].style.width = progress + "%";
                    }
                    if (child[i].className == "detail_status") {
                        child[i].innerText = progress + "%";
                    }
                }
            }
        }
        http.httpPostFile(obj)
    }
}


