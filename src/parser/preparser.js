let glob = require('glob');
let path = require('path');
let debug = require('../debugger')('preparser');
let fileUtils = require('../utils/fileUtils.js');


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
 * Loads from the config.compiler.data_folder into the site attribute
 * @param {projectConfig} config project configuration
 */
function loadData(config) {
    debug('Loading Data');

    let root = config.compiler.data_folder_absolute;
    let files = fileUtils.globFiles(root, '**/*.json');
    if (files.length > 0)
        config.site.data = {};

    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        if (config.compiler.ignore.indexOf(file.absolute) != -1)
            continue;

        // Create namespace for found file
        let namespace = _convertFolderStructureToNamespace(config.site.data, file.relativeDirname);
        let key = file.name.replace('.json', '');

        // Try to read, and parse the contents of the files
        try {
            let contents = fileUtils.readEntireFileSync(file);
            namespace[key] = JSON.parse(contents);
            debug(`\tFound data '${key}'`);
        }
        catch (e) {
            console.warn(e);
            namespace[key] = null;
        }
    }
}


/**
 * Load the templates available
 * @param {projectConfig} config project configuration
 */
function loadTemplates(config) {
    debug('Loading Templates');

    let root = config.compiler.templates_folder_absolute;

    let files = fileUtils.globFiles(root, '*.*');
    if (files.length > 0) {
        config.site.templates = [];
        config.__cache.templates = {};
    }

    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        let name = (file.name).replace(path.extname(file.name), '');

        config.site.templates.push(name);

        config.__cache.templates[name] = {
            info: file,
            contents: fileUtils.readEntireFileSync(file.absolute)
        };
    }
}

/**
 * Load files in the includes directory
 * @param {projectConfig} config project configuration
 */
function loadIncludes(config) {
    debug('Loading Includes');

    let root = config.compiler.includes_folder_absolute;

    let files = fileUtils.globFiles(root, '*.*');
    if (files.length > 0) {
        config.site.includes = [];
        config.__cache.includes = {};
    }

    for (let i = 0; i < files.length; i++) {
        let file = files[i];
        let name = (file.name).replace(path.extname(file.name), '');

        config.site.includes.push(name);

        config.__cache.includes[name] = {
            info: file,
            contents: fileUtils.readEntireFileSync(file.absolute)
        };
    }
}


/**
 * Load plugins the project is using
 * @param {projectConfig} config project configuration
 */
function loadPlugins(config) {
    if (!config.compiler.allow_plugins)
        return;

    debug('Loading Plugins');

    let root = config.compiler.plugins_folder_absolute;
    let plugins = fileUtils.globFiles(root, '*/plugin.json');
    if (plugins.length > 0) {
        config.site.plugins = [];
        config.__cache.plugins = {};
    }

    for (let i = 0; i < plugins.length; i++) {
        let name = plugins[i].relativeDirname;

        config.site.plugins.push(name);

        config.__cache.plugins[name] = {
            info: plugins[i],
            manifest: fileUtils.readEntireFileSync(plugins[i].absolute)
        };
    }
}


/**
 * Finds folders that begin with underscore and adds that to the site properties
 * @param {projectConfig} config project configuration file
 * @returns {string[]} array of new properties
 */
function _discoverSiteProperties(config) {
    debug('Discovering site properties');
    let root = config.site.source_absolute;
    let result = [];

    let ignore = [
        config.compiler.plugins_folder, config.compiler.templates_folder,
        config.compiler.includes_folder, config.compiler.data_folder,
        config.site.output
    ];

    // Get folders that begin with '_' at the root of the project
    let found = fileUtils.globFiles(root, '_*/');
    for (let i = 0; i < found.length; i++) {
        let folder = found[i];

        // Don't use certain folders
        if (ignore.indexOf(folder.name) != -1)
            continue;

        // Don't use ignored folders
        if (config.compiler.ignore_absolute.indexOf(folder.absolute) != -1)
            continue;

        result.push(folder.name);
        debug(`\tFound property '${folder.name}'`);
        //config.site[relativeName.replace('_', '')] = []; // It's an array
    }

    return result;
}


/**
 * Discovers all of the files that should be compiled
 * @param {projectConfig} config project configuration
 */
function discoverFiles(config) {
    let newLocations = _discoverSiteProperties(config);

    let discover = (cfg, key, directory) => {
        debug(`Exploring '${directory}'`);

        if (cfg[key] == undefined || !Array.isArray(cfg[key]))
            cfg[key] = [];

        let root = path.join(config.site.source, directory);

        let found = fileUtils.globFiles(root, '**/*.*');
        config.compiler.__files = found;
        for (let i = 0; i < found.length; i++) {
            debug(`\tFound: ${found[i].relative}`);
        }
    };

    //discover(config.site, 'plugins', config.compiler.plugins_folder);
    //discover(config.site, 'templates', config.compiler.templates_folder);
    //discover(config.site, 'includes', config.compiler.includes_folder);

    for (let i = 0; i < newLocations.length; i++) {
        discover(config.site, newLocations[i].replace('_', ''), newLocations[i]);
    }
}



module.exports = {
    loadData,
    loadTemplates,
    loadIncludes,
    loadPlugins,
    discoverFiles
};