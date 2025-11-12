import 'wasm_exec'; // import for side effects

export const initWasm = async () => {
    if (typeof Go === "undefined") {
        return;
    }
    const go = new Go(); // Defined in wasm_exec.js
    const WASM_URL = '/wasm.wasm';

    var wasm;
    if ('instantiateStreaming' in WebAssembly) {
        wasm = await WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
    } else {
        wasm = await fetch(WASM_URL)
            .then(resp => resp.arrayBuffer())
            .then(bytes => WebAssembly.instantiate(bytes, go.importObject));
    }

    go.run(wasm.instance);

    go.importObject.env = {
        'add': function(x, y) {
            return x + y
        }
        // ... other functions
    }
    return { toHtml: go_nasheetToJson };
}
