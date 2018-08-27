let debug = require('../debugger')('parser:second');
let fs = require('fs');
let path = require('path');
let fileUtils = require('../utils/fileUtils');
let directoryUtils = require('../utils/directoryUtils');

/**
 * Second pass of the parser
 * Second pass actually expands the body of the file
 * @param {projectConfig} config project config
 */
function secondPass(config) {
    debug('Second Pass');

    for (let i = 0; i < config.__cache.files.length; i++) {
        let cache = config.__cache.files[i];

        directoryUtils.createDirectoryStructure(config.site.output_absolute, cache.info.relativeDirname);

        let finalPath = path.join(config.site.output_absolute, cache.info.relative);

        if (!cache.shouldParse) {
            fileUtils.copyFile(finalPath, cache.info.absolute);
            continue;
        }

        // Parse the file and write it
        fs.writeFileSync(finalPath, cache.contents);
    }
}

module.exports = secondPass;