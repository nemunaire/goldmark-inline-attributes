# goldmark-inline-attributes

[GoldMark](https://github.com/yuin/goldmark/) inline attributes extension.

```markdown
[Attention]{.underline} some text
```

```html
<p><span class="underline">Attention</span> some text</p>
```

```go
var md = goldmark.New(attributes.Enable)
var source = []byte("[Text]{#id .class1}\nother text")
err := md.Convert(source, os.Stdout)
if err != nil {
    log.Fatal(err)
}
```
