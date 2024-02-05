# Leaky Vessels Image SBOM Generator

Inspired by (https://github.com/snyk/leaky-vessels-static-detector)


A tool to generate an image SBOM from Dockerfiles. Currently only extracts base images.

## Overview

### Why?

When assessing impact to security threats such as the Docker runc container vulnerability, the first thing is to inventory deployed assets. This tool scans your Dockerfiles and reports the detected base images in each, allowing for a better supply chain catalog of image dependencies.

This is typically found in SCA tools like Snyk or Semgrep, or platforms like Docker Scout. However, parsing the AST is quite easy in a pinch without these costly vendors.

### Quickstart

1. Collect Dockerfiles from Github:
  - Navigate to `gh_crawler/src/gh_enumerator` and run `npm install`. 
  - I wasn't able to get Snyk's bash script working so improvising with JS.
  - Run `node sbom_generator_github_crawler.js`
2. Verify outputs
  - Verify the script created an output file `gh_crawler/src/gh_enumerator/dockerfiles_paths.csv`
  - Verify the script created folders for each repo where it detected a Dockerfile under `gh_crawler/src/gh_enumerator/gh_data`

3. Run the SBOM Generator Go script.
  - `go run cmd/sbom-generator/sbom_generator.go dockerfile --base`


### Support/Security

1. No expectation of support is provided with this and I assume no liability for running it. 