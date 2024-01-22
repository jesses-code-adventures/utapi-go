# contributing

Currently this package has parity with the [UTApi class in Uploadthing](https://github.com/pingdotgg/uploadthing/blob/main/packages/uploadthing/src/sdk/index.ts#L39), so this package doesn't require new features.

If you have modifications you'd like to make to the code for cleanliness or other improvements, I'd recommend getting in touch before making the changes.

## scope

This package should do nothing more than the [UTApi class in Uploadthing](https://github.com/pingdotgg/uploadthing/blob/main/packages/uploadthing/src/sdk/index.ts#L39)

## utapi features

- [x] requestUploadThing
- [x] deleteFiles
- [x] getFileUrls
- [x] listFiles
- [x] renameFiles
- [x] getSignedURL
- [x] getUsageInfo

## dev setup

1. If you'd like to make a change, first clone the [Uploadthing](https://github.com/pingdotgg/uploadthing) repo and find the original function you plan to modify.

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
