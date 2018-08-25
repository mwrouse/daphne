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
        let contents = fs.readFileSync(filePath);
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




module.exports = {
    readEntireFileSync,
    globFiles,
    canFileBeParsed
}