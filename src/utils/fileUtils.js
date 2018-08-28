let fs = require('fs');
let glob = require('glob');
let path = require('path');
let isBinaryFile = require('isbinaryfile');


function makeSafeForReplace(content) {
    content = content.replace(/(?:\r\n|\r)/g, '\n');

     // Dollar signs are special .replace() parameters
    // (https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/String/replace#Specifying_a_string_as_a_parameter#Specifying_a_string_as_a_parameter)
    // so we need to avoid any mishaps, so we replace every dollar sign with two dollar signs
    content = content.replace(/\$/g, "$$$");

    return content;
}

/**
 * Reads an entire file
 * @param {string} filePath Path to the file
 * @returns {string} The contents of the file
 */
function readEntireFile(filePath) {
    return new Promise((resolve, reject) => {
        // Read the contents and resolve the promise
        fs.readFile(filePath, { encoding: 'utf-8' }, (err, content) => {
            if (err != null)
                reject(err);

            resolve(makeSafeForReplace(content));
        });
    });
}


/**
 *
 * @param {string} root Root path
 * @param {string} globString glob format string
 * @returns things found in the glob
 */
function globFiles(root, globString) {
    return new Promise((resolve, reject) => {
        try {
            let globPath = path.join(root, globString);

            // Get all the files
            glob(globPath, (err, files) => {
                if (err != null)
                    reject(err);

                let result = [];
                for (let i = 0; i < files.length; i++) {
                    result.push({
                        absolute: path.normalize(files[i]),
                        relative: path.relative(root, files[i]),
                        absoluteDirname: path.dirname(files[i]),
                        relativeDirname: path.dirname(path.relative(root, files[i])),
                        name: path.basename(files[i])
                    });
                }

                resolve(result);
            });

        } catch (e) {
            resolve();
        }
    });




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
 * Copies a file to the destination
 * @param {string} destination Target file
 * @param {string} source source file
 */
function copyFile(destination, source) {
    fs.copyFileSync(source, destination);
}




module.exports = {
    readEntireFile,
    globFiles,
    canFileBeParsed,
    copyFile,
    makeSafeForReplace
}