// This is heavily ChatGPT'd - apologies in advance 

import axios from 'axios';
import { promises as fs } from 'fs';
import { createObjectCsvWriter } from 'csv-writer';

// Environment variables
const token = process.env.GH_TOKEN; // Ensure GH_TOKEN is set in your environment
const org = process.env.TARGET_ORG; // Ensure TARGET_ORG is set in your environment
const baseURL = `https://api.github.com`;
const headers = {
    Authorization: `token ${token}`,
    'User-Agent': 'request'
};

// Base directory for output
const baseDir = './gh_data';

// Create the base directory if it doesn't exist
async function ensureBaseDir() {
    try {
        await fs.access(baseDir);
    } catch (error) {
        await fs.mkdir(baseDir, { recursive: true });
    }
}

// CSV writer setup
const csvWriter = createObjectCsvWriter({
    path: `${baseDir}/dockerfiles_paths.csv`,
    header: [
        { id: 'repo', title: 'REPOSITORY' },
        { id: 'path', title: 'PATH' }
    ]
});

// Fetch repositories from GitHub API
async function fetchRepos(org) {
    try {
        const response = await axios.get(`${baseURL}/orgs/${org}/repos`, { headers });
        return response.data.map(repo => ({ name: repo.name, url: repo.trees_url.replace('{/sha}', '/main?recursive=1') })); // main branch, put master if that's what you use
    } catch (error) {
        console.error('Error fetching repositories:', error);
        return [];
    }
}

// Search for Dockerfiles in repositories
async function searchDockerfiles(repos) {
    let dockerfilesPaths = [];
    for (const repo of repos) {
        try {
            const response = await axios.get(repo.url, { headers });
            const files = response.data.tree.filter(file => file.path.includes('Dockerfile'));
            if (files.length > 0) {
                const repoDir = `${baseDir}/${repo.name}`;
                await fs.mkdir(repoDir, { recursive: true });
                for (const file of files) {
                    const contentResponse = await axios.get(file.url, { headers });
                    const content = Buffer.from(contentResponse.data.content, 'base64').toString('utf8');
                    console.log(`Found Dockerfile in ${repo.name}/${file.path}`);
                    await fs.writeFile(`${repoDir}/${file.path.replace(/\//g, '_')}`, content);
                    dockerfilesPaths.push({ repo: repo.name, path: file.path });
                }
            }
        } catch (error) {
            console.error(`Error searching Dockerfiles in ${repo.name}:`, error);
        }
    }
    return dockerfilesPaths;
}

// Main function to run the script
async function main() {
    await ensureBaseDir();
    const repos = await fetchRepos(org);
    const dockerfilesPaths = await searchDockerfiles(repos);
    await csvWriter.writeRecords(dockerfilesPaths);
    console.log('CSV file has been written with paths of each Dockerfile.');
}

main();