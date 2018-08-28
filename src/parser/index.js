let preparser = require('./preparser.js');

let firstPass = require('./firstPass');
let secondPass = require('./secondPass');

/**
 * Performs the two pass parsing
 * @param {project config} config configuration
 */
function parse(config) {
    return Promise.resolve();/*firstPass(config)
            .then(() => {
                secondPass(config);
            });*/
    /*
    return new Promise((resolve, reject) => {
        //twoPassParser.firstPass(config);
        //twoPassParser.secondPass(config);

        resolve();
    });*/
}


/**
 * Perfrom the preparse operations
 * @param {project config} config
 */
function preparse(config) {
    return Promise.all([
        preparser.loadData(config),
        preparser.loadTemplates(config),
        preparser.loadIncludes(config),
        preparser.loadPlugins(config),
        preparser.loadCustomProperties(config),
        preparser.loadPosts(config),
        preparser.discoverFiles(config),
    ]);
}


module.exports = {
    preparse,
    parse
};