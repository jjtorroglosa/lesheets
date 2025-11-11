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
    localStorage.setItem("code", val);
}
const renderNasheet = (abc) => {
    try {
        localStorage.setItem("code", JSON.stringify(abc));
    } catch {
        localStorage.setItem("code", "");
    }
    const html = go_nasheetToJson(abc);
    const body = document.getElementById("root")
    body.innerHTML = html;
    renderAbcScripts();
}

const initWasm = async () => {
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

    render();
    go.importObject.env = {
        'add': function(x, y) {
            return x + y
        }
        // ... other functions
    }
    return Promise.resolve();
}


const initTextarea = () => {
    const textarea = document.getElementById("editor-text");
    textarea.addEventListener("input", debounce(handleTextChange, 200));
}

let editor;
const vimmode = () => localStorage.getItem("vimmode") == "true";
const toggleVimMode = () => {
    if (!vimmode()) {
        editor.setKeyboardHandler("ace/keyboard/vim");
        localStorage.setItem("vimmode", "true");
    } else {
        editor.setKeyboardHandler(null);
        localStorage.setItem("vimmode", "false");
    }
}
const getEditorContents = () => {
    if (typeof ace !== "undefined") {
        return editor.getValue();
    }
    const textarea = document.getElementById("editor-text");
    return textarea.textContent
}

const render = () => {
    if (typeof abc2svg === "undefined") {
        return;
    }
    setTimeout(() => renderNasheet(getEditorContents()));
}

const initAce = () => {
    let prevCode;
    try {
        prevCode = JSON.parse(localStorage.getItem("code"));
    } catch {
        prevCode = "";
    }
    if (prevCode) {
        document.getElementById("editor-text").textContent = prevCode;
    }
    editor = ace.edit("editor-contents");
    editor.setTheme("ace/theme/dracula");
    editor.session.setMode("ace/mode/markdown");
    if (vimmode()) {
        editor.setKeyboardHandler("ace/keyboard/vim");
    } else {
        editor.setKeyboardHandler(null);
    }
    editor.setFontSize("14px");
    // Listen for changes
    editor.session.on('change', debounce(() => render(), 100));
    editor.focus();
    // Store initial content
    const initialContent = editor.getValue();
    document.addEventListener("DOMContentLoaded", () => {
        const toggle = document.getElementById("vim-mode-toggler");
        toggle.onclick = toggleVimMode;
        toggle.checked = vimmode()
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
    const abcInstance = new abc2svg.Abc(user);
    abcInstance.tosvg("", code, 0, code.length);
    // Hack to remove the undefined coming from the title, or something
    const prefix = "undefined";
    if (svg.startsWith(prefix)) {
        svg = svg.slice(prefix.length);
    }
    return svg;
};

const renderAbcScripts = () => {
    const abcScripts = document.querySelectorAll('script[type="text/vnd.abc"]');

    abcScripts.forEach(i => {
        const code = i.textContent.trim();
        const svg = renderAbc(code);
        i.outerHTML = svg;
    });
}

const init = () => {
    initWasm();
    if (typeof ace !== "undefined") {
        initAce();
    } else {
        initTextarea();
    }
}

init();
