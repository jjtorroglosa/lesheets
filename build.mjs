// build.js
import { build } from 'esbuild'

const isWatch = process.argv.includes('--watch')

build({
    entryPoints: ['js/editor.js', 'js/sheet.js', 'js/livereload.js'],
    bundle: true,
    minify: !isWatch,
    minify: true,
    splitting: true,
    chunkNames: '[name]-[hash]',
    sourcemap: isWatch,
    format: 'esm',
    outdir: 'build',
    target: ['es2017'],
    loader: { '.css': 'css' },
    alias: {
        wasm_exec: './vendorjs/wasm_exec_go.js',
        abc2svg: './vendorjs/abc2svg-compiled.js',
    },
}).then(() => {
    console.log(isWatch ? '👀 Watching for changes...' : '✅ Build complete!')
})

