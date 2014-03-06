$(function () {
    function addToList(id, value, text) {
        return $('#' + id).append('<option value="' + value + '">' + text + '</option>');
    }

    var cityHolder = 0;
    var cityHolder2 = 0;

    cityHolder = $('#holder').ChinaCitySelect({'prov': '#state', 'city': '#city', 'dist': '#district', url: '/public/libs/city/city.json'});
    cityHolder2 = $('#holder').ChinaCitySelect({url: '/public/libs/city/city.json'});

    if (_p) {
        setTimeout(function () {
            $('#state').val(_p).change();
        }, 300);
    }

    if (_c) {
        setTimeout(function () {
            $('#city').val(_c).change();
        }, 500);
    }

    if (_a) {
        setTimeout(function () {
            $('#district').val(_a).change();
        }, 800);
    }

    function doTest() {
        $('#code').html(cityHolder.getCurrValue());
    }

    function doGetLoc() {
        $('#code_area').html(cityHolder.getCurrLoc());
    }

    function doTestLoc() {
        $('#code_loc').html(cityHolder2.doParseLoc('110229'));
    }
});