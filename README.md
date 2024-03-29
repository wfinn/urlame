# urlame

This tool can reduce a list of URLs in a way which should be useful for pentesting / bug bounty.  
E.g., when searching interesting URLs in the output of tools like `waymore`, this can do some initial filtering.

![image](https://user-images.githubusercontent.com/42862612/215803684-227232ff-97f7-4fea-af7e-86099da87de6.png)

`urlame` aims to print one URL per feature of the website in addition to blocking known lame URLs.  
This is done by converting a URL into a pattern and matching that against the patterns seen before.

## Things urlame considers lame

As a first step, `urlame` will filter out:
- lame directories like `/docs`
- files with lame extentensions like `.png`
- URLs that look like blog posts
- user profile/referral pages like `/user/FooBar`

This tool also ignores query values, so that only if a new parameter appears on a specific endpoint, the URL will be listed.  
This means once `/foo?id=bar` was seen, `/foo?id=baz` will not be printed.  
Certain URL query parameters are ignored completely, so that `/foo` and `/foo?utm_source=twitter` are considered equal.

It further can detect some patterns in parts of URLs which are ignored when comparing URLs.

- language codes
- numeric IDs
- hashes
- UUIDs

This means that `/en-US/upload/item/1` and `/de-DE/upload/item/5` are considered equal, so only the first will be printed.

## Usage

If you don't have Go installed read [this](https://go.dev/doc/install).

```sh
# installation
go install github.com/wfinn/urlame@latest
# basic usage
urlame < many_urls.txt > less_urls.txt
# practical example
waymore example.org | tee all_urls.txt | urlame > filtered_urls.txt
```
---

If you have ideas for more stuff to filter out or find a bug, [let me know](https://github.com/wfinn/urlame/issues/new).

Inspired by [uro](https://github.com/s0md3v/uro)
