let glob = require('glob');
let path = require('path');
let debug = require('../debugger')('preparser');
let fileUtils = require('../utils/fileUtils.js');

/**
 * Converts a folder structure to object namespace,
 * like /foo/bar/baz would become foo.bar.baz on an object
 * @param {object} configRoot Where to place namespace
 * @param {string} folderPath folder structure
 */
function convertFolderStructureToNamespace(configRoot, folderPath) {
    if (folderPath == '.' || folderPath == '')
        return configRoot;

    let paths = folderPath.split(path.sep);
    let next_namespace = paths.shift();
    configRoot[next_namespace] = {};

    return convertFolderStructureToNamespace(configRoot[next_namespace], paths.join(path.sep));
}

/**
 * Loads from the config.compiler.data_folder into the site attribute
 * @param {projectConfig} config project configuration
 */
function loadData(config) {
    debug('Loading data');

    let dataPath = path.join(config.compiler.root, config.compiler.data_folder);
    let dataGlobPath = path.join(dataPath, '**/*.json');

    let files = glob.sync(dataGlobPath);
    if (files.length > 0)
        config.site.data = {};

    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        if (config.compiler.ignore.indexOf(file) != -1)
            continue;

        let relativePath = path.relative(dataPath, file); // Convert to relative path
        let folder = path.dirname(relativePath);

        // Create namespace for found file
        let namespace = convertFolderStructureToNamespace(config.site.data, folder);
        let key = path.basename(relativePath, '.json');

        // Try to read, and parse the contents of the files
        try {
            let contents = fileUtils.readEntireFileSync(file);
            namespace[key] = JSON.parse(contents);
        }
        catch (e) {
            console.warn(e);
            namespace[key] = null;
        }
    }
}


/**
 * Discovers all of the files that are waiting to be compiled
 * @param {projectConfig} config project configuration
 */
function discoverFiles(config) {
    debug('Discovering files');

    let discover = (cfg, key, directory) => {

    };

    discover(config.compiler, 'plugins', config.compiler.plugins_folder);
    discover(config.compiler, 'templates', config.compiler.templates_folder);
    discover(config.compiler, 'includes', config.compiler.includes_folder);

}



module.exports = {
    loadData,
    discoverFiles
};