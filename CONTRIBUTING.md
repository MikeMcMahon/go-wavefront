# Contributing

We welcome contributors to this Golang Client, and we'll do our best to review and merge all requests.
Adding missing features as per the [Wavefront API](https://www.wavefront.com/api/) or bug fixes will be welcomed.
Any functional changes will require discussion.

## Opening Issues
If you encounter a bug or you are making a feature request, please open an issue in this repo.

## Making Pull Requests
1. Fork the repository
1. Create a new branch for your change
1. Make your changes and submit a [Pull Request](https://help.github.com/articles/creating-a-pull-request-from-a-fork/)

Before submitting a pull request, please ensure that unit tests pass. Refer to the [README.md](./README.md) for instructions on running unit tests.

We will review your pull request and provide feedback.

## Versioning

We use [Semantic Versioning](http://semver.org/) on this project. The version is located inside the `version` file, in the root of the repository, in the format `vMajor.Minor.Patch`. Update this version as required.

## Creating a new Release

1. Update the CHANGELOG.md
1. Update the `version` file to X.Y.Z
1. Commit changes.
1. Make a new tag (`git tag vX.Y.Z`)
1. Push changes / tag vX.Y.Z (`git push --tags`)
1. Publish GitHub a release using that tag.
    1. Add a summary of the release, as in prior releases.