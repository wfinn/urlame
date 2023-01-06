# urlame

removes lame urls from a list

```sh
go install github.com/wfinn/urlame@latest
cat many_urls.txt | urlame > less_urls.txt
```
The core idea is to normalize URLs, ignoring certain parts when comparing them against each other.  
Additionally some easy to identify lame URLs are filtered.

## Intended use

This tool can reduce a list of URLs in a way which should be good for pentesting / bug bounty.  
E.g., when searchign intersting URLs in the output of tools like `waybackurls`, this can do some initial filtering.  
Then you can feed the resulting list into other tools or manually search interesting URLs in the list.

## Things urlame considers lame

`urlame` will filter out lame directories like `/docs`, user profile pages and files with lame extentensions like `.png`.  
`urlame` also ignores query values, so that only if a new parameter appears on a specific endpoint, the URL will be listed.

It further can detect some patters in parts of URLs. In the following examples, only the first occurance would be listed.

| Pattern | Example |
| ------- | ---------- |
| langage codes | /en-US/admin & /de-DE/admin |
| numeric IDs  | /1000/details & /1001/details |
| hashes | /file/e4a25f7b052442a076b02ee9a1818d2e & /file/bed128365216c019988915ed3add75fb |
| UUIDs | /id/123e4567-e89b-12d3-a456-426614174000 & /id/34c764cf-b13b-4d36-ab93-4474f5b91848|
| common post titles | /common-post-title & /another-post-title |

*In a way* this means a URL is considered lame if the same feature of the website has been seen before.

---

If you have ideas for more stuff to filter out, [let me know](https://github.com/wfinn/urlame/issues/new).

Inspired by [uro](https://github.com/s0md3v/uro)
