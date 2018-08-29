let utils = require('./utils');
let parser = require('./parser');
let path = require('path');
let debug = require('./debugger')('\b');
let projectConfig = require('./managers/projectConfig');
let templateManager = require('./managers/templateManager');
let postManager = require('./managers/postManager');
let dataManager = require('./managers/dataManager');
let includesManager = require('./managers/includesManager');
let fileManager = require('./managers/otherFileManager');

/**
 * Reads an entire config file and parses it as JSON
 * @param {string} filePath Path to config file
 * @returns {json} parsed file
 */
function readAndParseConfigFile(filePath) {
    var fileContents = utils.files.readEntireFileSync(filePath);

    return JSON.parse(fileContents);
}

// Fatal error
function fatal(e) {
    console.log('\n');
    console.error(e);
    console.log('\n');
    process.exit();
}


/**
 * Main class for static site generator
 */
class Daphne {

    constructor(projectPath) {
        this._projectPath = projectPath;

    }


    /**
     * Begins the entire process of building a static website
     */
    buildSite() {
        this._preBuild()
            .then(() => {
                //console.log(JSON.stringify(postManager.posts));
                return parser.parse(this._projectConfig);
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

        try {
            projectConfig.setProjectRoot(this._projectPath);
        } catch (e) {
            fatal(e);
        }


        let preBuildRoutine = [
            templateManager.load(),
            postManager.load(),
            dataManager.load(),
            includesManager.load(),
            fileManager.load(),
        ];

        return Promise.all(preBuildRoutine)
            .catch((error) => {
                fatal(error);
            });

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