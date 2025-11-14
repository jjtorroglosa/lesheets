export const createActions = (editor) => {
    const addText = (str) => {
        const position = editor.getCursorPosition();
        editor.session.insert(position, str);
        editor.focus()
    }
    return {
        section: () => {
            return addText("# ");
        },
        header: () => {
            return addText("---\ntitle: Title\nsubtitle: Subtitle\ntempo: 123bpm\nkey\nL: 1/8\n---\n");
        },
        push: () => {
            return addText("!push!");
        },
        hold: () => {
            return addText("!hold!");
        },
        fermata: () => {
            return addText("!fermata!");
        },
        diamond: () => {
            return addText("!diamond!");
        },
        diamondFermata: () => {
            return addText("!diamond-fermata!");
        },
        repeatStart: () => {
            return addText("||: ");
        },
        repeatEnd: () => {
            return addText(" :||");
        },
        toggleCollapse: () => {
            // if (document.getElementById("editor").classList.contains("collapsed")) {
            //     document.getElementById('editor').style.width = "50%";
            // } else {
            //     document.getElementById('editor').style.width = "auto";
            // }
            document.getElementById("layout").classList.toggle("collapsed");
            document.getElementById("resizer").classList.toggle("collapsed");
            document.getElementById("editor").classList.toggle("collapsed");
        },
    }
};
