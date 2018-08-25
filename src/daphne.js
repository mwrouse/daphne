let utils = require('./utils');
let parser = require('./parser');
let path = require('path');
let debug = require('./debugger')('root');


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
        this._projectPath = '';
        this._projectConfig = {};
        this._projectConfigFile = '';
    }


    /**
     * Begins the entire process of building a static website
     * @param {string} projectPath Path to project root
     */
    buildSite(projectPath) {
        this._preBuild(projectPath);

    }


    /**
     * Watches a project for changes and builds when changes are made
     * @param {string} projectPath Path to the project root
     */
    watchSite(projectPath) {
        this._preBuild(projectPath);
        console.log("Watching");
    }


    /**
     * Watches a project for changes while serving a localhost website
     * @param {string} projectPath Path to project root
     */
    serveSite(projectPath) {
        this._preBuild(projectPath);
        console.log("Serving");
    }


    /**
     * PreBuild task (setups everything needed)
     * @param {string} projectPath Path to project root
     */
    _preBuild(projectPath) {
        debug(`Build project at ${projectPath}`);

        this._projectPath = projectPath;

        if (!utils.config.doesProjectHaveConfigFile(this._projectPath))
            throw new Error(`No 'config.daphne' in ${this._projectPath}`);

        this._projectConfigFile = path.join(this._projectPath, 'config.daphne');

        // Parse and apply defaults
        this._projectConfig = utils.config.parseConfigFile(this._projectConfigFile);
        this._projectConfig.compiler.root = this._projectPath;
        utils.config.applyDefaultConfiguration(this._projectConfig);

        // Preparse
        parser.preparser.loadData(this._projectConfig);
        parser.preparser.loadTemplates(this._projectConfig);
        parser.preparser.loadIncludes(this._projectConfig);
        parser.preparser.loadPlugins(this._projectConfig);
        //parser.preparser.discoverFiles(this._projectConfig);
    }

}




module.exports = Daphne;