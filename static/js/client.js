function loadPage(resturl){
	var container = 'content';
	promise.get(resturl,{},{"Accept":"application/json"}).then(function(error, data, xhr) {
	    if (error) {console.log('Error ' + xhr.status);return;}    
	    var parsed_data = JSON.parse(data);
	    var content = parsed_data.content;
	    var title = parsed_data.title;
	    top.document.title = title;
	    top.document.getElementById(container).innerHTML = content;
	});
	return
}
page('/', function(ctx, next) { loadPage('/x/') } );
page('/page1/', function(ctx, next) { loadPage('/page1/') } );
page();