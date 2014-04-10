$(function () {
    $('#btn-test-mail').click(function () {
        var email = $("input[data-key='site.mail.test_address']").val();
        if (email) {
            if (confirm("确认发送邮件到" + email)) {
                $.post(testMailURL, {"email": email}, function (ret) {
                    alert(ret.message);
                }, "json");
            }
        } else {
            alert("邮件地址为空")
        }
    });
});