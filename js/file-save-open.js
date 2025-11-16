const saveFilePicker = async () => {
    const opts = {
        types: [{
            description: 'Text Files',
            accept: { 'text/plain': ['.txt'] }
        }]
    };

    try {
        const handle = await window.showSaveFilePicker(opts);
        const writable = await handle.createWritable();
        await writable.write(content);
        await writable.close();
    } catch (err) {
        console.error(err);
    }
}

async function saveFile(content) {
    if ('showSaveFilePicker' in window) {
        // Modern Chrome/Edge/Opera
        const opts = {
            types: [{
                description: 'Text Files',
                accept: { 'text/plain': ['.txt'] }
            }]
        };

        try {
            const handle = await window.showSaveFilePicker(opts);
            const writable = await handle.createWritable();
            await writable.write(content);
            await writable.close();
        } catch (err) {
            console.error("Save cancelled or failed", err);
        }
    } else {
        // Fallback for unsupported browsers
        let filename = prompt("Enter filename", "Song.lesheet");
        if (!filename) return;
        const blob = new Blob([content], { type: 'text/plain' });
        const link = document.createElement('a');
        link.href = URL.createObjectURL(blob);
        link.download = filename;
        link.click();
        URL.revokeObjectURL(link.href);
    }
}
const initOpenFile = (editor) => {
    document.getElementById('openBtn').addEventListener('click', function() {
        // Create a hidden file input
        const input = document.createElement('input');
        input.type = 'file';
        input.accept = '.lesheet'; // only allow text files, you can adjust as needed

        input.onchange = async (event) => {
            const file = event.target.files[0];
            if (!file) return;

            // Read the file contents
            const text = await file.text();

            // Load into the textarea
            // Get the div content
            const div = document.getElementById("editor-text");
            if (div) {
                div.textContent = text; // or textContent
            } else {
                editor.setValue(text)
            }
        };

        // Trigger the file picker
        input.click();
    });
}

const initSaveFile = (editor) => {
    document.getElementById('saveBtn').addEventListener('click', async function() {
        // Get the div content
        const div = document.getElementById("editor-text");
        let text = ""
        if (div) {
            text = div.textContent; // or textContent
        } else {
            text = editor && editor.getValue();
        }
        saveFile(text);
    });
}

export const initSaveAndOpenButtons = (editor) => {
    initOpenFile(editor);
    initSaveFile(editor);
}
