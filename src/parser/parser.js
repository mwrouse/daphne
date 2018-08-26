let fileUtils = require('../utils/fileUtils');
let directoryUtils = require('../utils/directoryUtils');
let debugParser = require('../debugger')('parser');
let fs = require('fs');
let path = require('path');


/**
 * First pass of the parser
 * First pass discovers metadata about the file
 * @param {projectConfig} config project config
 */
function firstPass(config) {
    debugParser('First Pass');
    let debug = debugParser.new('first');

    for (let i = 0; i < config.__cache.files.length; i++) {
        fileUtils.getMetadataHeader(config.__cache.files[i], config, debug);
    }
}



/**
 * Second pass of the parser
 * Second pass actually expands the body of the file
 * @param {projectConfig} config project config
 */
function secondPass(config) {
    debugParser('Second Pass');
    let debug = debugParser.new('second');

    for (let i = 0; i < config.__cache.files.length; i++) {
        let cache = config.__cache.files[i];

        directoryUtils.createDirectoryStructure(config.site.output_absolute, cache.info.relativeDirname);

        let finalPath = path.join(config.site.output_absolute, cache.info.relative);

        if (!cache.shouldParse) {
            fileUtils.copyFile(finalPath, cache.info.absolute);
            continue;
        }

        // Parse the file and write it
        fs.writeFileSync(finalPath, cache.content);
    }
}



module.exports = {
    firstPass,
    secondPass
};