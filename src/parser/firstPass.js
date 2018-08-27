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
    debug(`Expanding file includes for ${file.info.name}`);



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

    if (file.info.name == '404.html')
        console.log(newContents);
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