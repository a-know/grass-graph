var getQueryVars = function() {
  // return用の配列
  var vars = [];

  // クエリ文字列を取得して「&」で分割
  var query_list = window.location.search.substring(1).split('&');

  // 値取得用のテンポラリ変数
  var tmp_arr;

  // 分割したクエリ文字列の配列から、値を取り出す
  query_list.forEach( function(e, i, a) {
    tmp_arr = e.split('=');
    vars[ tmp_arr[0] ] = tmp_arr[1];
  })

  return vars;
}

$(function(){

  // query string
  var query_vars = getQueryVars();

  // visitor notify
  $.ajax({
        type : 'POST',
        url : '/knock',
        data : { 'user_agent' : navigator.userAgent, 'language' : navigator.language, 'admin' : query_vars['admin'] },
        cache : false,
        dataType : 'json',

        success : function(json) {
          // no operation
        },
        complete : function() {
          // no operation
        }
    });
});