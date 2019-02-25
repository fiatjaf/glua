## glua [![npm badge](https://img.shields.io/npm/v/glua.svg)](https://www.npmjs.com/package/glua)

`glua` is what happens when you compile https://github.com/J-J-J/goluajit, a Lua VM written in Go (based on https://github.com/yuin/gopher-lua), to Javascript. It works right now and you can use it for most awesomeness. You don't have to know Go or even click on the link above, just use this library in your favorite JS environment.

### example:

```js
const glua = require('glua')

glua.run(`
  print(12, 'lala', true)
`) // will print these values

var result

glua.runWithGlobals({
  diff: function (a, b) {
    return Math.abs(Math.abs(b) - Math.abs(a))
  },
  saveResult: function (value) {
    result = value
  }
}, `
  local a = 23
  local b = 74
  local difference = diff(a, b)
  saveResult(difference)
`)

console.log('the result is: ', result)

glua.runWithModules({
  fooprinter: `
local fooprinter = {}

function fooprinter.print (foo)
  print('foo value is: ', foo)
end

return fooprinter
  `
}, {
  foo: 264857
}, `
local fooprinter = require('fooprinter')
print('printing foo...')
fooprinter.print(foo)
`)
```

### try it now

Visit https://raw.githack.com/fiatjaf/glua/master/try.html and use your console.

## how do I

1. Return multiple values from a JavaScript function?

  Return the special object `{_glua_multi: []}` with an array of the multiple values you want your function to return.

2. Get values out from the Lua environment?

  Call `.runWithGlobals` passing a function that takes the parameters from Lua and saves them to a JavaScript variable. It works.
