let glob = require('glob');
let path = require('path');
let preparserDebug = require('../debugger')('preparser');
let fileUtils = require('../utils/fileUtils.js');
let pluginManager = require('../utils/pluginManager');


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

//
/**
 * Generic file loader
 * @param {object} namespace
 * @param {string} key
 * @param {string} filePath
 * @param {debugger} debug
 * @param {string} contentsKey
 * @param {function} modifier
 */
function _loadFile(namespace, key, filePath, debug, contentsKey, modifier) {
    contentsKey = (contentsKey == undefined) ? 'contents' : contentsKey;
    modifier = (modifier == undefined || typeof(modifier) != 'function') ? null : modifier;

    (namespace[key])[contentsKey] = null;

    return fileUtils.readEntireFile(filePath)
            .then((contents) => {
                try {
                    if (modifier != null)
                        (namespace[key])[contentsKey] = modifier(contents);
                    else
                        (namespace[key])[contentsKey] = contents;
                } catch (e) {
                    console.warn(e);
                }

                debug(`Found ${filePath}`);
            });
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
 * Load the templates available
 * @param {projectConfig} config project configuration
 */
function loadTemplates(config) {
    let debug = preparserDebug.new('templates');
    debug('Loading Templates');

    config.__site.templates = [];
    config.__cache.templates = {};

    let root = config.compiler.templates_folder_absolute;

    return fileUtils.globFiles(root, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];
                    let name = (file.name).replace(path.extname(file.name), '');

                    config.__site.templates.push(name);
                    config.__cache.templates[name] = {
                        info: file,
                        shouldParse: true,
                        contents: null
                    };

                    awaiting.push(_loadFile(config.__cache.templates, name, file.absolute, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Load files in the includes directory
 * @param {projectConfig} config project configuration
 */
function loadIncludes(config) {
    let debug = preparserDebug.new('includes');
    debug('Loading Includes');

    config.__site.includes = [];
    config.__cache.includes = {};

    let root = config.compiler.includes_folder_absolute;

    return fileUtils.globFiles(root, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];
                    let name = file.name;

                    config.__site.includes.push(name);
                    config.__cache.includes[name] = {
                        info: file,
                        shouldParse: true,
                        contents: null
                    };

                    awaiting.push(_loadFile(config.__cache.includes, name, file.absolute, debug));
                }

                return Promise.all(awaiting);
            });
}


/**
 * Load plugins the project is using
 * @param {projectConfig} config project configuration
 */
function loadPlugins(config) {
    config.__site.plugins = [];
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

                    awaiting.push(pluginManager.register(name, plugin.absolute));
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
function _loadCustomProperty(namespace, key, prop, debug, config) {
    namespace[key] = [];

    return fileUtils.globFiles(prop.absolute, '*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];

                    // Don't use ignored folders
                    if (config.compiler.ignore_absolute.indexOf(file.absolute) != -1)
                        continue;

                    let extension = path.extname(file.name).replace('.', '');
                    let shouldBeParsed = (fileUtils.canFileBeParsed(file.absolute) && config.compiler.extensions_to_parse.indexOf(extension) != -1);

                    namespace[key].push({
                        info: file,
                        shouldParse: shouldBeParsed,
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

    let ignore = [
        config.compiler.plugins_folder, config.compiler.templates_folder,
        config.compiler.includes_folder, config.compiler.data_folder,
        config.site.output, config.compiler.posts_folder
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

                    config.__site.properties.push(name);

                    awaiting.push(_loadCustomProperty(config.__cache.site, name, prop, debug, config));
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

    return fileUtils.globFiles(root, '{!(_)*/**/,}*.*')
        .then((files) => {
            let awaiting = [];
            let counter = 0;

            for (let i = 0; i < files.length; i++) {
                let file = files[i];

                // Don't take ignored files :)
                if (config.compiler.ignore_absolute.indexOf(file.absolute) != -1)
                    continue;

                let extension = path.extname(file.name).replace('.', '');
                let shouldBeParsed = (fileUtils.canFileBeParsed(file.absolute) && config.compiler.extensions_to_parse.indexOf(extension) != -1);

                config.__cache.files.push({
                    info: file,
                    shouldParse: shouldBeParsed,
                    contents: null
                });

                awaiting.push(_loadFile(config.__cache.files, counter, file.absolute, debug));
                counter++;
            }

            return Promise.all(awaiting);
        });
}


/**
 * Still loads a post
 */
function __loadPost(config, key, file, debug) {
    return fileUtils.readEntireFile(file.absolute)
        .then((contents) => {
            let post = config.__cache.posts[key];
            let postFile = path.join(post.info.absolute, 'index.html');

            file.relative = path.relative(config.compiler.posts_folder_absolute, file.absolute);
            file.relativeDirname = path.dirname(file.relative);

            if (file.absolute == postFile) {
                config.__cache.posts[key].post.info = file;
                config.__cache.posts[key].post.contents = contents;
                debug(`Found post '${post.slug}'`);
            }
            else {
                config.__cache.posts[key].assets.push({
                    info: file,
                    shouldParse: true,
                    contents: contents
                });
            }

        });
}

/**
 * Loads a single post, and all of it's files
 */
function _loadPost(config, key, debug) {
    let post = config.__cache.posts[key];
    let postFile = path.join(post.info.absolute, 'index.html');

    return fileUtils.globFiles(post.info.absolute, '**/*.*')
            .then((files) => {
                let awaiting = [];

                for (let i = 0; i < files.length; i++) {
                    let file = files[i];

                    // Don't take ignored files :)
                    if (config.compiler.ignore_absolute.indexOf(file.absolute) != -1)
                        continue;

                    let extension = path.extname(file.name).replace('.', '');
                    let shouldBeParsed = (fileUtils.canFileBeParsed(file.absolute) && config.compiler.extensions_to_parse.indexOf(extension) != -1);

                    if (!shouldBeParsed) {
                        config.__cache.posts[key].assets.push({
                            info: file,
                            shouldParse: shouldBeParsed,
                            contents: null
                        });
                        continue;
                    }

                    awaiting.push(__loadPost(config, key, file, debug));
                }

                return Promise.all(awaiting);
            });
}

/**
 * Find blog posts
 * @param {projectConfig} config
 */
function loadPosts(config) {
    let debug = preparserDebug.new('posts');
    debug('Loading Posts');

    let root = config.compiler.posts_folder_absolute;

    return fileUtils.globFiles(root, '*/')
        .then((posts) => {
            let awaiting = [];
            let counter = 0;

            for (let i = 0; i < posts.length; i++) {
                let post = posts[i];

                // Don't take ignored files
                if (config.compiler.ignore_absolute.indexOf(post.absolute) != -1)
                    continue;

                config.__cache.posts.push({
                    info: post,
                    slug: post.name,
                    post: {
                        info: null,
                        contents: null,
                        shouldParse: true
                    },
                    assets: [],
                    contents: null
                });
                awaiting.push(_loadPost(config, config.__cache.posts.length - 1, debug));
            }

            return Promise.all(awaiting);
        });
}



module.exports = {
    loadData,
    loadTemplates,
    loadIncludes,
    loadPlugins,
    loadCustomProperties,
    discoverFiles,
    loadPosts
};