/**
 * Created by sanmu on 2017-05-16.
 */
//一些通用接口
(function (window, $) {
    var app = {
        serverName: "/",
        key: "X9dsf_)#&334$R(",
        module: '/admin/',

        //post request 统一接口
        execPostRequest: function (options) {
            var url = this.serverName + options.url;
            var data = JSON.stringify(options.data);
            var syn = md5(data + this.key);
            $.post(url, {
                data: data,
                syn: syn,
            }, function (res) {
                options.callback(res, options.callback_params);
            });
        },
        //get request 获取文件统一接口
        execGetRequest: function (options) {
            var url = this.serverName + options.url;
            var data = "";
            var syn = "";
            if (options.data) {
                data = JSON.stringify(options.data);
                syn = md5(data + this.key);
            }
            var filename = options.filename;
            var method = 'get';
            if (options.method) {
                var method = options.method;
            }

            var xmlHttpReg = null;
            if (window.ActiveXObject) {//如果是IE
                xmlHttpReg = new ActiveXObject("Microsoft.XMLHTTP");

            } else if (window.XMLHttpRequest) {
                xmlHttpReg = new XMLHttpRequest(); //实例化一个xmlHttpReg
            }

            //如果实例化成功,就调用open()方法,就开始准备向服务器发送请求
            if (xmlHttpReg != null) {
                var formdata = new FormData();
                if (data != "") {
                    url += "?data=" + data;
                    url += "&syn=" + syn;
                }
                xmlHttpReg.open(method, url, true);
                xmlHttpReg.responseType = "blob";
                xmlHttpReg.onreadystatechange = function () {
                    if (xmlHttpReg.readyState === 4 && xmlHttpReg.status === 200) {
                        if (xmlHttpReg.response.size > 0) {
                            if (typeof window.chrome !== 'undefined') {
                                // Chrome version
                                var link = document.createElement('a');
                                link.href = window.URL.createObjectURL(xmlHttpReg.response);
                                link.download = filename;
                                link.click();
                            } else if (typeof window.navigator.msSaveBlob !== 'undefined') {
                                // IE version
                                var blob = new Blob([xmlHttpReg.response], { type: 'application/force-download' });
                                window.navigator.msSaveBlob(blob, filename);
                            } else {
                                // Firefox version
                                var file = new File([xmlHttpReg.response], filename, { type: 'application/force-download' });
                                window.open(URL.createObjectURL(file));
                            }
                        } else {
                            app.error("导出失败");
                        }
                    }
                };
                xmlHttpReg.send(formdata);
            }

        },
        /******************************************弹框Begin*****************************************/
        info: function (msg) {
            toastr.options = {
                positionClass: 'toast-bottom-center',
                timeOut: 2000,
            }
            toastr.clear();
            toastr.info(msg, "提示");
        },
        warning: function (msg) {
            toastr.options = {
                positionClass: 'toast-bottom-center',
                timeOut: 3000,
            }
            toastr.clear();
            toastr.warning(msg, "警告");
        },
        error: function (msg) {
            toastr.options = {
                positionClass: 'toast-bottom-center',
                timeOut: 3000,
            }
            toastr.clear();
            toastr.error(msg, "错误");
        },
        success: function (msg, callback) {
            toastr.options = {
                positionClass: 'toast-bottom-center',
                timeOut: 1000,
            }
            toastr.clear();
            toastr.success(msg, "成功");

            setTimeout(function () {
                callback();
            }, 1000);
        },
        /******************************************弹框End*******************************************/


        /*******************jqGrid*******************************************************************/
        rowList: [10, 20, 30, 50, 100],
        fitGrid: function () {
            var width = $('.jqGrid_wrapper').width();
            $('#table').setGridWidth(width);

            // var form_height = $("#search-form").outerHeight();
            // var height = $(window).height() - 290 - form_height;
            // if (height < 250) {
            //     $("#table").setGridHeight(250);
            // } else {
            //     $("#table").setGridHeight(height);
            // }
        },
        /*******************jqGrid*******************************************************************/

        /*****************************win************************************************************/
        setWinUrl: function (obj) {
            app.win_url = $(obj).attr("href");
        },
        //设置每月工作概览里工单详情的查询参数
        setWinOptions: function (options) {
            app.win_options = options;
        },
        openWinFromWin: function (obj) {
            app.win_url = $(obj).data("href");
            $("#win").modal("hide");
            $("#win").modal({
                remote: app.win_url
            });
        },
        win_url: "",
        win_action: "add",
        win_key: "",
        win_options: {},
        /*****************************win************************************************************/

        saveCallbackWindow: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success("操作成功", function () {
                    $('#win').modal('hide');
                    if (params.action == "window") {
                        app.win_url = params.href;
                        $("#win").modal({
                            remote: params.href
                        });
                    } else {
                        if (params.key == "new_key") {
                            params.key = res.data;
                        }
                        params.refresh(params.key, params.action);
                    }
                });
            } else {
                app.error(res.msg)
            }
        },

        saveCallbackElevatorWindow: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success("操作成功", function () {
                    $('#win').modal('hide');
                    if (params.action == "window") {
                        app.win_url = params.href;
                        $("#win").modal({
                            remote: params.href
                        });
                    } else {
                        if (params.key == "new_key") {
                            params.key = res.data;
                        }
                        params.refresh(params.key, params.oldKey, params.action);
                    }
                });
            } else {
                app.error(res.msg)
            }
        },

        saveCallbackChangeStockWindow: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success("操作成功", function () {
                    $('#win').modal('hide');
                    if (params.action == "window") {
                        app.win_url = params.href;
                        $("#win").modal({
                            remote: params.href
                        });
                    } else {
                        if (params.key == "new_key") {
                            params.key = res.data;
                        }
                        params.refresh(params.key, params.oldKey, params.action);
                    }
                });
            } else {
                app.error(res.msg)
            }
        },

        openRightFrameFromWin: function (url) {
            $('#win').modal('hide');
            window.location.href = url;
        },

        /**********************************数据源控件***********************************************/
        initCombobox: function (opt) {
            var options = {
                url: '',
                data: {},
                callback: function (res, params) {
                    res = eval("(" + res + ")");
                    if (res.code == 10000) {
                        if (res.data) {
                            if (params.url == "FindPropertyPerson" || params.url == "FindNormalPerson") {
                                //物业人员
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.UserName + '">' + value.Name + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "GetDealMemberName" || params.url == "GetDealMemberNameWithElevatorNum") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.MemberUserName + '">' + value.Name + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "GetTeamInfoWithType" || params.url == "GetTeamInfo") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.TeamName + '">' + value.TeamName + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "GetMaintainProjectInfo") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.ID + '">' + value.ProjectName + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "GetMaintainTicketStateList"
                                || params.url == "GetCheckAndFixTicketStateList"
                                || params.url == "GetBigCheckAndFixTicketStateList"
                                || params.url == "GeYearCheckTicketStateList"
                                || params.url == "GetSelfCheckTicketStateList") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.StateID + '">' + value.StateName + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "FindTeamMemberInfoWithKey") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.UserName + '">' + value.PersonName + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else if (params.url == "GetSupportInfo" || params.url == "GetRepertoryInfo" || params.url == "GetPartNumList") {
                                res.data.forEach(function (value, index, arr) {
                                    var html = '<option value="' + value.ID + '">' + value.Name + '</option>';
                                    $("#" + params.id).append(html);
                                });
                            } else {
                                res.data.forEach(function (value, index, arr) {
                                    var html = "";
                                    if (isFinite(value)) {
                                        html = '<option value="' + value + '">' + value.trim() + '</option>';
                                    } else {
                                        html = '<option value="' + value.trim() + '">' + value.trim() + '</option>';
                                    }
                                    $("#" + params.id).append(html);
                                });
                            }
                        }
                    }


                    if (params.value && params.value != "") {
                        if (isFinite(params.value)) {
                            $("#" + params.id).val(params.value);
                        } else {
                            $("#" + params.id).val(params.value.trim());
                        }
                    }

                    if (params.refresh) {
                        $("#" + params.id).combobox("refresh");
                        if (params.value == "" || params.value == null) {
                            $("#" + params.id).combobox("clearTarget");
                            $("#" + params.id).combobox("clearElement");
                        }
                    } else if (params.clear) {
                        $("#" + params.id).combobox("refresh");
                        $("#" + params.id).combobox("clearTarget");
                        $("#" + params.id).combobox("clearElement");
                    } else {
                        $("#" + params.id).combobox();
                    }
                },
                callback_params: {}
            };
            $.extend(options, opt);

            options.callback_params.url = options.url;
            this.execPostRequest(options);
        },

        templateFileDownload: function (opt) {
            var options = {
                url: 'res/templatefile/' + opt.url,
                filename: opt.filename,
            };
            this.execGetRequest(options);
        },

        /**********************************数据源控件***********************************************/


        /**********************************删除***********************************************/
        deleteProperty: function (key) {
            var options = {
                url: 'DelProperty',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    key: key
                }
            };

            app.execPostRequest(options);
        },
        deleteBuilding: function (key) {
            var options = {
                url: 'DelBuilding',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    key: key
                }
            };

            app.execPostRequest(options);
        },
        deleteProject: function (key) {
            var options = {
                url: 'DelProject',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    key: key
                }
            };

            app.execPostRequest(options);
        },
        deleteBrand: function (key) {
            var options = {
                url: 'DelManufacturer',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    key: key
                }
            };

            app.execPostRequest(options);
        },
        deleteMaintain: function (key) {
            var options = {
                url: 'DelItemMaintainWithID',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteElevator: function (key) {
            var options = {
                url: 'DelElevatorWithKey',
                data: {
                    ElevatorNum: key,
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteCheck: function (key) {
            var options = {
                url: 'DelItemSelfCheckWithID',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteItem: function (key) {
            var options = {
                url: 'DelMaintainProject',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        abandonMaintainProject: function (key) {
            var options = {
                url: 'AbandonMaintainProject',
                data: {
                    ID: parseInt(key),
                },
                callback: function (res) {
                    res = eval("(" + res + ")");
                    if (res.code == 10000) {
                        app.success("弃用成功！", function () {
                        });
                    } else {
                        app.error(res.msg)
                    }
                },
                callback_params: {
                }
            };
            app.execPostRequest(options);
        },
        deleteRole: function (key) {
            var options = {
                url: 'DelRole',
                data: {
                    RoleName: key,
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deletePreson: function (key) {
            var options = {
                url: 'SetPersonQuitState',
                data: {
                    UserName: key,
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteGroup: function (key) {
            var options = {
                url: 'DelTeam',
                data: {
                    TeamName: key,
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteRepertory: function (key) {
            var options = {
                url: 'DelRepertory',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },
        deleteSupport: function (key) {
            var options = {
                url: 'DelSupport',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    key: key,
                }
            };
            app.execPostRequest(options);
        },

        deleteByKey: function (key, callback) {
            swal({
                title: "警告",
                text: "确定删除？",
                type: "warning",
                showCancelButton: true,
                cancelButtonColor: 'white',
                cancelButtonText: '取消',
                confirmButtonColor: "#18a689",
                confirmButtonText: "确认",
                closeOnConfirm: true
            }, function () {
                callback(key);
            });
        },
        abandonByKey: function (key, callback) {
            swal({
                title: "警告",
                text: "确定弃用？",
                type: "warning",
                showCancelButton: true,
                cancelButtonColor: 'white',
                cancelButtonText: '取消',
                confirmButtonColor: "#18a689",
                confirmButtonText: "确认",
                closeOnConfirm: true
            }, function () {
                callback(key);
            });
        },
        delete: function (callback) {
            swal({
                title: "警告",
                text: "确定删除？",
                type: "warning",
                showCancelButton: true,
                cancelButtonColor: 'white',
                cancelButtonText: '取消',
                confirmButtonColor: "#18a689",
                confirmButtonText: "确认",
                closeOnConfirm: true
            }, function () {
                callback();
            });
        },

        stopProject: function (key, callback) {
            swal({
                title: "警告",
                text: "确定停保？",
                type: "warning",
                showCancelButton: true,
                cancelButtonColor: 'white',
                cancelButtonText: '取消',
                confirmButtonColor: "#18a689",
                confirmButtonText: "确认",
                closeOnConfirm: true
            }, function () {
                callback(key);
            });
        },

        delCallback: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success("删除成功！", function () {
                    $("#table").delRowData(params.key);
                });
            } else {
                app.error(res.msg)
            }
        },
        /**********************************删除***********************************************/

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
            if (separator == "zh-CN") {
                if (M < 10) {
                    format += '年' + "0" + M;
                } else {
                    format += '年' + M;
                }

                if (D < 10) {
                    format += '月' + "0" + D + '日';
                } else {
                    format += '月' + D + '日';
                }
            } else {
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
            }
            return format;
        },

        //把Unix时间戳转换成日期时间格式
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

        //获得登录用户名
        getLoginUserName: function () {
            return localStorage.getItem('loginUserName') ? localStorage.getItem('loginUserName') : "";
        },
        getLoginName: function () {
            return localStorage.getItem('loginName') ? localStorage.getItem('loginName') : app.getLoginUserName();
        },
        //判断用户是否登录
        isLogin: function () {
            var loginUserName = localStorage.getItem('loginUserName') ? localStorage.getItem('loginUserName') : "";
            if (loginUserName.trim() != "") {
                return true;
            }
            return false;
        },

        //登录
        doLogin: function (obj) {
            var options = {
                url: 'BackendLogin',
                data: {
                    UserName: obj.userName,
                    Password: obj.password,
                },
                callback: function (res, params) {
                    res = eval("(" + res + ")");
                    if (res.code == 10000) {
                        var options = {
                            url: 'UpGetBackendLoginInfo',
                            data: {},
                            callback: function (res, params) {
                                res = eval("(" + res + ")");
                                if (res.code == 10000) {
                                    app.success("登录成功，马上跳转...", function () {
                                        localStorage.setItem('loginUserName', res.data.UserName);
                                        localStorage.setItem('loginName', res.data.Name);
                                        window.location.href = app.module + "index";
                                    })
                                } else {
                                    localStorage.clear();
                                    app.error(res.msg);
                                }
                            },
                            callback_params: {}
                        };
                        app.execPostRequest(options);
                    }
                    else {
                        localStorage.clear();
                        app.error(res.msg);
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
                        window.location.href = app.module + "index/login.html";
                    } else {
                        app.error(res.msg);
                    }
                },
                callback_params: {}
            };
            app.execPostRequest(options);
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
})(window, jQuery);

