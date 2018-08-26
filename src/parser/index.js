let twoPassParser = require('./parser.js');
let preparser = require('./preparser.js');


/**
 * Performs the two pass parsing
 * @param {project config} config configuration
 */
function parse(config) {
    return new Promise((resolve, reject) => {
        //twoPassParser.firstPass(config);
        //twoPassParser.secondPass(config);

        resolve();
    });
}


/**
 * Perfrom the preparse operations
 * @param {project config} config
 */
function preparse(config) {
    let promises = [
        preparser.loadData(config),
        preparser.loadTemplates(config),
        preparser.loadIncludes(config),
        preparser.loadPlugins(config),
        preparser.loadCustomProperties(config),
        preparser.discoverFiles(config),
    ];

    return Promise.all(promises).then(() => {
        console.log('Done!');
    });
}


module.exports = {
    preparse,
    parse
};