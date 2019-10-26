let socketUrl = "ws://127.0.0.1:7000/ws";
let apiUrl = "http://127.0.0.1:7070";
let websocket = new WebSocket(socketUrl);
$(document).ready(function () {
    let auth = getLocalStorage("authToken");
    let jsonData = {"authToken": auth};
    $.ajax({
        type: "POST",
        contentType: "application/json",
        dataType: "json",
        url: apiUrl + "/user/checkAuth",
        data: JSON.stringify(jsonData),
        success: function (result) {
            if (result.code == 0) {
                document.getElementById("chatroom-verified").style.display = "block";
                document.getElementById("chatroom-anonymous").style.display = "none";
                $("#nickName").text(result.data.UserName);
                document.getElementById("nickName").style.display = "inline";
                document.getElementById("chatroom-login").style.display = "none";
            } else {
                document.getElementById("chatroom-verified").style.display = "none";
                document.getElementById("chatroom-anonymous").style.display = "block";
                document.getElementById("nickName").style.display = "none";
                document.getElementById("chatroom-login").style.display = "inline";
                document.getElementById("chatroom-logout").style.display = "none";
            }
        },
        error: function () {
            alert("异常！");
        }
    });

    let msg = document.getElementById("msg");
    let data = {"authToken": getLocalStorage("authToken"), "roomId": 1};
    //websocket onopen
    websocket.onopen = function (evt) {
        websocket.send(JSON.stringify(data));
        getRoomInfo();
    };

    websocket.onmessage = function (evt) {
        let data = JSON.parse(evt.data);
        if (data.op == 3) {
            let innerInfo = '<div class="item" ><p class="nick guest j-nick " data-role="guest" data-account="">' + data.fname + '</p><p class="text">' + data.msg + '</p></div>';
            msg.innerHTML += innerInfo + '<br>';
        } else if (data.op == 4) {
            $("#roomOnlineMemberNum").text(data.count);
            document.getElementById('member_info').innerHTML = "";
            let member = document.getElementById("member_info");
            let innerInfo = "";
            for (let k in data.roomUserInfo) {
                innerInfo = innerInfo + '<div class="item" data-id="' + k + '"><div class="avatar"><img src="/static/9.jpeg"> </div> <div class="nick">' + data.roomUserInfo[k] + '</div> </div>';
            }
            member.innerHTML += innerInfo;
        } else if (data.op == 5) {

        }
    };
});

function getRoomInfo() {
    let jsonData = {roomId: 1};
    $.ajax({
        type: "POST",
        dataType: "json",
        url: apiUrl + "/push/getRoomInfo",
        data: JSON.stringify(jsonData),
        success: function (result) {
            if (result.code != 0) {
                alert("请求出错，请稍后重试！");
            }
        },
        error: function () {
            alert("异常！");
        }
    });
}


function send() {
    let msg = document.getElementById('editText').value;
    document.getElementById('editText').value = '';
    let jsonData = {op: 3, msg: msg, roomId: 1, authToken: getLocalStorage("authToken")};
    $.ajax({
        type: "POST",
        dataType: "json",
        url: apiUrl + "/push/pushRoom",
        data: JSON.stringify(jsonData),
        success: function (result) {
            if (result.code == 0) {

            } else {
                alert("请先简单注册登录");
                window.location.href = "/register.html";
            }
        },
        error: function () {
            alert("异常！");
        }
    });
}

function getLocalStorage(name) {
    return localStorage.getItem(name);
}

function logout() {
    let jsonData = {authToken: getLocalStorage("authToken")};
    $.ajax({
        type: "GET",
        dataType: "json",
        url: apiUrl + "/user/logout",
        data: JSON.stringify(jsonData),
        success: function (result) {
            if (result.code == 0) {
                window.location.href = "/login.html";
            } else {
                alert("请求出错，请稍后重试！");
            }
        },
        error: function () {
            alert("异常！");
        }
    });
}

function changeTab(type) {
    if (type == "chat") {
        document.getElementById("tab_chat").className = "crt j-tab";
        document.getElementById("tab_member").className = "j-tab";
        document.getElementById("msg").className = "chat j-pannel j-chat";
        document.getElementById("member_list").className = "member j-pannel hide";
    } else {
        document.getElementById("tab_chat").className = "j-tab";
        document.getElementById("tab_member").className = "crt j-tab";
        document.getElementById("member_list").className = "member j-pannel";
        document.getElementById("msg").className = "chat j-pannel j-chat hide";
    }
}