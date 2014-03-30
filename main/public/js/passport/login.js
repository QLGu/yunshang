$(function () {
    var $form = $('#loginForm');
    $("a.log-btn2").click(function () {
        $form.submit();
    });

    $("form input[name=validateCode]").on('keyup', function (e) {
        if (event.which == 13 || event.keyCode == 13) {
            e.preventDefault();
            $form.submit();
        }
    });

    if(typeof(_errorFields) == 'undefined'){
        $("input[name=login]").focus();
    }
});