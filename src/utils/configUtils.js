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
        template_folder: "_templates",
        includes_folder: "_includes",
        data_folder: "_data",

        include: [], /* Folders/files to include in compile (globs) */
        ignore: [],  /* Folders/files to NOT include in compile (globs) */

        include_no_compile: [], /* Folders/Files to copy to output, but not to compile */

        allow_plugins: false,
        ignore_dot_names: true,

        tags: {
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
function expandGlobs(config) {
    let expandThem = (cfg, key) => {
        let new_files = [];
        for (let i = 0; i < cfg[key].length; i++) {
            let globPath = path.join(config.root, cfg[key][i]);
            globPath = path.normalize(globPath);

            let files = glob.sync(globPath);
            if (files.length > 0)
                new_files = new_files.concat(files);
        }

        // Overwrite key, but keep olld copy
        cfg[key + '_globs'] = cfg[key].concat([]);
        cfg[key] = new_files.concat([]);
    };

    expandThem(config, 'include');
    expandThem(config, 'ignore');
    expandThem(config, 'include_no_compile');
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
 * @param {json} projectConfig Parsed config file object
 */
function applyDefaultConfiguration(projectConfig) {
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
    applyConfiguration(defaultConfig, projectConfig);

    expandGlobs(projectConfig.compiler);
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