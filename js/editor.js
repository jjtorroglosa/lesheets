import { createActions } from './actions.js';
import { time, timeEnd } from './logger.js';
import './resizer.js';
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
        time("render")
        try {
            localStorage.setItem("code", JSON.stringify(nasheet));
        } catch {
            localStorage.setItem("code", "");
        }
        time("renderHtml")
        const html = toHtml(nasheet);
        timeEnd("renderHtml")
        const body = document.getElementById("root")
        body.innerHTML = html;
        time("renderAbcScripts")
        renderAbcScripts();
        timeEnd("renderAbcScripts")
        timeEnd("render")
    }
    return debounce(render, 200);
}

const init = async () => {
    const wasm = await initWasm();
    const abc2svg = await import('abc2svg');
    const render = createRenderer(abc2svg.RenderSvgFromAbc, wasm.toHtml);
    const ace = await import("./ace.js");
    const { editor, setOnChange } = ace.initAce();
    setOnChange(() => render(editor.getValue()));
    // Trigger first render();
    render(editor.getValue());
    document.actions = createActions(editor);
}

init();
