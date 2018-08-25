let fs = require('fs');


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



module.exports = {
    readEntireFileSync
}