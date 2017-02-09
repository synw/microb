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

Create a database in Rethinkdb, ex: `localhost`. Create tables `pages` and `commands` and `hits`. Set a compound index 
for table `pages` with this reql query:

   ```javascript
r.db("localhost").table("pages").indexCreate("key", [r.row("uri"), r.row("domain")])
   ```

   ```json
{
	"centrifugo_host":"localhost",
	"centrifugo_port":"8001",
	"centrifugo_secret_key":"mycentrifugosecretkey",
	"db_host":"localhost",
	"db_port":"28015",
	"db_user":"admin",
	"db_password":"pwd",
	"http_host":":8080",
	"domain": "localhost",
	"hits_log": true,
	"hits_monitor":true,
	"hits_channel":"microb_hits"
}
   ```

[Install the Centrifugo websockets server](https://github.com/centrifugal/centrifugo) if you turned monitoring on

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
- [ ] Decent UI for the default page
