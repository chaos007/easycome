/**
 * Created by sanmu on 2017-05-16.
 */
//一些通用接口
(function (window, $) {
    //扩展easyui-validatebox=>手机验证
    $.extend($.fn.validatebox.defaults.rules, {
        mobile: {
            validator: function (value, param) {
                var reg = /(^1[3|4|5|7|8][0-9]{9}$)/;
                if (reg.test(value)) {
                    return true;
                } else {
                    return false;
                }
                ;
            },
            message: '请输入正确的手机号码.'
        }
    });

    var _serverName = "/";

    var _key = "X9dsf_)#&334$R(";

    // 读取 cookie
    function getCookie(c_name) {
        if (document.cookie.length > 0) {
            var all_val = "; " + document.cookie;
            c_start = all_val.indexOf("; " + c_name + "=");
            if (c_start != -1) {
                c_start = c_start + ("; " + c_name).length + 1;
                c_end = all_val.indexOf(";", c_start);
                if (c_end == -1)
                    c_end = all_val.length;
                var val = all_val.substring(c_start, c_end);

                return unescape(val);
            }
        }
        return "";
    }

    //设置 cookie
    function setCookies(name, value, time) {
        var cookieString = name + "=" + escape(value) + ";";
        if (time != 0) {
            var Times = new Date();
            Times.setTime(Times.getTime() + time);
            cookieString += "expires=" + Times.toGMTString() + ";"
        }
        document.cookie = cookieString;
    }

    //初始化物业Combobox
    function init_combobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];
            if (params.select == 0) {
                dataSource.push({name: "物业名称", value: "0"})
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });

            }


            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    //初始化员工Coombobox
    function init_employeeCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];
            if (params.select == 0) {
                dataSource.push({
                    name: "选择员工",
                    value: "0",
                });
            }
            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.Name,
                        value: value.UserName,
                    });
                });
            }


            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });
        }
    }

    //初始化物业人员Combobox
    function init_personCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "选择物业人员",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.Name,
                        value: value.UserName,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });
        }
    }

    function init_buildCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "楼号",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_projectCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "项目",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_brandCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "制造厂家",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_brandTypeCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "型号",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_dutyTeamNameCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "责任小组",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_dealMemberNameCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "责任人",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.Name,
                        value: value.MemberUserName,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_teamListCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "小组",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.TeamName,
                        value: value.TeamName,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_maintainProject(res, params) {
        //alert(res);
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "选择维保项目",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.ProjectName,
                        value: value.ID,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
                disabled: params.disabled ? true : false,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
                disabled: params.disabled ? true : false,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_dealMemberNameOfNumberCombobox(res, params) {
        //alert(res);
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "责任人",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.Name,
                        value: value.MemberUserName,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
                disabled: params.disabled ? true : false,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
                disabled: params.disabled ? true : false,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_eleNoCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "电梯号",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueField: 'value',
                textField: 'name',
                data: dataSource,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_maintainTicketStateCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "工单状态",
                    value: "-1",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value.StateName,
                        value: value.StateID,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueFiled: 'value',
                textField: 'name',
                data: dataSource,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                if (params.select == 0) {
                    $('#' + params.id).combobox('select', '-1');
                } else {
                    $('#' + params.id).combobox('select', params.select + '');
                }
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    function init_RoleCombobox(res, params) {
        res = eval("(" + res + ")");
        if (res.code == 10000) {
            var dataSource = [];

            if (params.select == 0) {
                dataSource.push({
                    name: "角色",
                    value: "0",
                });
            }

            if (res.data == null) {

            } else {
                res.data.forEach(function (value, index, arr) {
                    dataSource.push({
                        name: value,
                        value: value,
                    });
                });
            }

            $("#" + params.id).combobox({
                valueFiled: 'value',
                textField: 'name',
                data: dataSource,
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }

            if (params.select != -1) {
                $('#' + params.id).combobox('select', params.select + '');
            }
        } else {
            $("#" + params.id).combobox({
                data: [],
            });

            //给combobox添加onblur事件
            if (params.onBlur != null) {
                $("#" + params.id).next("span").find("input[type='text']").blur(function () {
                    setTimeout(function () {
                        params.onBlur();
                    }, 300);
                });
            }
        }
    }

    var app = {
        serverName: _serverName,
        key: _key,
        //获取物业名字->combobx
        getPropertyName: function (opt) {
            var options = {
                url: 'GetPropertyName',
                data: {},
                callback: init_combobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取员工列表->combobx
        getEmployeeList: function (opt) {
            var options = {
                url: 'FindNormalPersonWithOutTime',
                data: {},
                callback: init_employeeCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取物业人员列表->combobx
        getPersonList: function (opt) {
            var options = {
                url: 'FindPropertyPerson',
                data: {
                    PropertyName: "",
                    SearchName: "",
                },
                callback: init_personCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取楼号列表->combobox
        getBulidList: function (opt) {
            var options = {
                url: 'GetBuildingName',
                data: {
                    PropertyName: "",
                },
                callback: init_buildCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取项目列表->combobox
        getProjectList: function (opt) {
            var options = {
                url: 'GetProjectName',
                data: {
                    PropertyName: "",
                },
                callback: init_projectCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取品牌列表->combobox
        getBrandList: function (opt) {
            var options = {
                url: 'GetManufacturerName',
                data: {},
                callback: init_brandCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取品牌型号->combobox
        getBrandTypeList: function (opt) {
            var options = {
                url: 'GetManufacturerType',
                data: {
                    ManufacturerName: "",
                },
                callback: init_brandTypeCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取责任大组->combobox
        getDutyTeamName: function (opt) {
            var options = {
                url: 'GetDutyTeamName',
                data: {},
                callback: init_dutyTeamNameCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得责任人->combobox
        getDealMemberName: function (opt) {
            var options = {
                url: 'GetDealMemberName',
                data: {},
                callback: init_dealMemberNameCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得小组->combobox
        getTeamList: function (opt) {
            var options = {
                url: 'GetTeamInfo',
                data: {},
                callback: init_teamListCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取特定类别小组
        getTeamInfoWithType: function (opt) {
            var options = {
                url: 'GetTeamInfoWithType',
                data: {
                    Type: '',
                },
                callback: init_teamListCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },

        //获得保养项目
        getMaintainProject: function (opt) {
            var options = {
                url: 'GetMaintainProjectInfo',
                data: {
                    SearchName: "",
                    BigType: "维保",
                    IntervalType: "",
                },
                callback: init_maintainProject
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取电梯责任人
        getDealMemberNameWithElevatorNum: function (opt) {
            var options = {
                url: 'GetDealMemberNameWithElevatorNum',
                data: {
                    ElevatorNum: "",
                    SmallTeamType: "",
                },
                callback: init_dealMemberNameOfNumberCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获取电梯号
        getEleNoWithBuildName: function (opt) {
            var options = {
                url: 'GetElevatorNumByBuilding',
                data: {
                    BuildingName: "",
                },
                callback: init_eleNoCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得维保工单状态
        getMaintainTicketStateList: function (opt) {
            var options = {
                url: 'GetMaintainTicketStateList',
                data: {},
                callback: init_maintainTicketStateCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得检修工单状态
        getCheckAndFixTicketStateList: function (opt) {
            var options = {
                url: 'GetCheckAndFixTicketStateList',
                data: {},
                callback: init_maintainTicketStateCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得年检工单状态
        geYearCheckTicketStateList: function (opt) {
            var options = {
                url: 'GeYearCheckTicketStateList',
                data: {},
                callback: init_maintainTicketStateCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得自检工单状态
        getCheckFilterCheckTicketStateList: function (opt) {
            var options = {
                url: 'GetSelfCheckTicketStateList',
                data: {},
                callback: init_maintainTicketStateCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得大修整改工单状态
        getBigCheckAndFixTicketStateList: function (opt) {
            var options = {
                url: 'GetBigCheckAndFixTicketStateList',
                data: {},
                callback: init_maintainTicketStateCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },
        //获得角色列表
        getRoleList: function (opt) {
            var options = {
                url: 'GetRoleList',
                data: {},
                callback: init_RoleCombobox
            };
            $.extend(options, opt);
            this.execPostRequest(options);
        },

        //post request 同一接口
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

        //消息提示
        //警告
        alert: function (msg) {
            $.messager.alert('警告', msg, 'warning');
        },
        //成功
        success: function () {
            $.messager.show({
                title: '提示',
                msg: '操作成功',
                timeout: 1000,
            });
        },
        //删除
        delete: function (key, id, rowIndex, callback) {
            $.messager.confirm('提示', '确定删除？', function (r) {
                if (r) {
                    callback(key, id, rowIndex);
                }
            });
        },

        //datagrid客户端分页
        pagerFilter: function (data) {
            if (typeof data.length == 'number' && typeof data.splice == 'function') {	// is array
                data = {
                    total: data.length,
                    rows: data
                }
            }
            var dg = $(this);
            var opts = dg.datagrid('options');
            var pager = dg.datagrid('getPager');
            pager.pagination({
                onSelectPage: function (pageNum, pageSize) {
                    opts.pageNumber = pageNum;
                    opts.pageSize = pageSize;
                    pager.pagination('refresh', {
                        pageNumber: pageNum,
                        pageSize: pageSize
                    });
                    dg.datagrid('loadData', data);
                }
            });
            if (!data.originalRows) {
                data.originalRows = (data.rows);
            }
            var start = (opts.pageNumber - 1) * parseInt(opts.pageSize);
            var end = start + parseInt(opts.pageSize);
            data.rows = (data.originalRows.slice(start, end));
            return data;
        },

        //获取tab url key,来区分新增还是编辑info
        getTabUrlKey: function () {
            return $('#tab').tabs('getSelected').panel("options").href;
        },

        //获取window上的url
        getWindowUrlKey: function () {
            return $('#win').window("options").href;
        },

        //获取tab title
        getTabTitle: function () {
            return $('#tab').tabs('getSelected').panel("options").title;
        },

        //保存之后跳转tab
        saveCallback: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success();
                setTimeout(function () {
                    updateTabWithNewUrl(params.href, params.title);
                }, 1000)
            } else {
                app.alert(res.msg)
            }
        },

        saveCallbackWindow: function (res, params) {
            // alert(res);
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success();
                setTimeout(function () {
                    if (params.action == "window") {
                        openSubWindow(params.href);
                    } else if (params.action == "refresh") {
                        if (params.key == "new_id") {
                            params.key = res.data;
                        }
                        closeSubWindow();
                        params.refresh(params.key);
                    } else {
                        closeSubWindow();
                    }
                }, 1000)
            } else {
                app.alert(res.msg)
            }
        },

        //删除之后刷新tab
        delCallback: function (res, params) {
            res = eval("(" + res + ")");
            if (res.code == 10000) {
                app.success();
                setTimeout(function () {
                    // updateTabWithNewUrl(params.href, params.title);
                    $('#' + params.id).datagrid('deleteRow', params.rowIndex);
                }, 1000)
            } else {
                app.alert(res.msg)
            }
        },

        //删除物业人员
        _delPropertyPerson: function (key, id, rowIndex) {
            var options = {
                url: 'DelPropertyPerson',
                data: {
                    UserName: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };

            app.execPostRequest(options);
        },
        //删除小组
        _delTeam: function (key, id, rowIndex) {
            var options = {
                url: 'DelTeam',
                data: {
                    TeamName: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };

            app.execPostRequest(options);
        },
        //删除物业
        _delProperty: function (key, id, rowIndex) {
            var options = {
                url: 'DelProperty',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };

            app.execPostRequest(options);
        },
        //删除楼号
        _delBuild: function (key, id, rowIndex) {
            var options = {
                url: 'DelBuilding',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除项目
        _delProject: function (key, id, rowIndex) {
            var options = {
                url: 'DelProject',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除品牌
        _delBrand: function (key, id, rowIndex) {
            var options = {
                url: 'DelManufacturer',
                data: {
                    Name: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除电梯
        _delLift: function (key, id, rowIndex) {
            var options = {
                url: 'DelElevatorWithKey',
                data: {
                    ElevatorNum: key,
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除保养单项
        _delMaintain: function (key, id, rowIndex) {
            var options = {
                url: 'DelItemMaintainWithID',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除自检单项
        _delChecking: function (key, id,rowIndex) {
            var options = {
                url: 'DelItemSelfCheckWithID',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除维保项目
        _delItem: function (key, id,rowIndex) {
            var options = {
                url: 'DelMaintainProject',
                data: {
                    ID: parseInt(key),
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除员工(离职)
        _delEmployee: function (key, id, rowIndex) {
            var options = {
                url: 'SetPersonQuitState',
                data: {
                    UserName: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        //删除角色
        _delRole: function (key, id, rowIndex) {
            var options = {
                url: 'DelRole',
                data: {
                    RoleName: key + "",
                },
                callback: app.delCallback,
                callback_params: {
                    id: id,
                    rowIndex: rowIndex,
                }
            };
            app.execPostRequest(options);
        },
        deleteByKey: function (obj, id, rowIndex, callback) {
            var key = $(obj).data("key");
            app.delete(key, id, rowIndex, callback);
        },

        //设置datagrid高度
        setDataGridHeight: function () {
            var sw = $(window).height();
            $(".tab_center").height(sw - 150);
            $(window).resize(function () {
                var sw = $(window).height();
                $(".tab_center").height(sw - 150);
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
                url: 'WebLogin',
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
                                    localStorage.setItem('loginUserName', res.data.UserName);
                                    localStorage.setItem('loginName', res.data.Name);
                                    window.location.href = "/admin/index";
                                } else {
                                    localStorage.clear();
                                    app.alert(res.msg);
                                }
                            },
                            callback_params: {}
                        };
                        app.execPostRequest(options);
                    } else {
                        localStorage.clear();
                        app.alert(res.msg);
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
                        window.location.href = "/admin/index/login.html";
                    } else {
                        app.alert(res.msg);
                    }
                },
                callback_params: {}
            };
            app.execPostRequest(options);
        }
    };

    window.app = app;

    // $.getJSON('/config.txt').done(function (data) {
    //     window.app.serverName = data.WebDomain;
    // });

})(window, jQuery);

