# Go Linkk

Custom URL redirection on Go AppEngine.

`go.blinkk.com/grow` --> `grow.io`

# Usage

## Requirements

This project requires the following:

-  [Grow](https://grow.io) to build the ui.
-  [gcloud](https://cloud.google.com/sdk/gcloud/) to work with AppEngine and deploy.

## Initial Setup

After cloning the repository copy the `app-example.yaml` to `app.yaml` and update the `AUTH_DOMAINS` variable to match the domain name for administrators. Multiple domains can be separated by a `|`.

Generate the UI by running `grow install` followed by `grow build` in the `ui/` directory. Any time files change in the `ui/` directory the `grow build` command will need to be run for AppEngine development and deployments to pick up the changes.

## Deploying

Deploy to an AppEngine project like so:

    gcloud app deploy --project=<project-id>

## Local Development

Run a local server for development by running:

    dev_appserver.py app.yaml
