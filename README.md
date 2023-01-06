# urlame

removes lame urls from a list

```sh
go install github.com/wfinn/urlame@latest
cat many_urls.txt | urlame > less_urls.txt
```

The core idea is to normalize URLs, ignoring certain parts when comparing them against each other.

## Intended Use

This tool can reduce a list of URLs in a way which should be good for pentesting / bug bounty.  
E.g., when viewing and decluttering the output of tools like `waybackurls`, this can do some initial filtering.  

## Things urlame considers lame

`urlame` ignores query values, so that only if a new parameter appears on a specific endpoint, the URL will be listed.

`urlame` will also block certain URLs based on:

| Thing | Example |
| ------- | ---------- |
| file extensions | /xyz.woff2 |
| (currently optionally) path names | /docs |

It further can detect some patters in parts of URLs, in the following examples, only the first occurance would be listed.

| Pattern | Example |
| ------- | ---------- |
| langage codes | /en-US/admin & /de-DE/admin |
| numeric IDs  | /1000/details & /1001/details |
| hashes | /file/e4a25f7b052442a076b02ee9a1818d2e & /file/bed128365216c019988915ed3add75fb |
| UUIDs | /id/123e4567-e89b-12d3-a456-426614174000 & /id/34c764cf-b13b-4d36-ab93-4474f5b91848|
| profile pages | /user/max & /user/moritz |
| common post titles | /common-post-title & /another-post-title |

*In a way* this means a URL is considered lame if the same feature of the website has been seen before.

---

Inspired by [uro](https://github.com/s0md3v/uro)
