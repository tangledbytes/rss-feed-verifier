# RSS Feed Verifier
Yes, I use RSS to consume a lot of content. I use YARR (hosted on my homelab) as my feeds aggregator. I cobbled together this script a while back when I was collecting
good engineering blogs. The program does one thing and only one thing and that is, it consumes an OPML file and verifies if the RSS feeds inside the file are valid or not.

I run this program infrequently to verify if any of the blog's URLs has changed (happens).

## Usage
```
$ rss-feed-verifier ~/Downloads/engineering_blogs.opml
```
