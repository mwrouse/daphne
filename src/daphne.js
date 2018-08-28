let utils = require('./utils');
let parser = require('./parser');
let path = require('path');
let debug = require('./debugger')('\b');
let projectConfig = require('./managers/projectConfig');
let templateManager = require('./managers/templateManager');
let postManager = require('./managers/postManager');

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

    constructor(projectPath) {
        this._projectPath = projectPath;
       // this._projectConfig = {};
        //this._projectConfigFile = '';
    }


    /**
     * Begins the entire process of building a static website
     */
    buildSite() {
        this._preBuild()
            .then(() => {
                console.log(postManager.posts);
                return parser.parse(this._projectConfig);
            })
            .then(() => {
                console.log('Done!');
            });
        //console.log(this._projectConfig.__cache.site);
    }


    /**
     * Watches a project for changes and builds when changes are made
     */
    watchSite() {
        this._preBuild();
        console.log("Watching");
    }


    /**
     * Watches a project for changes while serving a localhost website
     */
    serveSite() {
        this._preBuild();
        console.log("Serving");
    }


    /**
     * PreBuild task (setups everything needed)
     */
    _preBuild() {
        debug(`Building project at ${this._projectPath}`);

        if (!utils.config.doesProjectHaveConfigFile(this._projectPath))
            throw new Error(`No 'config.daphne' in ${this._projectPath}`);

        //this._projectConfigFile = path.join(this._projectPath, 'config.daphne');

        projectConfig.setProjectRoot(this._projectPath);

        let routine = [
            templateManager.loadTemplates(),
            postManager.loadPosts(),
        ];

        return Promise.all(routine);
       // return Promise.resolve();
/*
        // Parse and apply defaults
        this._projectConfig = utils.config.parseConfigFile(this._projectConfigFile);
        this._projectConfig.compiler.root = this._projectPath;
        utils.config.applyDefaultConfiguration(this._projectConfig);

        // Clean output directory
        utils.directories.removeFolder(this._projectConfig.site.output_absolute);

        // Preparse
        return parser.preparse(this._projectConfig);*/
    }

}




module.exports = Daphne;