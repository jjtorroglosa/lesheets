import 'ace-builds/src-min-noconflict/ace.js';
import 'ace-builds/src-min-noconflict/keybinding-vim.js';
//import 'ace-builds/src-noconflict/theme-dracula.js';


export const initAce = () => {
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
    const editor = ace.edit("editor-contents");
    // editor.setTheme("ace/theme/dracula");
    const vimmode = () => localStorage.getItem("vimmode") == "true";
    if (vimmode()) {
        editor.setKeyboardHandler("ace/keyboard/vim");
    } else {
        editor.setKeyboardHandler(null);
    }

    const toggleVimMode = () => {
        if (!vimmode()) {
            editor.setKeyboardHandler("ace/keyboard/vim");
            localStorage.setItem("vimmode", "true");
        } else {
            editor.setKeyboardHandler(null);
            localStorage.setItem("vimmode", "false");
        }
    }
    editor.setFontSize("14px");
    // Listen for changes
    editor.focus();

    const toggle = document.getElementById("vim-mode-toggler");
    toggle.onclick = toggleVimMode;
    toggle.checked = vimmode()

    return { editor, setOnChange: (onChange) => editor.session.on('change', onChange) };
}
