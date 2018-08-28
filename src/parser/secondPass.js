let debug = require('../debugger')('parser:second');
let fs = require('fs');
let path = require('path');
let fileUtils = require('../utils/fileUtils');
let directoryUtils = require('../utils/directoryUtils');


/**
 * Generates data to be sent when rendering the page
 * @param {projectConfig} config
 * @param {object} page
 */
function _getPageData(config, page) {
    let data = {
        site: config.site,
        page: page.metadata
    };

    // Load site properties
    for (let i = 0; i < config.__site.properties.length; i++) {
        let key = config.__site.properties[i];
        data.site[key] = config.__cache.site[key];
    }

    // Load site posts
    data.site.posts = [];
    for (let i = 0; i < config.__cache.posts.length; i++) {
        data.site.posts.push(config.__cache.posts[i]);
    }

    console.log(JSON.stringify(data));
    return data;
}



/**
 * Performs the second pass on the posts
 * @param {projectConfig} config
 */
function _secondPassPosts(config) {
    return new Promise((resolve, reject) => {
        resolve();
    });
}

/**
 * Performs the second pass on pages
 * @param {projectConfig} config
 */
function _secondPassPages(config) {
    return new Promise((resolve, reject) => {
        for (let i = 0; i < config.__cache.files.length; i++) {
            let cache = config.__cache.files[i];

            directoryUtils.createDirectoryStructure(config.site.output_absolute, cache.info.relativeDirname);

            let finalPath = path.join(config.site.output_absolute, cache.info.relative);

            if (!cache.shouldParse) {
                fileUtils.copyFile(finalPath, cache.info.absolute);
                continue;
            }

            let data = _getPageData(config, cache);

            // Parse the file and write it
            fs.writeFileSync(finalPath, cache.contents);
        }

        resolve();
    });
}

/**
 * Second pass of the parser
 * Second pass actually expands the body of the file
 * @param {projectConfig} config project config
 */
function secondPass(config) {
    debug('Second Pass');

    return Promise.all([
        _secondPassPosts(config),
        _secondPassPages(config)
    ]);
}

module.exports = secondPass;