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
        debug(`\tFound ${name}`);
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
        debug(`\tFound ${name}`);
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
        debug(`\tFound ${name}`);
    }
}



/**
 * Loads all files in a properties folder
 * @param {projectConfig} config configuration
 * @param {fileInfo} prop information about the folder
 */
function _readCustomProperties(config, prop) {
    let key = prop.name.replace('_', '');

    let found = fileUtils.globFiles(prop.absolute, '*.*');
    if (found.length > 0) {
        config[key] = [];
    }

    for (let i = 0; i < found.length; i++) {
        let file = found[i];
        let name = (file.name).replace(path.extname(file.name), '');

        config[key].push({
            info: file,
            contents: fileUtils.readEntireFileSync(file.absolute)
        });
    }
}

/**
 * Load custom site properties (folders starting with _)
 * @param {projectConfig} config configuration
 */
function loadCustomProperties(config) {
    debug('Loading Custom Site Properties');

    let ignore = [
        config.compiler.plugins_folder, config.compiler.templates_folder,
        config.compiler.includes_folder, config.compiler.data_folder,
        config.site.output
    ];

    let root = config.site.source_absolute;
    let properties = fileUtils.globFiles(root, '_*/');
    if (properties.length > 0) {
        config.site.__properites = [];
    }

    for (let i = 0; i < properties.length; i++) {
        let prop = properties[i];
        let name = prop.name.replace('_', '');
        // Don't use certain folders
        if (ignore.indexOf(name) != -1)
            continue;

        // Don't use ignored folders
        if (config.compiler.ignore_absolute.indexOf(prop.absolute) != -1)
            continue;

        config.site.__properites.push(name);
        _readCustomProperties(config.__cache.site, prop);
        debug(`\tFound ${prop}`);
    }
}


/**
 * Discovers all of the files that should be compiled
 * @param {projectConfig} config project configuration
 */
function discoverFiles(config) {
    debug('Loading other files');

    let root = config.site.source_absolute;
    let found = fileUtils.globFiles(root, `!(_)*/**/*.*`);
    let alsoFound = fileUtils.globFiles(root, '*.*');
    found = found.concat(alsoFound);
    if (found.length > 0) {
        config.__cache.files = [];
    }

    for (let i = 0; i < found.length; i++) {
        if (found[i].name == "config.daphne")
            continue;

        // Don't use ignored files
        if (config.compiler.ignore_absolute.indexOf(found[i].absolute) != -1)
            continue;

        let extenion = path.extname(found[i].name).replace('.', '');
        let parse = (fileUtils.canFileBeParsed(found[i].absolute) == true && config.compiler.extensions_to_parse.indexOf(extenion) != -1);
        config.__cache.files.push({
            info: found[i],
            shouldParse: parse,
            content: (parse) ? fileUtils.readEntireFileSync(found[i].absolute) : null
        });
        debug(`\tFound ${found[i].relative}`);
    }


}



module.exports = {
    loadData,
    loadTemplates,
    loadIncludes,
    loadPlugins,
    loadCustomProperties,
    discoverFiles
};