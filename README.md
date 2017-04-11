## glua [![npm badge](https://img.shields.io/npm/v/glua.svg)](https://www.npmjs.com/package/glua)

`glua` is what happens when you compile https://github.com/yuin/gopher-lua, a Lua VM written in Go, to Javascript. It works right now and you can use it for most awesomeness. You don't have to know Go or even click on the link above, just use this library in your favorite JS environment.

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
```

### try it now

Visit https://rawgit.com/fiatjaf/glua/master/try.html and use your console.
