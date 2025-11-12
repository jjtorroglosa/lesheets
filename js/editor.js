import { initWasm } from './wasm.js';

const handleTextChange = (ev) => {
    const val = ev.target.value.trim();
    render(val);
    localStorage.setItem("code", val);
}


const debounce = (fn, delay) => {
    let timeoutId;
    return (...args) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => fn(...args), delay);
    };
}

const createRenderer = (renderSvgFromAbc, toHtml) => {
    const renderAbcScripts = () => {
        const abcScripts = document.querySelectorAll('script[type="text/vnd.abc"]');

        abcScripts.forEach(i => {
            const code = i.textContent.trim();
            const svg = renderSvgFromAbc(code);
            i.outerHTML = svg;
        });
    }
    const render = (nasheet) => {
        console.time("render")
        try {
            localStorage.setItem("code", JSON.stringify(nasheet));
        } catch {
            localStorage.setItem("code", "");
        }
        const html = toHtml(nasheet);
        const body = document.getElementById("root")
        body.innerHTML = html;
        renderAbcScripts();
        console.timeEnd("render")
    }
    return render;
}

const init = async () => {
    const wasm = await initWasm();
    const abc2svg = await import('abc2svg');
    const render = createRenderer(abc2svg.RenderSvgFromAbc, wasm.toHtml);
    const ace = await import("./ace.js");
    const onChange = debounce(() => render(ace.getEditorContents()), 200);
    ace.initAce(onChange);
    // Trigger first render();
    onChange();
}

init();
