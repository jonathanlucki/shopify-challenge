# Shopify Image Repo Challenge
This API provides simple functionality for the uploading/naming/deletion of images in a repository. The uploading, retreival, and deletion actions can be done via the REST endpoints listed below. These endpoints can be accessed via a test server running at `http://shopify.jonathanlucki.ca/`. Image uploading can also be tested via [the upload form here](http://shopify.jonathanlucki.ca/public/), and you may also view all images currently in the repo and also delete them [here](http://shopify.jonathanlucki.ca/public/images.html).

## API Endpoints
### Get all images
GET `/images`
##### Example response:
```
{
    "images": [
        {"Url": "https://storage.googleapis.com/image-repo-jlucki/ZbyJPQW6wane7KMKa.png",
        "Id": "ZbyJPQW6wane7KMKa",
        "Name": "ImageOneName",
        "Date": "Mon May 9 08:00:00 UTC 2021"},
        {"Url": "https://storage.googleapis.com/image-repo-jlucki/JPQW6wane7KMKahjk.png",
        "Id": "JPQW6wane7KMKahjk",
        "Name": "ImageTwoName",
        "Date": "Mon May 9 09:00:00 UTC 2021"},
        ]
}
```

### Upload image
POST `/images`
(Upload by multipart/form-data content-type)
##### Example response:
```
{
    "Url": "https://storage.googleapis.com/image-repo-jlucki/MVqxD5vxsQHrKnr2yoL.png",
    "Id": "MVqxD5vxsQHrKnr2yoL"
}
```

### Get image (using id)
GET `/images/:id`
##### Example response:
```
{
    "Url": "https://storage.googleapis.com/image-repo-jlucki/ZbyJPQW6wane7KMKa.png",
    "Id": "ZbyJPQW6wane7KMKa",
    "Name": "ImageOneName",
    "Date": "Mon May 9 08:00:00 UTC 2021"
}
```

### Delete image (using id)
DELETE `/images/:id`
##### Example response:
```
{
    "Id": "MVqxD5vxsQHrKnr2yoL"
}
```

## Possible future changes/additions
- Cache to speed up response time for frequently requested images
- Pagination of image data responses
- Ability to search images via name
- Users/permissions with ownership of images
- Support for additional image/file formats
- Ability to comment on images