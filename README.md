# md-multi-toc

```
md-multi-toc [-update FILE.md] file1.md file2.md
```

If `-update` is passed, will update `FILE.md` in place with a table of contents
between the markers

```
<!-- ts -->
<!-- te -->
```

## Motivation

The goal of this project is to accept multiple files and to create a TOC with links to the headings in their respective files - optionally updating another file between some comment markers.

There are a few table of contents cli applications, but they either write to stdout OR just operate on 1 file.
