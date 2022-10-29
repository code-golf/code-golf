var Elm = require('./elm').Elm;
var main = Elm.Main.init();
main.ports.send.subscribe(console.log);
