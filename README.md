# Concurrent Web Crawler

## Task

We'd like you to write a simple web crawler in a programming language of your choice. Feel free to either choose one you're very familiar with or, if you'd like to learn some Go, you can also make this your first Go program! The crawler should be limited to one domain - so when you start with https://monzo.com/, it would crawl all pages within monzo.com, but not follow external links, for example to the Facebook and Twitter accounts. Given a URL, it should print a simple site map, showing the links between pages.

Ideally, write it as you would a production piece of code. Bonus points for tests and making it as fast as possible!

## Run Instructions

```bash
go run main.go --target=https://jigsaw.xyz --concurrency=1 --singleDomain=true
```

## High Level Design

Utilise channels/go routines by spinning up a set of routines up to `concurrency` which are responsible for fetching URLs. The resultant pages are sent to another set of go routines which are allowed to scale as wide as needed to deal with tokenisation of the html and extraction of the links. This design allows us to have a limited set of concurrent connections to the target site, while dealing with the tokenisation (which can be more expensive than the fetch) in a manner that maximises efficiency.

Once the links have been extracted,they are returned to the main thread which is responsible for maintaining state of which links have been seen, and the relationship between pages via a graph data structure. New links are sent to the worker threads so they can be scraped, and the process continues until we no longer have any pages to scrape.

At that point a sitemap is printed, via a BFS traversal of the graph data structure.