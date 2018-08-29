let path = require('path');
let debug = require('../debugger')('includes_manager');
let fileUtils = require('../utils/fileUtils');
let FileToCompile = require('../models/fileToCompile');
let config = require('./projectConfig');

var __includesCache = {};

/**
 * Load file
 * @param {string} file
 */
function _loadIncludes(file) {
    __includesCache[file.relative] = null;

    return fileUtils.readEntireFile(file.absolute)
        .then((content) => {
            __includesCache[file.realtive] = new FileToCompile(file.absolute, content);
        });
}


class IncludesManager {
    constructor() {

    }

    /**
     * Load all files that could be included in others
     */
    load() {
        debug('Loading Includes');
        let root = config.compiler.includes_folder_absolute;

        return fileUtils.globFiles(root, '**/*.*')
            .then((files) => {
                let waiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];

                    if (config.isFileIgnored(file.absolute))
                        continue; // Ignored

                    waiting.push(
                        _loadIncludes(file)
                    );
                }

                return Promise.all(waiting);
            });
    }


    /**
     * Check if an include is valid
     * @param {string} filePath
     */
    isIncludeValid(filePath) {
        filePath = path.normalize(filePath);
        return __includesCache[filePath] != undefined;
    }

    /**
     * Get the file of an include
     * @param {stirng} filePath
     */
    getInclude(filePath) {
        filePath = path.normalize(filePath);

        if (!this.isIncludeValid(filePath))
            return null;

        return __includesCache[filePath];
    }
}


let instance = new IncludesManager();
module.exports = instance;