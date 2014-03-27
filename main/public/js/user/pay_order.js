$(function () {
    var ractive = new Ractive({
        el: "logs",
        template: "#logs_tpl",
        data: {
            logs: [],
            format: yunshang.mRenderTime
        }
    });
    ractive.on({
        "load": function () {
            $.getJSON(orderLogsDataURL + "?code=" + orderCode, function (ret) {
                ractive.set("logs", ret.data);
            });
        }
    })
    $('#commentForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            ractive.fire("load");
        }
    });

    ractive.fire("load");
});