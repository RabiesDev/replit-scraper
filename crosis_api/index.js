const { Crosis } = require('crosis4furrets');
const express = require('express');
const router = express();

function asyncWrapper(fn) {
    return (req, res, next) => {
        return Promise.resolve(fn(req))
            .then((result) => res.send(result))
            .catch((err) => res.status(400).send(err));
    };
}

function hasValidFileExtension(filePath) {
    const extensionIndex = filePath.lastIndexOf('.');
    if (extensionIndex === -1 || extensionIndex === filePath.length - 1) {
        return 0;
    }

    const fileName = filePath.substring(0, extensionIndex).toLowerCase();
    const extension = filePath.substring(extensionIndex + 1).toLowerCase();
    if (fileName.length === 0 || extension.length <= 1) {
        return 0;
    }

    const badExtensions = ["nix", "mp3", "mp4", "avi", "png", "jpg", "jpeg", "gif", "css", "lock", "exe", "jar", "class", "toml", "md", "sh", "dockerfile", "draw", "ico", "deployment", "development", "dll", "ttf", "utf"];
    const badFileNames = ["requirements", "package-lock", "package", "tsconfig", "tsconfig-base", "cypress-install", "cypress", "docker-compose", "pom", "pyproject", "bukkit", "pnpm-lock", "dockerfile", "webpack", "config.example", "pm2.example", "tslint", "yarn-error"];
    if (badExtensions.includes(extension) || badFileNames.includes(fileName)) {
        return -1;
    }
    return 1;
}

async function processDirectory(client, repo, path) {
    const directories = await client.readdir(path);
    const readDirectoryContent = async (directory) => {
        const extensionResult = hasValidFileExtension(directory);
        if (extensionResult === 1) {
            console.log("[+] Reading file... (" + repo + ", " + directory + ")");
            const buffer = await client.read(directory);
            return { path: directory, content: buffer.toString() };
        } else if (extensionResult === 0) {
            return { path: directory, content: "" };
        }
    };
    return await Promise.all(directories.map(readDirectoryContent));
}

router.get('/directory', asyncWrapper(async (req) => {
    if (!req.query.repo || !req.query.token || !req.query.path) {
        throw new Error("Missing parameters");
    }

    try {
        const client = new Crosis({ replId: req.query.repo, token: req.query.token });
        try {
            await client.connect();
            if (!client.connected) {
                throw new Error("Connection failure");
            }
            console.log("[+] Reading directory... (" + req.query.repo + ", " + req.query.path + ")");
            return JSON.stringify(await processDirectory(client, req.query.repo, req.query.path));
        } catch (err) {
            throw err; // Propagate the error to be caught by asyncWrapper
        } finally {
            client.close(); // Make sure to close the connection
        }
    } catch (err) {
        throw err;
    }
}));

router.get('/file', asyncWrapper(async (req) => {
    if (!req.query.repo || !req.query.token || !req.query.path) {
        throw new Error("Missing parameters");
    }

    try {
        const client = new Crosis({ replId: req.query.repo, token: req.query.token });
        try {
            await client.connect();
            if (!client.connected) {
                throw new Error("Connection failure");
            }
            console.log("[+] Reading file... (" + req.query.repo + ", " + req.query.path + ")");
            return await client.read(req.query.path);
        } catch (err) {
            throw err; // Propagate the error to be caught by asyncWrapper
        } finally {
            client.close(); // Make sure to close the connection
        }
    } catch (err) {
        throw err;
    }
}));

router.listen(3000, () => {
    console.log("[+] Start on port 3000");
});
