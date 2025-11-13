// build.js
import { context } from 'esbuild';

const isWatch = process.argv.includes('--watch');
const isDev = isWatch || process.argv.includes('--dev');
const isProd = !isDev;

const ctx = await context({
    entryPoints: ['js/editor.js', 'js/sheet.js', 'js/livereload.js'],
    bundle: true,
    minify: isProd,
    splitting: isProd,
    chunkNames: '[name]-[hash]',
    sourcemap: isProd,
    format: 'esm',
    outdir: 'build',
    target: ['es2017'],
    loader: { '.css': 'css' },
    alias: {
        wasm_exec: './vendorjs/wasm_exec_go.js',
        abc2svg: './vendorjs/abc2svg-wrapper.js',
    },
    define: {
        'process.env.NODE_ENV': JSON.stringify(isProd ? 'production' : 'development'),
    },
})

if (isWatch) {
    await ctx.watch();
    console.log('👀 Watching for changes...');
} else {
    await ctx.rebuild();
    await ctx.dispose();
    console.log('✅ Build complete');
}
