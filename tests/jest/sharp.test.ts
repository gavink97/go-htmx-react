import {describe, expect, test} from "@jest/globals";
import {minifyImages, inputDir, outputDir} from "../../build/sharp";

import fs from 'fs';
import path from 'path';
import sharp from 'sharp';

jest.mock('fs');
jest.mock('sharp', () => {
    const mockSharpInstance = {
        jpeg: jest.fn().mockReturnThis(),
        png: jest.fn().mockReturnThis(),
        toFile: jest.fn(),
    };
    return jest.fn(() => mockSharpInstance);
});

const mockReaddirSync = fs.readdirSync as jest.Mock;
const mockExistsSync = fs.existsSync as jest.Mock;
const mockLstatSync = fs.lstatSync as jest.Mock;
const mockMkdirSync = fs.mkdirSync as jest.Mock;
const mockToFile = sharp().toFile as jest.Mock;

describe('minifyImages', () => {
    beforeEach(() => {
        jest.resetAllMocks();
    });

    test('should minify images correctly', async () => {
        mockReaddirSync.mockReturnValue(['image1.jpg', 'image2.png']);
        mockExistsSync.mockImplementation((filePath: string) => !filePath.includes('image1'));
        mockLstatSync.mockImplementation((filePath: string) => ({
            isFile: () => filePath.includes('.jpg') || filePath.includes('.png')
        }));
        mockMkdirSync.mockImplementation(() => {});

        console.log('Mocked readdirSync:', mockReaddirSync());
        console.log('Sharp instance:', sharp());

        await minifyImages();

        expect(sharp).toHaveBeenCalledWith(path.join(inputDir, 'image1.jpg'));
        expect(sharp).toHaveBeenCalledWith(path.join(inputDir, 'image2.png'));

        expect(sharp().jpeg).toHaveBeenCalledWith({ quality: 60 });
        expect(sharp().png).toHaveBeenCalledWith({ quality: 80, compressionLevel: 9 });

        expect(mockToFile).toHaveBeenCalledWith(path.join(outputDir, 'image1.jpg'));
        expect(mockToFile).toHaveBeenCalledWith(path.join(outputDir, 'image2.png'));
    });
});
