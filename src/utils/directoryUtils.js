var fs = require('fs');

/**
 * Checks if a folder has a daphne.config file
 * @param {string} directory Path to static site files
 * @returns bool
 */
function doesDirectoryHaveConfigFile(directory) {
    try
    {
        let files = fs.readdirSync(directory);
        for (let i = 0; i < files.length; i++)
        {
            if (files[i].toLocaleLowerCase() == "config.daphne")
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
    doesDirectoryHaveConfigFile
};