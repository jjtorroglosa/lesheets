import 'ace-builds/src-min-noconflict/ace.js';
import 'ace-builds/src-min-noconflict/keybinding-vim.js';
//import 'ace-builds/src-noconflict/theme-dracula.js';

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

export const getEditorContents = () => {
    if (typeof ace !== "undefined") {
        return editor.getValue();
    }
    const textarea = document.getElementById("editor-text");
    return textarea.textContent
}

export const initAce = (onChange) => {
    let prevCode;
    try {
        prevCode = JSON.parse(localStorage.getItem("code"));
    } catch {
        prevCode = "";
    }
    if (prevCode) {
        document.getElementById("editor-text").textContent = prevCode;
    }
    ace.config.set("useWorker", false);
    editor = ace.edit("editor-contents");
    // editor.setTheme("ace/theme/dracula");
    if (vimmode()) {
        editor.setKeyboardHandler("ace/keyboard/vim");
    } else {
        editor.setKeyboardHandler(null);
    }
    editor.setFontSize("14px");
    // Listen for changes
    editor.focus();

    const toggle = document.getElementById("vim-mode-toggler");
    toggle.onclick = toggleVimMode;
    toggle.checked = vimmode()

    editor.session.on('change', onChange);
    return editor;
}
