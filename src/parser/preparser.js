let glob = require('glob');
let path = require('path');
let preparserDebug = require('../debugger')('preparser');
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
 * Loads individual data
 * @param {object} namespace
 * @param {string} key
 * @param {string} file
 */
function _loadData(namespace, key, file, debug) {
    return new Promise((resolve, reject) => {
        namespace[key] = null;

        fileUtils.readEntireFile(file)
            .then((contents) => {
                try {
                    namespace[key] = JSON.parse(contents);
                    debug(`Found data '${key}'`);
                } catch (e) {
                    console.warn(e);
                }
            }).catch((err) => {
                console.error(err);
            }).finally(() => {
                resolve(); // Resolve no matter what
            });
    });
}


/**
 * Loads from the config.compiler.data_folder into the site attribute
 * @param {projectConfig} config project configuration
 */
function loadData(config) {
    let debug = preparserDebug.new('data');
    debug('Loading Data');

    config.site.data = {};

    let root = config.compiler.data_folder_absolute;

    return fileUtils.globFiles(root, '**/*.json')
            .then((files) => {
                let awaiting = [];
                for (let i = 0; i < files.length; i++) {
                    let file = files[i];
                    if (config.compiler.ignore.indexOf(file.absolute) != -1)
                        continue;

                    // Create namespace for found file
                    let namespace = _convertFolderStructureToNamespace(config.site.data, file.relativeDirname);
                    let key = file.name.replace('.json', '');

                    // Try to read, and parse the contents of the files
                    awaiting.push(_loadData(namespace, key, file.absolute, debug));
                }

                // Wait for all of them
                return Promise.all(awaiting);
            });
}


/**
 * Load an individual template
 * @param {object} namespace
 * @param {string} key
 * @param {string} file
 */
function _loadTemplate(namespace, key, file, debug) {
    namespace[key].contents = null;

    return fileUtils.readEntireFile(file)
        .then((content) => {
            namespace[key].contents = content;
            debug(`Found ${key}`);
        });
}

/**
 * Load the templates available
 * @param {projectConfig} config project configuration
 */
