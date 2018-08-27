let debug = require('../debugger')('parser:first');
let fs = require('fs');
let path = require('path');
let fileUtils = require('../utils/fileUtils');



/**
 * Expand the include statements inside of a file
 * @param {object} config
 * @param {object} file
 */
function _expandIncludes(config, file) {
    debug(`Expanding includes for ${file.info.name}`);

    let cmd = new RegExp(config.compiler.tags.opening + '(?:\\s*)include\\s(.*)?(?:\\s*)' + config.compiler.tags.closing, 'gi');
    let includes = file.contents.match(cmd);
    if (includes == null || includes == undefined)
        includes = []; // Don't break things

    let actually_included = [];

    for (let i = 0; i < includes.length; i++) {
        let fileName = includes[i].replace(config.compiler.tags.opening, '')
                            .replace(config.compiler.tags.closing, '')
                            .replace('include', '')
                            .trim();

        if (config.__site.includes.indexOf(fileName) == -1) {
            debug(`Unkown include of '${fileName}' in ${file.info.relative}`);
            continue;
        }

        if (actually_included.indexOf(fileName) != -1)
            continue; // Two of the same include statements in the file


        // Make sure file has valid cache
        if (config.__cache.includes[fileName] == undefined || config.__cache.includes[fileName].contents == null) {
            debug(`Error including '${fileName}'`);
            // TODO: retrieve the file
            continue;
        }

        actually_included.push(fileName);
        file.contents = file.contents.replace(new RegExp(includes[i], 'gi'), config.__cache.includes[fileName].contents);

        debug(`Including '${fileName}' into ${file.info.relative}`);
    }

    if (actually_included.length > 0)
        return _expandIncludes(config, file); // Re-run

    return Promise.resolve(file);
}

/**
 * Expand the object into the entire template
 * @param {object} config
 * @param {object} file
 */
function _expandTemplate(config, file) {
    debug(`Expanding file template for ${file.info.name}`);

    let template = file.metadata.template;
    if (template == undefined)
        return Promise.reject(`No template specified for ${file.info.relative}`);

    if (config.__site.templates.indexOf(template) == -1)
        return Promise.reject(`Unkown template '${template}' in ${file.info.relative}`);

    let cachedTemplate = config.__cache.templates[template];
    if (cachedTemplate == undefined || cachedTemplate.contents == null)
        return Promise.reject(`Unkown error while loading template '${template}'`);

    let cmd = new RegExp(config.compiler.tags.print_opening + '(?:\\s*)content(?:\\s*)' + config.compiler.tags.print_closing, 'g');
    let newContents = cachedTemplate.contents;
    newContents = newContents.replace(cmd, file.contents);

    file.contents = newContents; // Replace the contents of the file

    return Promise.resolve(file);
}



/**
 * Retrieves metadata about the file
 * @param {object} config
 * @param {object} file
 */
function _discoverFileMetadata(config, file) {
    debug(`Discoving file metadata for ${file.info.name}`);
    return new Promise((resolve, reject) => {
        fileUtils.getMetadataHeader(file, config, debug);
        resolve(file);
    });
}


/**
 * Executes the first pass on a file
 * @param {object} config
 * @param {object} file
 */
function _performFirstPass(config, file) {
    return _discoverFileMetadata(config, file)
            .then((f) => { return _expandTemplate(config, f); })
            .then((f) => { return _expandIncludes(config, f); })
            .catch((e) => {
                console.error(e);
                return Promise.resolve();
            });

}


/**
 * First pass of the parser
 * First pass discovers metadata about the file
 * @param {projectConfig} config project config
 */
function firstPass(config) {
    debug('First Pass');

    let awaiting = [];

    for (let i = 0; i < config.__cache.files.length; i++) {
        let file = config.__cache.files[i];
        if (!file.shouldParse) {
            debug(`Skipping ${file.info.name}`);
            continue;
        }
        awaiting.push(_performFirstPass(config, file));
        //fileUtils.getMetadataHeader(config.__cache.files[i], config, debug);
    }

    return Promise.all(awaiting);
}

module.exports = firstPass;