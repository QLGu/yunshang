var LoginLogs = (function () {
    return {"init": function () {
        $('#sample_2').dataTable({
            "aLengthMenu": [
                [5, 15, 20, -1],
                [5, 15, 20, "All"] // change per page values here
            ],
            // set the initial value
            "iDisplayLength": 5,
            "sDom": "<'row-fluid'<'span6'l><'span6'f>r>t<'row-fluid'<'span6'i><'span6'p>>",
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
            "aoColumnDefs": [
                {
                    'bSortable': false,
                    'aTargets': [0]
                }
            ]
        });
    }
    }
})();

$(function () {
    LoginLogs.init();
});
