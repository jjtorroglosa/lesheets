export const createActions = (editor) => {
    const addText = (str) => {
        const position = editor.getCursorPosition();
        editor.session.insert(position, str);
        editor.focus()
    }
    const menu = document.getElementById("tools-dropdown");
    const btn = document.getElementById("tools-dropdown-open-button");
    const onClickOutside = (e) => {
        if (!menu.classList.contains("hidden") && !menu.contains(e.target) && !btn.contains(e.target)) {
            menu.classList.add("hidden");
            document.removeEventListener("click", onClickOutside)
        }
    }

    return {
        section: () => {
            return addText("# ");
        },
        header: () => {
            return addText("---\ntitle: Title\nsubtitle: Subtitle\ntempo: 123bpm\nkey: C\nL: 1/8\n---\n");
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
        toggleEditorCollapse: () => {
            document.getElementById("editor").classList.add("transition-all");
            document.getElementById("editor").classList.add("duration-200");
            document.getElementById("editor").style.removeProperty("width");

            setTimeout(() => {
                document.getElementById("layout").classList.toggle("editor-collapsed");
                document.getElementById("resizer").classList.toggle("editor-collapsed");
                document.getElementById("editor").classList.toggle("editor-collapsed");
            });

            setTimeout(() => {
                document.getElementById("editor").classList.remove("transition-all");
                document.getElementById("editor").classList.remove("duration-200");
            }, 200);
        },
        toggleToolbarCollapse: () => {
            if (menu.classList.contains("hidden")) {
                menu.classList.remove("hidden");
                document.addEventListener("click", onClickOutside);
            } else {
                document.removeEventListener("click", onClickOutside);
                menu.classList.add("hidden");
            }
        },
    }
};
