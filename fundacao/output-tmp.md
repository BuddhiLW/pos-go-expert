# Output

`struFoo`: _{Foo 256 true}_ `main.StruFoo`
`struFooBar`: _{[Foo Bar] [256 512] [true false] {Foo 256 true}}_ `main.StruFooBar`
`struBarfoo`: _{[Foo Bar] [256 512] [true false] {Foo 256 true}}_ `main.StruBarfoo`

Struct instances of StruFooBar and StruBarfoo will look identical, but aren't.
Composition vs Standard type (implies direct access vs not direct access):
`struFooBar.S` (_Foo_) == `struBarfoo.SFoo.S` (_Foo_)? `true`
