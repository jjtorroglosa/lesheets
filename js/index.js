const debounce = (fn, delay) => {
    let timeoutId;
    return (...args) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => fn(...args), delay);
    };
}

const handleTextChange = (ev) => {
    const val = ev.target.value.trim();
    renderNasheet(val);
    renderAbcs();

}
const renderNasheet = (abc) => {
    const html = go_nasheetToJson(abc);
    const body = document.getElementById("root")
    body.innerHTML = html;
    renderAbcs();
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
}


const initTextarea = () => {
    const textarea = document.getElementById("editor-text");
    textarea.addEventListener("input", debounce(handleTextChange, 200));
}
let editor;
let vimmode = false;
const toggleVimMode = () => {
    console.log("toggling", vimmode);
    if (!vimmode) {
        editor.setKeyboardHandler("ace/keyboard/vim");
        vimmode = true;
    } else {
        editor.setKeyboardHandler(null);
        vimmode = false;
    }
}

const initAce = () => {
    editor = ace.edit("editor-contents");
    editor.setTheme("ace/theme/dracula");
    editor.session.setMode("ace/mode/markdown");
    if (vimmode) {
        editor.setKeyboardHandler("ace/keyboard/vim");
    } else {
        editor.setKeyboardHandler(null);
    }
    editor.setFontSize("14px");
    // Listen for changes
    editor.session.on('change', debounce(() => renderNasheet(editor.getValue()), 100));
    editor.focus();
    // Store initial content
    const initialContent = editor.getValue();
    document.addEventListener("DOMContentLoaded", () => {
        const button = document.getElementById("toggleVimMode");
        button.onclick = toggleVimMode;
    });

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

const renderAbc = (code) => {
    let svg = "";
    const user = {
        read_file: function(fn) {
            console.log("read_file", fn);
        }, // read_file()
        errmsg: function(msg, l, c) {	// get the errors
            console.log("errmsg", fn);
        },
        img_out: function(p) {		// image output
            svg += p
        }
    }
    const abc2 = new abc2svg.Abc(user);
    abc2.tosvg("", code, 0, code.length);
    const prefix = "undefined";
    if (svg.startsWith(prefix)) {
        svg = svg.slice(prefix.length);
    }
    return svg;
};

const renderAbcs = () => {
    const abcScripts = document.querySelectorAll('script[type="text/vnd.abc"]');

    abcScripts.forEach(i => {
        const code = i.textContent.trim();
        const svg = renderAbc(code);
        i.outerHTML = svg;
    });
}
initWasm();
initAce();
//initTextarea();
