let path = require('path');
let debug = require('../debugger')('data_manager');
let fileUtils = require('../utils/fileUtils');
let FileToCompile = require('../models/fileToCompile');
let config = require('./projectConfig');

var __dataCache = {};

/**
 * Converts a folder structure to object namespace,
 * like /foo/bar/baz would become foo.bar.baz on an object
 * @param {object} root Where to place namespace
 * @param {string} folderPath folder structure
 */
function _convertFolderStructureToNamespace(root, folderPath) {
    if (folderPath == '.' || folderPath == '')
        return root;

    let paths = folderPath.split(path.sep);
    let next_namespace = paths.shift();
    root[next_namespace] = {};

    return _convertFolderStructureToNamespace(root[next_namespace], paths.join(path.sep));
}

/**
 * Loads data from a file into the proper namespace
 * @param {string} file
 */
function _loadData(file) {
    let namespace = _convertFolderStructureToNamespace(__dataCache, file.relativeDirname);
    let key = file.name.replace('.json', ''); // Namespace key

    namespace[key] = null;

    return fileUtils.readEntireFile(file.absolute)
        .then((content) => {
            try {
                namespace[key] = JSON.parse(content);
                debug(`Found data '${key}'`);
            }
            catch (e) {
                console.warn(e);
            }
        });
}



class DataManager {
    constructor() {

    }

    /**
     * Load all of the site data
     */
    loadData() {
        debug('Loading Data');
        let root = config.compiler.data_folder_absolute;

        return fileUtils.globFiles(root, '**/*.json')
            .then((data_files) => {
                let waiting = [];

                for (let i = 0; i < data_files.length; i++) {
                    let file = data_files[i];

                    if (config.compiler.ignore_absolute.indexOf(file.absolute) != -1)
                        continue; // File is ignored

                    waiting.push(
                        _loadData(file)
                    );
                }

                return Promise.all(waiting);
            });
    }


    /**
     * Return the data
     */
    get data() {
        return JSON.parse(JSON.stringify(__dataCache)); // Deep copy
    }
}


let instance = new DataManager();
module.exports = instance;