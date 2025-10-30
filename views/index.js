const go = new Go(); // Defined in wasm_exec.js
const WASM_URL = 'wasm.wasm';

var wasm;

var instance;
if ('instantiateStreaming' in WebAssembly) {
    instance = WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
} else {
    instance = fetch(WASM_URL)
        .then(resp => resp.arrayBuffer())
        .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
}

instance.then(function(obj) {
    wasm = obj.instance;
    go.run(wasm);
}).then(function(obj) {
    const res = go_nasheetToJson("# section\nA | Dm7(#11)");
    console.log(JSON.parse(res));

});

go.importObject.env = {
    'add': function(x, y) {
        return x + y
    }
    // ... other functions
}
