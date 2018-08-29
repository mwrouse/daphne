let path = require('path');
let debug = require('../debugger')('file_manager');
let fileUtils = require('../utils/fileUtils');
let FileToCompile = require('../models/fileToCompile');
let config = require('./projectConfig');

var __fileCache = [];

/**
 * Class for managing other files
 */
class OtherFileManager {
    constructor(){
    }

    /**
     * Load all files
     */
    load() {
        debug('Loading Other Files');

        let root = config.site.source_absolute;

        return fileUtils.globFiles(root, '{!(_)*/**/,}*.*')
            .then((files) => {
                let waiting = [];
                let index = 0;

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];

                    if (config.isFileIgnored(file.absolute))
                        continue; // Ignore the file

                        waiting.push(
                            fileUtils.readEntireFile(file.absolute)
                                .then((content) => {
                                    debug(`Found file '${file.relative}'`);
                                    return new FileToCompile(file.absolute, content);
                                })
                        );
                }

                return Promise.all(waiting)
                    .then((loadedFiles) => {
                        __fileCache.concat(loadedFiles);
                    });
            });
    }
}


let instance = new OtherFileManager();
module.exports = instance;