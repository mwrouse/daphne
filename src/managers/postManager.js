let path = require('path');
let debug = require('../debugger')('post_manager');
let fileUtils = require('../utils/fileUtils');
let FileToCompile = require('../models/fileToCompile');
let config = require('./projectConfig');


var __postCache = [];

/**
 * Load an asset to a post
 * @param {string} file
 */
function __loadAsset(file) {
    return fileUtils.readEntireFile(file)
        .then((content) => {
            return new FileToCompile(file, content);
        });
}

/**
 * Load post and it's assets
 * @param {string} folder
 */
function __loadPostAndAssets(folder) {
    return fileUtils.globFiles(folder, '**/*.*')
        .then((files_found) => {
            let post_file_load = null;
            let assets_loading = [];

            for (let i = 0; i < files_found.length; i++) {
                let file = files_found[i];

                if (config.compiler.ignore_absolute.indexOf(file.absolute) != -1)
                    continue; // Ignore the file

                let extension = path.extname(file.absolute);
                let nameNoExtension = path.basename(file.absolute).replace(extension, '');
                if (nameNoExtension == 'index') {
                    post_file_load = __loadAsset(file.absolute, debug);
                }
                else {
                    assets_loading.push(
                        __loadAsset(file.absolute, debug)
                    );
                }
            }

            if (post_file_load == null)
                throw new Error(`Post does not contain an index.html file`);

            return post_file_load.then((post_file) => {
                // Wait until all assets are done loading
                return Promise.all(assets_loading)
                    .then((assets) => {
                        // Set assets on the post_file
                        post_file.__assets = post_file.__assets.concat(assets);
                        return post_file;
                    });
            });

        });
}


class PostManager {
    constructor(){
    }

    /**
     * Load all the posts
     */
    loadPosts() {
        let root = config.compiler.posts_folder_absolute;

        return fileUtils.globFiles(root, '*/')
            .then((post_folders) => {
                let waiting = [];
                for (let i = 0; i < post_folders.length; i++) {
                    let post = post_folders[i];

                    if (config.compiler.ignore_absolute.indexOf(post.absolute) != -1)
                        continue; // Ignore folder

                    debug(`Loading post in folder '${post.relative}'`);
                    waiting.push(
                        __loadPostAndAssets(post.absolute)
                    );
                }

                return Promise.all(waiting)
                    .then((posts) => {
                        __postCache = __postCache.concat(posts);
                    });
            });
    }


    /**
     * Get list of the posts
     */
    get posts() {
        let posts = [];
        for (let i = 0; i < __postCache.length; i++)
            posts.push(__postCache[i].context);

        return posts;
    }
}

// Singelton
let instance = new PostManager();
module.exports = instance;