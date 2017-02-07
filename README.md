Microb
======

Single page app engine using Go stdlib and Rethinkdb. This package provides:

- An http server for pages and API requests
- A process that listens to some Rethinkdb changefeeds to detect CRUD events in the database: these events will
trigger templates reparsing and client-side routes update at the http server level

What you get:

- Speed: fast server, light javascript: [Page.js](https://github.com/visionmedia/page.js) (7.8 Ko) and 
[Promise.js](https://github.com/stackp/promisejs) (2,3 Ko) for the client-side routing
- Easy horizontal scalability thanks to [Rethinkdb](http://www.rethinkdb.com)
- Decoupling of admin: to edit content use [Microb Manager](https://github.com/synw/microb-manager) or anything that can edit
documents in Rethinkdb

Configuration
-------------

Create a database in Rethinkdb, ex: `localhost`. Create tables `pages` and `commands`. Set a compound index for table
`pages` with this reql query:

   ```javascript
r.db("localhost").table("pages").indexCreate("key", [r.row("uri"), r.row("domain")])
   ```

   ```json
{
	"db_type": "rethinkdb",
	"db_host":"localhost",
	"db_port":"28015",
	"database":"microb",
	"table":"pages",
	"db_user":"admin",
	"db_password":"mypwd",
	"http_port":":8081",
	"domain": "localhost"
}
   ```
   
Usage
-----

Edit the main template: `templates/view.html`

Edit pages in [Microb Manager](https://github.com/synw/microb-manager) or directly in Rethinkdb. 
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
   
The `uri` field has to be unique with `domain` as it is used as a compound index in Rethinkdb.

Run: `./microb`

Go to `localhost:8081`

Todo
----

- [ ] Better error handling
- [ ] Logging
- [ ] Redis cache
- [ ] Monitoring
