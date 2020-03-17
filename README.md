# gobuzz
<p align="center">
  <img src="https://github.com/egonelbre/gophers/raw/master/.thumb/vector/superhero/gotham.png">
</p>

<p align="justify">
<b>GoBuzz</b> is an example of an HTTP server providing REST API functionality using in-memory storage. Project has been inspired by Kat Zien speech about 
hex-domain architecture proposed by <a href="https://github.com/katzien/go-structure-examples">@katezien</a>. Application is currently
in build process. Used: HTTP light-weight router <a href="https://github.com/go-chi/chi">@chi</a>, BDD framework: <a href="https://onsi.github.io/ginkgo/">@Ginkgo</a>.
</p>

<b>Assumptions:</b>
<ol>
<li>Server is listining on localhost with port 8080</li>
<li>Payload has been limited up to 1 MB per POST request</li>
<li>Worker has five seconds tiemout for fetching URL</li>
<li>Worker fetches data in background with provided interval time in seconds.</li>
<li>Request ID must be an int value.</li>
</ol>


<b>Creating new Post Request</b>:

```curl -si 127.0.0.1:8080/api/fetcher -X POST -d '{"url": "https://httpbin.org/range/15","interval":60}'```

<p align="justify">
For testing purposes only https://httpbin.org/range or https://httpbin.org.delay path are accepted. If duration for fetching
url content will be longer than 5s inside response storage response record will be stored as nil value.</p>

In progress:
<ol>
<li>Adding rest of funcionality for handling PUT, DELETE requests.</li>
<li>Listing global history of created requests and specific fetch response storage data.</li>
</ol>
