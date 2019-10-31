/**
 * Created by sanmu on 2017-06-06.
 */
//移动端，一些通用接口
(function (window, $) {
    var app = {
        serverName: "/",
        key: "X9dsf_)#&334$R(",
        pageSize: 10,

        getLoginUserName: function () {
            return localStorage.getItem('mobile_loginUserName') ? localStorage.getItem('mobile_loginUserName') : "";
        },
        getLoginName: function () {
            return localStorage.getItem('mobile_loginName') ? localStorage.getItem('mobile_loginName') : app.getLoginUserName();
        },
        //判断用户是否登录
        isLogin: function () {
            var loginUserName = localStorage.getItem('mobile_loginUserName') ? localStorage.getItem('mobile_loginUserName') : "";
            if (loginUserName.trim() != "") {
                return true;
            }
            return false;
        },

        //登录
        doLogin: function (obj) {
            var options = {
                url: 'Login',
                data: {
                    UserName: obj.userName,
                    Password: obj.password,
                },
                callback: function (res, params) {
                    res = eval("(" + res + ")");
                    if (res.code == 10000) {
                        var options = {
                            url: 'GetLoginInfo',
                            data: {},
                            callback: function (res, params) {
                                res = eval("(" + res + ")");
                                if (res.code == 10000) {
                                    localStorage.setItem('mobile_loginUserName', res.data.UserName);
                                    localStorage.setItem('mobile_loginName', res.data.Name);
                                    window.location.href = "/mobile/index";
                                } else {
                                    localStorage.clear();
                                    layer.open({
                                        content: res.msg
                                        , btn: '我知道了'
                                    });
                                    return;
                                }
                            },
                            callback_params: {}
                        };
                        app.execPostRequest(options);
                    } else {
                        localStorage.clear();
                        layer.open({
                            content: res.msg
                            , btn: '我知道了'
                        });
                        return;
                    }
                },
                callback_params: {}
            };
            app.execPostRequest(options);
        },

        //退出
        doLogout: function () {
            var options = {
                url: 'QuitLogin',
                data: {},
                callback: function (res, params) {
                    res = eval("(" + res + ")");
                    if (res.code == 10000) {
                        localStorage.clear();
                        window.location.href = "/mobile/index/login.html";
                    } else {
                        layer.open({
                            content: res.msg
                            , btn: '我知道了'
                        });
                        return;
                    }
                },
                callback_params: {}
            };
            app.execPostRequest(options);
        },

        //与后台交换的post基层接口
        execPostRequest: function (options) {
            var randomTime = 100000 + Math.round(Math.random() * 100000);
            var url = this.serverName + options.url + "?random_time=" + randomTime;
            var data = JSON.stringify(options.data);
            var syn = md5(data + this.key);
            // alert(url);
            $.post(url, {
                data: data,
                syn: syn,
            }, function (res) {
                options.callback(res, options.callback_params);
            });
        },

        //把Unix时间戳转成日期格式
        formatDateFromUnix: function (time, separator) {
            if (time == null || time == 0 || time == "") {
                return "";
            }

            if (separator == null) {
                separator = "-";
            }
            var time = new Date(time * 1000);
            var Y = time.getFullYear();
            var M = time.getMonth() + 1;
            var D = time.getDate();
            var format = Y;
            if (M < 10) {
                format += separator + "0" + M;
            } else {
                format += separator + M;
            }

            if (D < 10) {
                format += separator + "0" + D;
            } else {
                format += separator + D;
            }
            return format;
        },

        //把Unix时间戳转成日期时间格式
        formatDateTimeFromUnix: function (time, separator) {
            if (time == null || time == 0 || time == "") {
                return "";
            }

            if (separator == null) {
                separator = "-";
            }
            var time = new Date(time * 1000);
            var Y = time.getFullYear();
            var M = time.getMonth() + 1;
            var D = time.getDate();
            var H = time.getHours();
            var I = time.getMinutes();
            var S = time.getSeconds();
            var format = Y;
            if (M < 10) {
                format += separator + "0" + M;
            } else {
                format += separator + M;
            }

            if (D < 10) {
                format += separator + "0" + D;
            } else {
                format += separator + D;
            }

            if (H < 10) {
                format += " 0" + H;
            } else {
                format += " " + H;
            }

            if (I < 10) {
                format += ":0" + I;
            } else {
                format += ":" + I;
            }

            if (S < 10) {
                format += ":0" + S;
            } else {
                format += ":" + S;
            }

            return format;
        },

        /**********************************地图***********************************************/
        initMap: function (res) {
            var searchService, map, markers = [];
            var center = new qq.maps.LatLng(res.Latitude, res.Longitude);
            if (res.IsInit) {
                map = new qq.maps.Map(document.getElementById(res.ID), {
                    center: center,
                    zoom: 11
                });
            } else {
                map = new qq.maps.Map(document.getElementById(res.ID), {
                    center: center,
                    zoom: 13
                });
                var marker = new qq.maps.Marker({
                    map: map,
                    position: center
                });
                marker.setTitle("当前地址");

                markers.push(marker);
            }
            var latlngBounds = new qq.maps.LatLngBounds();
            //调用Poi检索类
            searchService = new qq.maps.SearchService({
                complete: function (results) {
                    var pois = results.detail.pois;
                    if (pois == null) {
                        app.warning("没有找到相关的地址");
                        return
                    }
                    for (var i = 0, l = pois.length; i < l; i++) {
                        var poi = pois[i];
                        latlngBounds.extend(poi.latLng);
                        var marker = new qq.maps.Marker({
                            map: map,
                            position: poi.latLng
                        });

                        marker.setTitle(i + 1);
                        var info = new qq.maps.InfoWindow({
                            map: map
                        });
                        //获取标记的点击事件
                        qq.maps.event.addListener(marker, 'click', function () {
                            info.open();
                            app.info("经纬度已经修改");
                            $("#w_lat").val(marker.position.lat);
                            $("#w_lng").val(marker.position.lng);
                            var thisPos = new qq.maps.LatLng(marker.position.lat, marker.position.lng);
                        });
                        markers.push(marker);
                    }
                    map.fitBounds(latlngBounds);
                }
            });
            var data = {
                Map: map,
                Markers: markers,
                SearchService: searchService,
            };
            return data;
        },

        //清除地图上的marker
        clearOverlays: function (overlays) {
            var overlay;
            while (overlay = overlays.pop()) {
                overlay.setMap(null);
            }
        }
    };

    window.app = app;
    // $.getJSON('/config.txt').done(function(data){
    //     window.app.serverName = data.WebDomain;
    // });

})(window, jQuery);

