//公共函数


//----选择搜索的项目--------
$(function () {
    //网页打开时根据location.hash打开对应的选择    
    //$("#btnradio" + location.hash.replace(/#/, "")).prop("checked", true);
    var jvl = "jingbu";
    if (window.location.href.indexOf("/lvbu") != -1)
        jvl = "lvbu"
    if (window.location.href.indexOf("/lunbu") != -1)
        jvl = "lunbu"
    if (window.location.href.indexOf("/ming") != -1)
        jvl = "ming"
    if (window.location.href.indexOf("/foxuecidian") != -1)
        jvl = "foxuecidian"
    $("#btnradio" + jvl).prop("checked", true);
    $("#param").attr("data-jlv", jvl);
      
    $('#slist .btn').on('click', function () {
        var num = $(this).attr("for").replace(/btnradio/, "");
        $("#param").attr("data-jlv", num);
        $("#param").attr("placeholder", $(this).attr("title"));
    });
});

$(function () {
    //-------绑定下拉框数据-----------
    
    $("#param").on('input', function () {      
        var myReg = /^[\u4e00-\u9fa5]+$/; //过滤非中文
        if (myReg.test($("#param").val())) {            
            var name = $("#param").val();
            var fj = (window.location.href.indexOf("l=") != -1) ? "&1=" : "";
            $("#datas").children().filter('option').remove();
            //var jlv = location.hash.replace(/#/, "");
            //if (jlv == "") jlv = "jingbu";            
            var jlv = $("#param").attr("data-jlv");
            var dir = $("#param").attr("data-dir");
            var url = "/cidian/?q=" + name + fj + "&a=" + jlv + "&dir=" + dir;
            
            $.get(url,
                function (data) {
                    //console.log(data);
                    $("#datas").append(data);
                });
        }
    }).focus(function () {
        $("#datas").children().filter('option').remove();
    });
});

$(function () {
    //--搜索按钮
    //------搜索按钮--------
    $("#bnsearch").click(function () {        
        var val = $("#param").val().trim();        
        if (val.length == 0) return;
        var jlv = $("#param").attr("data-jlv");
        if (jlv == "") jlv = "jingbu";
        var dir = $("#param").attr("data-dir");
        dir = (dir != "") ? "dir=" + dir : "";
        var l = (window.location.href.indexOf("l=") != -1) ? "?l=0" : "";
        if (l=="" && dir!="") dir="?"+dir
        if (l!="" && dir!="") dir="&"+dir
        window.location.href = "/" + jlv + "/" + val + l + dir;
    });
 

    $('#param').keydown(function (e) {
        if (e.keyCode == 13) {
            $('#bnsearch').click();
        }
    });
    $("#retop").click(function () {        
        scrollTo(0, 0);       
    });

    //----繁简体-------
    if ($.cookie('ft') == "1")
        $("#fjti").text("繁體");

    $('#jianti').click(function () {
        $.cookie('ft', "0", { expires: 365, path: '/' });        
        window.location.href = "/";
        //location.reload();
    });
    $('#fanti').click(function () {
        $.cookie('ft', "1", { expires: 365, path: '/' });        
        window.location.href = "/?l=0";
        //location.reload();
    });
});




$(function () {
    //------统计--------
    function tJi() {
        
        var url = window.location.href;
        url = url.replace("//", "");
        var lxing = getmindstr(url, "/", "/", false, false);
        if (lxing!="xianyan") return        
        var tid = getmindstr(url + ";", "/", ";", true, false);

        var td = getmindstr(tid, "?q=", "&", false, false);
        if (td != "")
            tid = td;
        else
            tid = tid.replace("?q=", "");
        var purl = document.referrer;

        var khd;
        khd = (("ontouchstart" in window) == false) ? 0 : 1; //--PC、手机端
        $.get("/llan/?l=" + lxing + "&i=" + tid + "&p=" + purl + "&k=" + khd);
    }
    $(window).load(function () {

        tJi();  
    });

    //---cookie----
    rdcookie = function (cname, cont) {        
        var fsc = $.cookie(cname);
        var ofcs = fsc;
        $.cookie(cname, cont + "," + fsc, {
            expires: 365,
            path: '/'
        });
        if (ofcs == $.cookie(cname)) //--cooke已经到达最大长度。
        {
            ofcs = getmindstr("." + ofcs, ".", "|", false, true);
            ofcs = getmindstr("." + ofcs, ".", "|", false, true); //--删除最后两个
            $.cookie(cname, cont + "," + ofcs, {
                expires: 365,
                path: '/'
            });
        }
    }

 

    $('#clearcooke').click(function () {
        $.removeCookie('fangwenlishi');
        $("#wode").hide();
        //alert("完成");
    });

   

    $(window).unload(function () {        
        rdcookie('fangwenlishi', window.location.href + '|' + document.title);
    });


});

//-----字符串截取----------
function getmindstr(con, l, r, ll, rl) //--获取字符中间的字符,ll,rl是否最后一个匹配
{
    var lp, rp, cp;
    lp = (l == "") ? 0 : ((ll == false) ? lp = con.indexOf(l) : lp = con.lastIndexOf(l));
    if (lp == -1) return "";
    lp = lp + l.length;
    rp = (r == "") ? con.length : ((rl == false) ? con.indexOf(r, lp) : rp = con.lastIndexOf(r));
    if (rp == -1) return "";
    cp = rp - lp;
    if (cp < 0) return "";
    return con.substr(lp, cp);
}
//-----打开对应的tab页面--------
function tabshow(tid) {
    var someTabTriggerEl = document.querySelector(tid)
    var tab = new bootstrap.Tab(someTabTriggerEl)
    tab.show()
};

//------播放文字转语音---------
function playv(text) {
    //console.log(text.replace(/\s+/g, "，")); //-替换英文或者中文空格
    var utterThis = new window.SpeechSynthesisUtterance(text.replace(/\s+/g, "，"));
    window.speechSynthesis.speak(utterThis);    
}

 function getcooknames() {
    //debugger;
    var fsc = $.cookie("fangwenlishi");
    if (fsc == null) return "";
    var arr = fsc.split(',');
    var fs = "", url, title;
    for (var i = 0; i < arr.length; i++) {
        //console.log(arr[i]);
        url = getmindstr('.' + arr[i], '.', "|", false, false);
        title = getmindstr(arr[i], '|', "-", false, false);
        console.log(url);
        console.log(title);
        if (title == "") continue;
        if (url == "http://soufoshuo.com/") continue;
        if (url == "http://soufoshuo.com/?l=0") continue;
        if (i > 21) break;
        fs = fs + "<a href='" + url + "'>" + decodeURI(title) + "</a> <br />";
    }
    return fs;
};