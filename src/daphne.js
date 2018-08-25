let utils = require('./utils');
const path = require('path');

/**
 * Reads an entire config file and parses it as JSON
 * @param {string} filePath Path to config file
 * @returns {json} parsed file
 */
function readAndParseConfigFile(filePath) {
    var fileContents = utils.files.readEntireFileSync(filePath);

    return JSON.parse(fileContents);
}

/**
 * Main class for static site generator
 */
class Daphne {

    constructor() {
        this._ProjectPath = '';
        this._ProjectConfigFile = '';
    }


    /**
     * Begins the entire process of building a static website
     * @param {string} projectPath
     */
    buildSite(projectPath) {
        this._ProjectPath = projectPath;

        if (!utils.directories.doesDirectoryHaveConfigFile(this._ProjectPath))
            throw new Error(`No 'config.daphne' in ${this._ProjectPath}`);

        this._ProjectConfigFile = path.join(this._ProjectPath, 'config.daphne');
        this._ProjectConfig = readAndParseConfigFile(this._ProjectConfigFile);

        console.log(this._ProjectConfig.site);
    }

}




module.exports = Daphne;