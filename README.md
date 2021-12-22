# Hitomezashi

Generate a [sashiko hitomezashi](https://en.wikipedia.org/wiki/Sashiko) pattern for a public github repository.

![](https://hitomezashi.witchert.co.uk/?owner=witchert&repo=hitomezashi&branch=master&size=25)

## Installation
```sh
> docker build --tag hitomezashi .
> docker run -p 1111:1111 hitomezashi
```
Navigate to http://localhost:1111/?owner=witchert&repo=hitomezashi&branch=master&size=24

### Query parameters
* **owner**: *required*, github repository owner, i.e. `witchert`
* **repo**: *required*, github repostirory name, i.e. `hitomezashi`
* **branch**: *required*, git branch name, i.e. `master`
* **size**: size of graphic, defaults to 24
