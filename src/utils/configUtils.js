let fs = require('fs');
let glob = require('glob');
let path = require('path');

let defaultConfig = {
    site: {
        title: "No title",
        description: "No description",
        author: "Unkown",
        author_email: "",
        url: "",

        source: ".",
        output: "_build",
        default_template: "default",

        show_drafts: false, /* If true, then drafts will show up same as finished posts */

        permalink: "/blog/%slug%",
        read_more: "<!-- more -->"
    },
    compiler: {
        plugins_folder: "_plugins",
        templates_folder: "_templates",
        includes_folder: "_includes",
        data_folder: "_data",

        include: [], /* Folders/files to include in compile (globs) */
        ignore: [],  /* Folders/files to NOT include in compile (globs) */

        include_no_compile: [], /* Folders/Files to copy to output, but not to compile */

        allow_plugins: false,
        ignore_dot_names: true,

        tags: {
            delimeter: "---",
            opening: "{%",
            closing: "%}",
            print_opening: "{{",
            print_closing: "}}"
        }
    }
};


/**
 * Expand globs and get files that match
 * @param {projectConfig} config Project configuration
 */
function _expandGlobs(config) {
    let expandThem = (cfg, key, root) => {
        let new_files = [];

        for (let i = 0; i < cfg[key].length; i++) {
            let globPath = path.join(root, cfg[key][i]);
            globPath = path.normalize(globPath);
            new_files.push(globPath);

            let files = glob.sync(globPath);
            if (files.length > 0)
                new_files = new_files.concat(files);
        }

        // Overwrite key, but keep olld copy
        cfg[key + '_absolute'] = new_files.concat([]);
    };

    expandThem(config.compiler, 'include', config.site.source_absolute);
    expandThem(config.compiler, 'ignore', config.site.source_absolute);
    expandThem(config.compiler, 'include_no_compile', config.site.source_absolute);
}



/**
 * Reads and parses a config file
 * @param {string} filePath Path to config file
 * @returns {json} Parsed object
 */
function parseConfigFile(filePath) {
    let contents = fs.readFileSync(filePath);
    let parsed = JSON.parse(contents);
    return parsed;
}


/**
 * Applies default values from defaultConfig to config
 * @param {json} config Parsed config file object
 */
function applyDefaultConfiguration(config) {
    let applyArray = (defaults, cfg) => {
        for (let i = 0; i < defaults.length; i++) {
            if (cfg.indexOf(defaults[i]) == -1)
                cfg.push(defaults[i]);
        }
    };

    // Function to apply everything
    let applyConfiguration = (defaults, cfg) => {
        for (let key in defaults) {
            if (!defaults.hasOwnProperty(key))
                continue;

            if (Array.isArray(defaults[key])) {
                if (cfg[key] == undefined || !Array.isArray(cfg[key]))
                    cfg[key] = [];
                applyArray(defaults[key], cfg[key]);
            }
            else if (typeof(defaults[key]) == 'object') {
                if (cfg[key] == undefined || typeof(cfg[key]) != 'object')
                    cfg[key] = {};
                applyConfiguration(defaults[key], cfg[key]);
            }
            else {
                cfg[key] = defaults[key];
            }


        }
    };

    // Start at root of both
    applyConfiguration(defaultConfig, config);

    // Expand site source and site root
    let expandPath = (cfg, key, root) => {
        cfg[key + '_absolute'] = path.join(root, cfg[key]);
    };
    expandPath(config.site, 'source', config.compiler.root);
    expandPath(config.site, 'output', config.compiler.root);
    expandPath(config.compiler, 'plugins_folder', config.compiler.root);
    expandPath(config.compiler, 'templates_folder', config.compiler.root);
    expandPath(config.compiler, 'includes_folder', config.compiler.root);
    expandPath(config.compiler, 'data_folder', config.compiler.root);

    _expandGlobs(config);
}


/**
 * Checks if a project has a config file
 * @param {string} projectRoot Project root path
 */
function doesProjectHaveConfigFile(projectRoot) {
    try
    {
        let files = fs.readdirSync(projectRoot);

        for (let i = 0; i < files.length; i++) {
            if (files[i].toLocaleLowerCase() == 'config.daphne')
                return true;
        }

        return false;
    }
    catch (e)
    {
        return false;
    }
}





module.exports = {
    parseConfigFile,
    applyDefaultConfiguration,
    doesProjectHaveConfigFile
};