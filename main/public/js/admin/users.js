var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            if (!jQuery().dataTable) {
                return;
            }

            // begin first table
            var sampleTable = $('#sample_1').dataTable({
                "bProcessing": true,
                "bServerSide": true,
                "sAjaxSource": "/admin/users/data",
                //"sDom":'T<"clear">lfrtip',
                "sDom": "T<'row-fluid'<'span3'l><'span3'r><'span6'f>>t<'row-fluid'<'span6'i><'span6'p>>",
                "oTableTools": {
                    "sSwfPath": "/public/media/swf/copy_csv_xls_pdf.swf",
                    //"sSelectedClass": "highlight",
                    "sRowSelect": "single",
                    "fnRowSelected": function (nodes) {
                        var oTT = TableTools.fnGetInstance('sample_1');
                        var aData = oTT.fnGetSelectedData();
                        var enabled = aData[0].enabled;
                        var loginName = aData[0].login_name;
                        if (enabled) {
                            $("#sample_editable_1_enabled").html('禁用 <i class="fa fa-minus-circle"></i>');
                        } else {
                            $("#sample_editable_1_enabled").html('激活 <i class="fa fa-check-square"></i>');
                        }
                        $(".toggleable").removeClass("disabled");
                        $(".toggleable").attr("disabled", false);
                        if (loginName == "admin") {
                            $("#sample_editable_1_enabled").attr("disabled", true);
                            $("#sample_editable_1_enabled").addClass("disabled");
                            $("#sample_editable_1_reset_password").attr("disabled", true);
                            $("#sample_editable_1_reset_password").addClass("disabled");
                        } else {
                            $("#sample_editable_1_enabled").attr("disabled", false);
                        }
                    },
                    "aButtons": [
                        "copy",
                        "print",
                        {
                            "sExtends": "collection",
                            "sButtonText": "Save",
                            "aButtons": [ "csv", "xls", "pdf" ]
                        }
                    ]
                },
                "aaSorting": [
                    [0, 'desc']
                ],
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "login_name", "bSortable": false},
                    { "mData": "email", "bSortable": false},
                    { "mData": "real_name", "bSortable": false},
                    { "mData": "from", "bSortable": false},
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            var status = (data == "true" ? "可用" : "不可用");
                            if (data) {
                                return '<span class="label label-success">可用</span>';
                            } else {
                                return '<span class="label label-warn">不可用</span>';
                            }
                        }},
                    { "mData": "last_sign_at", "bSortable": true,
                        "mRender": function (data) {
                            if (data == "0001-01-01T00:00:00Z") {
                                return "-"
                            }
                            return data;
                        }},
                ],
                "aLengthMenu": [
                    [2, 10, 20, 30, 50, -1],
                    ["2条", "10条", "20条", "30条", "50条", "全部"] // change per page values here
                ],
                // set the initial value
                "iDisplayLength": 10,

                "sPaginationType": "bootstrap",
                "oLanguage": {
                    "sLengthMenu": "每页显示_MENU_ 记录",
                    "sInfo": "共计 _TOTAL_ 条， 显示_START_ 到 _END_ 条",
                    "sInfoEmpty": "",
                    "sEmptyTable": "查询不到数据",
                    "sSearch": "搜 索:",
                    "oPaginate": {
                        "sFirst": "首页",
                        "sPrevious": "前一页",
                        "sNext": "后一页",
                        "sLast": "末页"
                    },
                    "sInfoEmtpy": "没有数据",
                    "sProcessing": "正在加载数据...",

                },
                "fnServerParams": function (aoData) {
                    aoData.push({ "name": "more_data", "value": "my_value" });
                },
                "fnRowCallback": function (nRow, aData, iDisplayIndex) {
                    //console.log(JSON.stringify(aData));
                    return nRow;
                }
            });

            //sampleTable.fnDraw();

            jQuery('#sample_1 .group-checkable').change(function () {
                var set = jQuery(this).attr("data-set");
                var checked = jQuery(this).is(":checked");
                jQuery(set).each(function () {
                    if (checked) {
                        $(this).attr("checked", true);
                    } else {
                        $(this).attr("checked", false);
                    }
                });
                jQuery.uniform.update(set);
            });

            jQuery('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            jQuery('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            //jQuery('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown


            $("#sample_editable_1_refresh").click(function () {
                sampleTable.fnDraw(true);
            });

            $("#sample_editable_1_enabled").click(function () {
                var oTT = TableTools.fnGetInstance('sample_1');
                var aData = oTT.fnGetSelectedData();
                var url = changeStatusUrl + "?id=" + aData[0].id;
                doAjaxPost(url, function () {
                    sampleTable.fnDraw(true);
                });
            });

            $("#sample_editable_1_reset_password").click(function () {
                var oTT = TableTools.fnGetInstance('sample_1');
                var aData = oTT.fnGetSelectedData();
                var url = resetPasswordUrl + "?id=" + aData[0].id;
                doAjaxPost(url, function () {
                    sampleTable.fnDraw(true);
                });
            });

            $('#sample_editable_1_loginlog1').click(function(){
                var oTT = TableTools.fnGetInstance('sample_1');
                var aData = oTT.fnGetSelectedData();
                $.fancybox.open({
                    href : showUserLoginLogs +  "?id=" + aData[0].id,
                    type : 'iframe',
                    padding : 5
                });
            });
        }
    };
}();

$(function () {
    TableUsers.init();
})