const sharp = require('sharp');
const fs = require('fs');
const path = require('path');


const inputDir = './public/images';
const outputDir = './bin/images';

if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
};

const minifyImages = async () => {
    const images = fs.readdirSync(inputDir)

    for (const image of images) {
        const inputFilePath = path.join(inputDir, image);
        const outputFilePath = path.join(outputDir, image);

        if (fs.existsSync(outputFilePath)) {
            console.log(`Skipping minified image: ${image}`);
            continue;
        }

        if (fs.lstatSync(inputFilePath).isFile()) {
            const ext = path.extname(image).toLowerCase();

            // look at sharp api to see all options for image minfication
            // https://sharp.pixelplumbing.com/api-utility
            if (ext === '.jpeg' || ext === '.jpg') {
                await sharp(inputFilePath)
                    .jpeg({ quality: 60 })
                    .toFile(outputFilePath);

            } else if (ext === '.png') {
                await sharp(inputFilePath)
                    .png({ quality: 80, compressionLevel: 9})
                    .toFile(outputFilePath);

            } else {
                console.log(`Skipping unsupported file format: ${image}`)
                continue;
            }

            console.log(`Minified: ${image}`);
        }
    }
};

minifyImages()
    .then(() => console.log('All images minified successfully.'))
    .catch(err => console.error('Error processing images:', err))

module.exports = { minifyImages, inputDir, outputDir};
