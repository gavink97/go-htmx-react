import * as fs from 'fs';
import * as path from 'path';
import * as esbuild from 'esbuild'

const isBuild = process.argv.includes('-build') || process.argv.includes('-b');

const parseInputs = () => {
    let inputDir = './src/assets';
    let outputDir = './assets/components';

    const args = process.argv.slice(2);
    args.forEach((arg, index) => {
        if (arg === '-i' || arg === '--input') {
            inputDir = args[index + 1];
        }
        if (arg === '-o' || arg === '--output') {
            outputDir = args[index + 1];
        }
    });
    return [inputDir, outputDir];
};

function getSubdirectories(dir) {
    return fs.readdirSync(dir).filter(file => {
        return fs.statSync(path.join(dir, file)).isDirectory();
    });
}


const [inputDir, outputDir] = parseInputs();
const subdirectories = getSubdirectories(inputDir);

var entryPoints = {};

subdirectories.forEach(subdir => {
    const entryFile = path.join(inputDir, subdir, 'index.ts');
    if (fs.existsSync(entryFile)) {
        entryPoints[subdir] = entryFile;
    }
});

const settings = {
    entryPoints: entryPoints,
    outdir: outputDir,
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
