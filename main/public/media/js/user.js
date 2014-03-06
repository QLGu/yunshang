var UserApp = (function(){
    return {
        init: function(){
            $('.layout-option')[0].selectedIndex = 1;
            $('.layout-option').change();
        }
    }
})();

$(function(){
   //UserApp.init();
});