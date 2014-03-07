$(function () {
    function addToList(id, value, text) {
        return $('#' + id).append('<option value="' + value + '">' + text + '</option>');
    }

    var cityHolder = 0;
    var cityHolder2 = 0;

    cityHolder = $('#holder').ChinaCitySelect({'prov': '#state', 'city': '#city', 'dist': '#district', url: '/public/libs/city/city.json'});
    cityHolder2 = $('#holder').ChinaCitySelect({url: '/public/libs/city/city.json'});

    setTimeout(function () {
        _p && $('#state').val(_p).change();
        _c && $('#city').val(_c).change();
        _a && $('#district').val(_a).change();
    }, 1500);


    function doTest() {
        $('#code').html(cityHolder.getCurrValue());
    }

    function doGetLoc() {
        $('#code_area').html(cityHolder.getCurrLoc());
    }

    function doTestLoc() {
        $('#code_loc').html(cityHolder2.doParseLoc('110229'));
    }

    $('#updateImageForm').ajaxForm({
        // dataType identifies the expected content type of the server response
        dataType: 'json',

        // success identifies the function to invoke when the server response
        // has been received
        success: function (ret) {
            alert(ret.message);
            if (ret.ok) {
                $img = $('img.mypic');
                $img.attr("src", $img.attr("src") + "?time=" + new Date().getTime());
            }
            $('.MultiFile-remove').click();
        }
    });
    $(".fancybox").fancybox();
});