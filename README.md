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
Edit pages in [Microb manager](https://github.com/synw/microb-manager) or directly in Rethinkdb. 
The json documents must have at least these fields:

   ```json
{
	"uri":"/page1/",
	"domain":"mysite.com",
	"fields":{
		"title":"Page title",
		"content":"Page content"
		}
}
   ```
   
The `uri` field has to be unique with domain as it is used as a compound index in Rethinkdb.

Edit the client side routes in `routes.js`

Run: `./microb`

Todo
----

- [ ] Better error handling
- [ ] Logging
- [ ] Redis cache
- [ ] Monitoring
