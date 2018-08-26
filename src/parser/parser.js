let fileUtils = require('../utils/fileUtils.js');
let debugParser = require('../debugger')('parser');


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

    }
}



module.exports = {
    firstPass,
    secondPass
};