let path = require('path');
let fs = require('fs');
let glob = require('glob');
let debug = require('../debugger')('project_config');


var __projectRoot = null;
var __configFilePath = null; // Path to the config file
var __parsed_config = null;

// defaults
var __default_site = {
    title: "No title",
    description: "No description",
    author: "Unkown",
    author_email: "",
    url: "",
    language: "en",

    source: ".",
    output: "_build",
    default_template: "default",

    show_drafts: false, /* If true, then drafts will show up same as finished posts */

    post_permalink: "/blog/%slug%", // %category%
    read_more: "<!-- more -->"
};

var __default_compiler = {
    plugins_folder: "_plugins",
    templates_folder: "_templates",
    includes_folder: "_includes",
    data_folder: "_data",
    posts_folder: "_posts",

    include: [], /* Folders/files to include in compile (globs) */
    ignore: ['*.daphne', 'README.md'],  /* Folders/files to NOT include in compile (globs) */

    include_no_compile: [], /* Folders/Files to copy to output, but not to compile */

    allow_plugins: false,
    ignore_dot_names: true,

    extensions_to_parse: ['html', 'htm', 'txt'],

    tags: {
        delimeter: "---",
        opening: "{%",
        closing: "%}",
        print_opening: "{{",
        print_closing: "}}"
    }
};




/**
 * Loads configuration file after setting project root
 * @param {string} root Project root
 */
function __loadConfigFile(root) {
    try {
        __configFilePath = path.join(root, 'config.daphne');
        let contents = fs.readFileSync(__configFilePath);
        __parsed_config = JSON.parse(contents);

        debug(`Read project config from '${__configFilePath}'`);

        if (__parsed_config.site == undefined || __parsed_config.site == null)
            __parsed_config.site = {};

        if (__parsed_config.compiler == undefined || __parsed_config.compiler == null)
            __parsed_config.compiler = {};

        // Deep copy two parts of the config
        this.site = JSON.parse(JSON.stringify(__parsed_config.site));
        this.compiler = JSON.parse(JSON.stringify(__parsed_config.compiler));
    } catch (e) {
        throw new Error(`Unable to find or read 'config.daphne' file in '${root}'`);
    }
}

/**
 * Applies default configuration to something
 * @param {object} target
 * @param {object} defaults
 */
function __applyDefaults(target, defaults) {
    if (Array.isArray(target) && Array.isArray(defaults)) {
        for (let i = 0; i < defaults.length; i++) {
            if (target.indexOf(defaults[i]) == -1)
                target.push(defaults[i]);
        }
    }
    else if (typeof(target) == 'object' && typeof(defaults) == 'object') {
        for (let key in defaults) {
            if (!defaults.hasOwnProperty(key))
                continue;

            if (Array.isArray(defaults[key])) {
                if (target[key] == undefined || !Array.isArray(target[key]))
                    target[key] = [];

                __applyDefaults(target[key], defaults[key]);
            }
            else if (typeof(defaults[key]) == 'object') {
                if (target[key] == undefined || typeof(target[key]) != 'object')
                    target[key] = {};

                __applyDefaults(target[key], defaults[key]);
            }
            else {
                if (target[key] == undefined)
                    target[key] = defaults[key];
            }
        }
    }
}

/**
 * Expands paths to absolute paths
 * @param {object} target
 * @param {string} key
 */
function __expandPaths(target, key) {
    if (Array.isArray(key)) {
        for (let i = 0; i < key.length; i++) {
            __expandPaths(target, key[i]);
        }
    }
    else {
        target[key] = path.normalize(target[key]);
        target[key + '_absolute'] = path.join(__projectRoot, target[key]);
    }
}

/**
 * Private method to expand globs of a config file after loading
 */
function __expandGlobs(target, key, root) {
    let new_files = [];

    for (let i = 0; i < target[key].length; i++) {
        target[key][i] = path.normalize(target[key][i]);

        let globPath = path.join(root, target[key][i]);
        new_files.push(globPath);

        let files = glob.sync(globPath);
        for (let j = 0; j < files.length; j++) {
            new_files.push(path.normalize(files[j]));
        }
    }

    target[key + '_absolute'] = new_files.concat([]); // Deep copy
}


/**
 * Class for managing a project configuration
 */
class ProjectConfig {

    constrcutor() {
        // Apply defaults
        this.site = {};
        this.compiler = {
            root: '',
        };

    }

    /**
     * Parses a configuration file for a project
     * @param {string} projectRoot folder path
     */
    setProjectRoot(projectRoot) {
        __projectRoot = path.normalize(projectRoot);

        // Load the config file and apply defaults
        __loadConfigFile.call(this, __projectRoot);
        __applyDefaults(this.site, __default_site);
        __applyDefaults(this.compiler, __default_compiler);

        // Expand all paths
        __expandPaths(this.site, [
            'source',
            'output'
        ]);

        __expandPaths(this.compiler, [
            'plugins_folder',
            'templates_folder',
            'includes_folder',
            'data_folder',
            'posts_folder'
        ]);

        // Expand globbed paths
        __expandGlobs(this.compiler, 'include', this.site.source_absolute);
        __expandGlobs(this.compiler, 'ignore', this.site.source_absolute);
        __expandGlobs(this.compiler, 'include_no_compile', this.site.source_absolute);

        this.compiler.root = __projectRoot;
    }


    /**
     * Returns true if absolute path of file is ignored
     * @param {string} filePath
     */
    isFileIgnored(filePath) {
        return (this.compiler.ignore_absolute.indexOf(filePath) != -1);
    }

    /**
     * Returns the project folder
     */
    get projectRoot() {
        return __projectRoot;
    }

    /**
     * Returns the website folder
     */
    get websiteRoot() {
        return this.site.source_absolute;
    }

    /**
     * Output folder
     */
    get outputRoot() {
        return this.site.output_absulte;
    }
}

// Export as a singleton
let instance = new ProjectConfig();
module.exports = instance;