function debounce(fn, delay) {
    let timeoutId;
    return (...args) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => fn(...args), delay);
    };
}

const handleTextChange = (ev) => {
    const val = ev.target.value.trim();
    renderAbc(val);

}
const renderAbc = (abc) => {
    console.log("render abc");
    const html = go_nasheetToJson(abc);
    const body = document.getElementById("root")
    body.innerHTML = html;
}

const initWasm = async () => {
    const go = new Go(); // Defined in wasm_exec.js
    const WASM_URL = '/wasm.wasm';

    var wasm;

    var instance;
    if ('instantiateStreaming' in WebAssembly) {
        instance = WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject);
    } else {
        instance = fetch(WASM_URL)
            .then(resp => resp.arrayBuffer())
            .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
    }

    await instance.then(function(obj) {
        wasm = obj.instance;
        go.run(wasm);
    })

    // const textarea = document.getElementById("editor-text");
    // textarea.addEventListener("input", debounce(handleTextChange, 200));

    go.importObject.env = {
        'add': function(x, y) {
            return x + y
        }
        // ... other functions
    }
    console.log("Wasm initialized");
}


const initAce = () => {
    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/dracula");
    editor.session.setMode("ace/mode/markdown");
    editor.setKeyboardHandler("ace/keyboard/vim");
    editor.setFontSize("14px");
    // Listen for changes
    editor.session.on('change', debounce(() => renderAbc(editor.getValue()), 100));
    editor.focus();
    // Store initial content
    const initialContent = editor.getValue();

    // Listen for page unload
    window.addEventListener("beforeunload", function(e) {
        const currentContent = editor.getValue();

        // Check if content has changed
        if (currentContent !== initialContent) {
            // Standard message for modern browsers (custom messages ignored)
            e.preventDefault();
            e.returnValue = ""; // Required for Chrome
        }
    });
}

initWasm();
initAce();
