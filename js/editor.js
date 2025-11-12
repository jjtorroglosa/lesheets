import { RenderSvgFromAbc } from '../vendorjs/abc2svg-compiled.js';
import '../vendorjs/wasm_exec_go.js'; // import for side effects

import ace from 'ace-builds/src-noconflict/ace.js';

import 'ace-builds/src-noconflict/keybinding-vim.js';
import 'ace-builds/src-noconflict/mode-markdown.js';
import 'ace-builds/src-noconflict/theme-dracula.js';

// optional: import workers
import 'ace-builds/src-noconflict/worker-javascript.js';

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
    if (vimmode()) {
        editor.setKeyboardHandler("ace/keyboard/vim");
    } else {
        editor.setKeyboardHandler(null);
    }
    editor.setFontSize("14px");
    // Listen for changes
    editor.session.on('change', debounce(() => render(), 100));
    editor.focus();
    document.addEventListener("DOMContentLoaded", () => {
        const toggle = document.getElementById("vim-mode-toggler");
        toggle.onclick = toggleVimMode;
        toggle.checked = vimmode()
    });
}

const renderAbcScripts = () => {
    const abcScripts = document.querySelectorAll('script[type="text/vnd.abc"]');

    abcScripts.forEach(i => {
        const code = i.textContent.trim();
        const svg = RenderSvgFromAbc(code);
        i.outerHTML = svg;
    });
}

const init = () => {
    initWasm();
    console.log("ace.edit", ace, ace.edit);
    if (typeof ace !== "undefined") {
        initAce();
    } else {
        initTextarea();
    }
}

init();
