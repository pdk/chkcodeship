# chkcodeship

Talks to codeship API to get the latest successful build of a particular branch.

Example:

    CODESHIP_LOGIN=email@example.com CODESHIP_PASSWORD=foobarackbbq chkcodeship orgname projname master
    fcb762e814ebd5679161ffdc75ea263a92afc002 2020-09-08 20:57:11.196 +0000 UTC

The output is the commit SHA followed by the finished at timestamp from
codeship.

If no successful build are found, program exits with an error.

Limitations:

1. Currently only fetches the first page of 50 builds.
2. Seems like codeship's go library is not module-ready. Have to build/run old
   skool.
