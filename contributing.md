# contributing

I made this package because I needed the functionality of deleting files, so that's what's implemented right now.

While I plan to add the rest of the Uploadthing api's functionality, it's not a priority of mine currently. You are more than welcome though!

## scope

This package should do nothing more than the [UTApi class in Uploadthing](https://github.com/pingdotgg/uploadthing/blob/main/packages/uploadthing/src/sdk/index.ts#L39)

## dev setup

1. If you'd like to add a feature, first clone the [Uploadthing](https://github.com/pingdotgg/uploadthing) repo and find its current implementation.

    ```bash
    git clone --depth 1 https://github.com/pingdotgg/uploadthing uploadthing_source;
    cd uploadthing_source;
    ## I recommend starting in the below file - the core functionality is there
    vim packages/uploadthing/src/sdk/index.ts
    ```

2. Clone this repo and start a branch for your feature.

    ```bash
    git clone https://github.com/jesses-code-adventures/utapi-go;
    cd utapi-go;
    git checkout -b my_feature;
    ```

3. Make your changes

4. Check for any upstream changes

    ```bash
    git pull --rebase origin master
    ```

5. Push your branch to origin

    ```bash
    git push origin my_feature
    ```

6. Go to Github
    - Create a pull request for your branch
    - In your pull request, explain your changes and link any relevant issues.
    - Wait for a review/test
