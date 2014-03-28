$(function () {

    var load_das = function () {
        $('#das').load(dasForSelect);
    };
    var load_ins = function () {
        $('#ins').load(insForSelect);
    };


    $('#btn-refresh-das').live("click", load_das);

    $('#btn-refresh-ins').live("click", load_ins);

    $('#ckb-invoice').click(function () {
        $this = $(this)
        if ($this.prop("checked")) {
            $('#ins').show()
        } else {
            $('#ins').hide()
        }
    });

    $('#btn-submit').click(function () {
        var daId = $("input:radio[name ='o.DaId']:checked").val();
        if (!daId) {
            alert("请填写收货地址！");
            return
        }

        if ($('#ckb-invoice').prop("checked")) {
            var inId = $("input:radio[name ='o.InvoiceId']:checked").val();
            if (!inId) {
                alert("请填写发票信息！");
                $('#ckb-invoice').focus();
                return
            }
        }

        $("form").submit();
    });

    load_das();
    load_ins();
});