function loadTemplates(config) {
    let debug = preparserDebug.new('templates');
    debug('Loading Templates');

    config.site.templates = [];
    config.__cache.templates = {};

    let root = config.compiler.templates_folder_absolute;

    return fileUtils.globFiles(root, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];
                    let name = (file.name).replace(path.extname(file.name), '');

                    config.site.templates.push(name);

                    config.__cache.templates[name] = {
                        info: file,
                        contents: null
                    };

                    awaiting.push(_loadTemplate(config.__cache.templates, name, file.absolute, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Loads an individual include file
 * @param {object} namespace
 * @param {string} key
 * @param {string} file
 * @param {debug logger} debug
 */
function _loadInclude(namespace, key, file, debug) {
    namespace[key].contents = null;

    return fileUtils.readEntireFile(file)
            .then((contents) => {
                namespace[key].contents = contents;
                debug(`Found ${key}`);
            });
}

/**
 * Load files in the includes directory
 * @param {projectConfig} config project configuration
 */
function loadIncludes(config) {
    let debug = preparserDebug.new('includes');
    debug('Loading Includes');

    config.site.includes = [];
    config.__cache.includes = {};

    let root = config.compiler.includes_folder_absolute;

    return fileUtils.globFiles(root, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];
                    let name = (file.name).replace(path.extname(file.name), '');

                    config.site.includes.push(name);
                    config.__cache.includes[name] = {
                        info: file,
                        contents: null
                    };

                    awaiting.push(_loadInclude(config.__cache.includes, name, file.absolute, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Loads an individual plugin
 * @param {object} namespace
 * @param {string} key
 * @param {string} file
 * @param {debugger} debug
 */
function _loadPlugin(namespace, key, file, debug) {
    namespace[key].manifest = null;

    return fileUtils.readEntireFile(file)
            .then((contents) => {
                try {
                    namespace[key].manifest = JSON.parse(contents);
                } catch (e) {
                    console.warn(e);
                }

                debug(`Found ${key}`);
            });
}

/**
 * Load plugins the project is using
 * @param {projectConfig} config project configuration
 */
function loadPlugins(config) {
    config.site.plugins = [];
    config.__cache.plugins = {};

    if (!config.compiler.allow_plugins)
        return;

    let debug = preparserDebug.new('plugins');
    debug('Loading Plugins');

    let root = config.compiler.plugins_folder_absolute;

    return fileUtils.globFiles(root, '*/plugin.json')
            .then((plugins) => {
                let awaiting = [];

                for (let i = 0; i < plugins.length; i++) {
                    let plugin = plugins[i];
                    let name = plugin.relativeDirname;

                    config.site.plugins.push(name);

                    config.__cache.plugins[name] = {
                        info: plugin,
                        manifest: null
                    };

                    awaiting.push(_loadPlugin(config.__cache.plugins, name, plugin.absolute, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Loads a file for a property
 * @param {object} namespace
 * @param {string} key
 * @param {string} file
 * @param {debugger} debug
 */
function __loadCustomProperty(namespace, key, name, file, debug) {
    namespace[key].contents = null;

    return fileUtils.readEntireFile(file)
            .then((contents) => {
                try {
                    namespace[key].contents = contents;
                } catch (e) {
                    console.warn(e);
                }

                debug(`Found ${name}`);
            });
}

/**
 * Loads all files in a properties folder
 * @param {object} namespace configuration
 * @param {string} key
 */
function _loadCustomProperty(namespace, key, prop, debug) {
    namespace[key] = [];

    return fileUtils.globFiles(prop.absolute, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];

                    namespace[key].push({
                        info: file,
                        contents: null
                    });

                    awaiting.push(__loadCustomProperty(namespace[key], i, key + '.' + file.name, file.absolute, debug));
                }

                debug(`Found ${key}`);

                return Promise.all(awaiting);
            });
}

/**
 * Load custom site properties (folders starting with _)
 * @param {projectConfig} config configuration
 */
function loadCustomProperties(config) {
    let debug = preparserDebug.new('properties');
    debug('Loading Custom Site Properties');

    config.site.__properites = [];
    config.__cache.site = {};

    let ignore = [
        config.compiler.plugins_folder, config.compiler.templates_folder,
        config.compiler.includes_folder, config.compiler.data_folder,
        config.site.output
    ];

    let root = config.site.source_absolute;

    return fileUtils.globFiles(root, '_*/')
            .then((folders) => {
                // TODO: Verify it is a folder
                let awaiting = [];

                for (let i = 0; i < folders.length; i++) {
                    let prop = folders[i];
                    let name = prop.name.replace('_', '');

                    // Don't use certain folders
                    if (ignore.indexOf(prop.name) != -1)
                        continue;

                    // Don't use ignored folders
                    if (config.compiler.ignore_absolute.indexOf(prop.absolute) != -1)
                        continue;

                    config.site.__properites.push(name);

                    awaiting.push(_loadCustomProperty(config.__cache.site, name, prop, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Discovers all of the files that should be compiled
 * @param {projectConfig} config project configuration
 */
function discoverFiles(config) {
    let debug = preparserDebug.new('files');
    debug('Loading other files');

    let root = config.site.source_absolute;
    //let found = fileUtils.globFiles(root, `!(_)*/**/*.*`);
    //let alsoFound = fileUtils.globFiles(root, '*.*');
    /*found = found.concat(alsoFound);
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
            content: null //(parse) ? fileUtils.readEntireFileSync(found[i].absolute) : null
        });
        debug(`Found ${found[i].relative}`);
    }*/

    return Promise.resolve(config);
}



module.exports = {
    loadData,
    loadTemplates,
    loadIncludes,
    loadPlugins,
    loadCustomProperties,
    discoverFiles
};