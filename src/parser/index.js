let twoPassParser = require('./parser.js');
let preparser = require('./preparser.js');


/**
 * Performs the two pass parsing
 * @param {project config} config configuration
 */
function parse(config) {
    twoPassParser.firstPass(config);
    twoPassParser.secondPass(config);
}


module.exports = {
    preparser,
    parse
};