function openWinFromIndex(url) {
    if (app.isLogin()) {
        window.location.href = url;
    } else {
        window.location.href = "/mobile/index/login.html";
    }
}

//保养工单电话联系
function doMaintainNoticePhone(ticket_id, name, phone) {
    var options = {
        url: 'AddMaintainPhoneCall',
        data: {
            UserName: app.getLoginUserName(),
            TicketID: parseInt(ticket_id),
            PersonWithPhone: name + "/" + phone,
        },
        callback: function (res, params) {
        },
        callback_params: {}
    };

    app.execPostRequest(options);
}

//检修工单电话联系
function doRushNoticePhone(ticket_id, name, phone) {
    var options = {
        url: 'AddCheckAndFixPhoneCall',
        data: {
            UserName: app.getLoginUserName(),
            TicketID: parseInt(ticket_id),
            PersonWithPhone: name + "/" + phone,
        },
        callback: function (res, params) {
        },
        callback_params: {}
    };

    app.execPostRequest(options);
}

//大修整改电话联系
function doBigNoticePhone(ticket_id, name, phone) {
    var options = {
        url: 'AddBigCheckAndFixPhoneCall',
        data: {
            UserName: app.getLoginUserName(),
            TicketID: parseInt(ticket_id),
            PersonWithPhone: name + "/" + phone,
        },
        callback: function (res, params) {
        },
        callback_params: {}
    };

    app.execPostRequest(options);
}

//年检工单电话联系
function doYearNoticePhone(ticket_id, name, phone) {
    var options = {
        url: 'AddYearCheckPhoneCall',
        data: {
            UserName: app.getLoginUserName(),
            TicketID: parseInt(ticket_id),
            PersonWithPhone: name + "/" + phone,
        },
        callback: function (res, params) {
        },
        callback_params: {}
    };

    app.execPostRequest(options);
}

//自检工单电话联系
function doCheckNoticePhone(ticket_id, name, phone) {
    var options = {
        url: 'AddSelfCheckPhoneCall',
        data: {
            UserName: app.getLoginUserName(),
            TicketID: parseInt(ticket_id),
            PersonWithPhone: name + "/" + phone,
        },
        callback: function (res, params) {
        },
        callback_params: {}
    };

    app.execPostRequest(options);
}