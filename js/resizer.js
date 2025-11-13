const resizer = document.getElementById('resizer');
const left = document.getElementById('editor');
const body = document.getElementById('layout');

let isResizing = false;
let startX = 0;
let startWidth = 0;

resizer.addEventListener('mousedown', e => {
    isResizing = true;
    startX = e.clientX;
    startWidth = left.offsetWidth;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
});

document.addEventListener('mousemove', e => {
    if (!isResizing) return;

    const delta = e.clientX - startX;
    const newWidth = startWidth + delta;

    const minWidth = 150;
    const maxWidth = window.innerWidth - 150;
    left.style.width = Math.min(Math.max(newWidth, minWidth), maxWidth) + 'px';
});

document.addEventListener('mouseup', () => {
    isResizing = false;
    document.body.style.cursor = 'default';
    document.body.style.userSelect = 'auto';
});
