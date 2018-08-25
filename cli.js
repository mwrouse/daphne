#!/usr/bin/env node
let Daphne = require('./src/daphne.js');


const directory = process.cwd();
const argv = process.argv.slice(2);
const argc = argv.length;


let siteGenerator = new Daphne();

// Enable debugging or not
switch (argv[argc - 1])
{
    case "--debug":
        let debug = require('./src/debugger');
        debug.enableAll();
        break;
}


switch (argv[0])
{
    case "build":
        siteGenerator.buildSite(directory);
        break;

    case "watch":
        siteGenerator.watchSite(directory);
        break;

    case "serve":
        siteGenerator.serveSite(directory);
        break;


    // Help, unkown argument, or no argument
    case "help":
    default:
        if (argv[0] != "help" && argv[0] != undefined)
            console.log(`Unkown argument '${argv[0]}'`);

        console.log("help");
        break;
}

