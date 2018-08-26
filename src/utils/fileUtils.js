let fs = require('fs');
let glob = require('glob');
let path = require('path');
let isBinaryFile = require('isbinaryfile');

/**
 * Reads an entire file
 * @param {string} filePath Path to the file
 * @returns {string} The contents of the file
 */
function readEntireFileSync(filePath) {
    try {
        let contents = fs.readFileSync(filePath, 'utf-8');
        return contents;
    }
    catch (e) {
        return null;
    }
}


/**
 *
 * @param {string} root Root path
 * @param {string} globString glob format string
 * @returns things found in the glob
 */
function globFiles(root, globString) {
    let globPath = path.join(root, globString);

    let result = [];

    let found = glob.sync(globPath);
    for (let i = 0; i < found.length; i++) {
        result.push({
            absolute: found[i],
            relative: path.relative(root, found[i]),
            absoluteDirname: path.dirname(found[i]),
            relativeDirname: path.dirname(path.relative(root, found[i])),
            name: path.basename(found[i])
        });
    }

    return result;
}


/**
 * Determins if a file can be parsed
 * @param {string} filePath Path to the file
 * @returns {boolean}
 */
function canFileBeParsed(filePath) {
    try {
        return !isBinaryFile.sync(filePath);
    } catch (e) {
        return false;
    }
}


/**
 * Returns what is between the delimeters at the top of the file
 * @param {object} file File to get metadata header from
 */
function getMetadataHeader(file, config, debug) {
    file.metadata = {};
    if (!file.shouldParse || file.content == null)
        return;

    debug(`Reading metadata header for ${file.info.relative}`);

    let lines = file.content.trim().split('\n');
    if (lines[0].trim() != config.compiler.tags.delimeter) {
        console.warn(`File '${file.info.relative}' doesn't begin with a proper metadata header`);
        return;
    }

    lines.shift(); // Remove first line (that was already read)

    while (lines.length > 0) {
        let line = lines.shift().trim();
        if (line == config.compiler.tags.delimeter)
            break; // Done reading metadata

        let parts = line.split(':');
        if (parts.length < 2)
            throw new Error(`Unable to parse metadata from file '${file.info.relative}', invalid line in header: ${line}`);

        let key = parts.shift().trim();
        let value = parts.join(':').trim();

        file.metadata[key] = value;
        debug(`\t${key} = ${value}`);
    }

    file.content = lines.join('\n');
    //console.log(`${file.info.relative}: ${lines.length}`);
}


/**
 * Copies a file to the destination
 * @param {string} destination Target file
 * @param {string} source source file
 */
function copyFile(destination, source) {
    fs.copyFileSync(source, destination);
}




module.exports = {
    readEntireFileSync,
    globFiles,
    canFileBeParsed,
    getMetadataHeader,
    copyFile
}