import * as esbuild from 'esbuild'
import * as fs from 'fs';
import * as path from 'path';

const isBuild = process.argv.includes('-build') || process.argv.includes('-b');

function getSubdirectories(dir) {
    return fs.readdirSync(dir).filter(file => {
        return fs.statSync(path.join(dir, file)).isDirectory();
    });
}

const assetsDir = './src/assets';
const subdirectories = getSubdirectories(assetsDir);

const entryPoints = {};
subdirectories.forEach(subdir => {
    const entryFile = path.join(assetsDir, subdir, 'index.ts');
    if (fs.existsSync(entryFile)) {
        entryPoints[subdir] = entryFile;
    }
});

const settings = {
    entryPoints: entryPoints,
    outdir: './assets/components',
    bundle: true,
    minify: true,
    splitting: true,
    format: 'esm',
    target: ['ESNext'],
    plugins: []
}

let ctx = await esbuild.context(settings)

if (isBuild) {
    await esbuild.build(settings)
    console.log('Build Complete');
    process.exit(0);
} else {
    await ctx.watch();
    console.log('Watching...');
}
