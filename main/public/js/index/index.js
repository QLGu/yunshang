/**
 * require(jQuery).
 */
$(function () {
    $('.log-btn').click(function () {
        $login = $('input[name=login]');
        $password = $('input[name=password]');
        $.post(loginURL,
            {login: $login.val(), password: $password.val()},
            function (ret) {
                if (ret.ok) {
                    window.location.reload();
                } else {
                    window.location.href = "/passport/login";
                }
            },
            "json");
    });
});