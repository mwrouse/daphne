let fs = require('fs');
let rimraf = require('rimraf');
let path = require('path');

/**
 * Recursively create a folder
 * @param {string} folder Folder to create
 */
function _recursiveCreate(folder) {
    if (fs.existsSync(folder))
        return true;

    _recursiveCreate(path.join(folder, '..'));
    fs.mkdirSync(folder);
}


/**
 * Makes folders for a directory tree
 * @param {string} buildPathAbsolute build output folder
 * @param {string} structure Structure of the folders
 */
function createDirectoryStructure(buildPathAbsolute, structure) {
    let deepestFolder = path.join(buildPathAbsolute, structure);
    _recursiveCreate(deepestFolder);
}


/**
 * Recursively delete a folder, until root is reached
 * @param {string} root
 * @param {string} folder
 */
function _recursiveDelete(root, folder) {
    let fullPath = path.join(root, folder);
    if (!fs.existsSync(fullPath))
        return true;

    if (folder != '.') {
        let nextFolder = path.dirname(folder);
        _recursiveDelete(root, nextFolder);
    }
}


/**
 * Removes a folder, and anything in it
 * @param {string} folder
 */
function removeFolder(folder) {
    rimraf.sync(path.join(folder, '**/*'));
}


module.exports = {
    createDirectoryStructure,
    removeFolder
};