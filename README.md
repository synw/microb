Microb
======

Single page app engine using Go stdlib and Rethinkdb

Configuration
-------------

   ```json
{
	"db_type": "rethinkdb",
	"db_host":"localhost",
	"db_port":"28015",
	"database":"microb",
	"table":"pages",
	"db_user":"admin",
	"db_password":"mypwd",
	"http_port":":8080",
	"domain": "localhost"
}
   ```
Edit pages in [Microb manager](https://github.com/synw/microb-manager) or directly in Rethinkdb. The json files must have at least these fields:

   ```json
{
	"uri":"/page1/",
	"fields":{
		"title":"Page title",
		"content":"Page content"
		}
}
   ```

Edit the client side routes in `routes.js`

Run: `./microb`

