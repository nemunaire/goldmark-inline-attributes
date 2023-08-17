# goldmark-inline-attributes

[GoldMark](https://github.com/yuin/goldmark/) inline attributes extension.

This implements the [`bracketed_spans`](https://pandoc.org/MANUAL.html#extension-bracketed_spans) of pandoc.

```markdown
[This is *some text*]{.class key="val"} outside text
```

```html
<p><span class="class" key="val">This is <em>some text</em></span> outside text</p>
```

```go
var md = goldmark.New(attributes.Enable)
var source = []byte("[Text]{#id .class1}\nother text")
err := md.Convert(source, os.Stdout)
if err != nil {
    log.Fatal(err)
}
```
