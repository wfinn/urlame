# urlame

This tool can reduce a list of URLs in a way which should be useful for pentesting / bug bounty.  
E.g., when searching intersting URLs in the output of tools like `waymore`, this can do some initial filtering.

`urlame` aims to print one URL per feature of the website in addition to blocking known lame URLs.  
This is done by converting a URL into a pattern and matching that against the patterns seen before.

```sh
go install github.com/wfinn/urlame@latest
urlame < many_urls.txt > less_urls.txt
```

## Example

The easiest way to understand what this does is to see it:

//TODO picture

For a more in depth understanding take a look at [the test TODO]() or the source code in general.

## Things urlame considers lame

`urlame` will filter out lame directories like `/docs`, user profile pages and files with lame extentensions like `.png`.  
This tool also ignores query values, so that only if a new parameter appears on a specific endpoint, the URL will be listed.

It further can detect some patters in parts of URLs. In the following examples, only the first occurance would be listed.

| Pattern | Example |
| ------- | ---------- |
| langage codes | /en-US/admin & /de-DE/admin |
| numeric IDs  | /1000/details & /1001/details |
| hashes | /file/e4a25f7b052442a076b02ee9a1818d2e & /file/bed128365216c019988915ed3add75fb |
| UUIDs | /id/123e4567-e89b-12d3-a456-426614174000 & /id/34c764cf-b13b-4d36-ab93-4474f5b91848|
| common post titles | /common-post-title & /another-post-title |

## Equivalences

Some websites have patterns which should be considered, but which do not apply for other targets, meaning `urlame` would filter stuff we do not want to filter out on these other websites.  
In such cases, you must modify the source code yourself to get decent results.

One mechanism explicitly for target specific filtering exists called "Equivalences".  
These are words you can define, which are kind of equivalent, from our view.  
For example, when filtering tesla URLs you could define `model-3`,`model-y` and so on,
so only the first `%carmodel%-release` and `/api/foo/%carmodel%` URLs are printed.

---

If you have ideas for more stuff to filter out, [let me know](https://github.com/wfinn/urlame/issues/new).

Inspired by [uro](https://github.com/s0md3v/uro)
